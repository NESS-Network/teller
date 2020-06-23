package sender

import (
	"errors"

	"github.com/MDLlife/MDL/src/api"
	"github.com/MDLlife/MDL/src/wallet"
	"github.com/MDLlife/MDL/src/readable"
	"strings"
	"github.com/MDLlife/MDL/src/util/droplet"
)

// APIError wraps errors from the mdl CLI/API library
type APIError struct {
	error
}

// NewAPIError wraps an err with APIError
func NewAPIError(err error) APIError {
	return APIError{err}
}

// RPC provides methods for sending coins
type API struct {
	walletFile string
	changeAddr string
	apiClient  *api.Client
}

// NewRPC creates RPC instance
func NewAPI(wltFile, apiAddr string) (*API, error) {
	wlt, err := wallet.Load(wltFile)
	if err != nil {
		return nil, err
	}

	if len(wlt.GetAddresses()) == 0 {
		return nil, errors.New("Wallet is empty")
	}

	apiClient := api.NewClient( "http://" + apiAddr + "/")

	wfs := strings.Split(wltFile,"/")
	wltFileName := wfs[len(wfs)-1]

	return &API{
		walletFile: wltFileName,
		changeAddr: wlt.GetAddresses()[0].String(),
		apiClient:  apiClient,
	}, nil
}

// CreateTransaction creates a raw MDL transaction offline, that can be broadcast later
func (c *API) CreateTransaction(recvAddr string, amount uint64) (*api.CreateTransactionResponse, error) {
	// TODO -- this can support sending to multiple receivers at once,
	// which would be necessary if the exchange was busy

	strCoins, err := droplet.ToString(amount)
	if err != nil {
		return nil, APIError{err}
	}

	strCoins = strCoins[:len(strCoins)-3]

	to := api.Receiver{Address: recvAddr, Coins: strCoins}
	req := api.WalletCreateTransactionRequest{ WalletID: c.walletFile }
	req.To = []api.Receiver{to}
	req.HoursSelection = api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.1"}

	createTxResp, err := c.apiClient.WalletCreateTransaction(req)

	if err != nil {
		return nil, APIError{err}
	}

	return createTxResp, nil
}

// BroadcastTransaction broadcasts a transaction and returns its txid
func (c *API) BroadcastTransaction(encodedTx string) (string, error) {
	txid, err := c.apiClient.InjectEncodedTransaction(encodedTx)
	if err != nil {
		return "", APIError{err}
	}

	return txid, nil
}

// GetTransaction returns transaction by txid
func (c *API) GetTransaction(txid string) (*readable.TransactionWithStatus, error) {
	txn, err := c.apiClient.Transaction(txid)
	if err != nil {
		return nil, APIError{err}
	}

	return txn, nil
}

// Balance returns the balance of a wallet
func (c *API) Balance() (*readable.BalancePair, error) {
	bal, err := c.apiClient.WalletBalance(c.walletFile)
	if err != nil {
		return nil, APIError{err}
	}

	return &bal.BalancePair, nil
}
