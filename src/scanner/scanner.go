package scanner

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/modeneis/waves-go-client/model"
	//"github.com/skycoin/skycoin/src/api"
	"github.com/MDLlife/MDL/src/readable"
	//"github.com/skycoin/skycoin/src/readable"
)

// Scanner provids apis for interacting with a scan service
type Scanner interface {
	AddScanAddress(string, string) error
	GetDeposit() <-chan DepositNote
}

// BtcRPCClient rpcclient interface
type BtcRPCClient interface {
	GetBlockVerboseTx(*chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlockHash(int64) (*chainhash.Hash, error)
	GetBlockCount() (int64, error)
	Shutdown()
}

// SkyRPCClient rpcclient interface
type SkyRPCClient interface {
	GetTransaction(txid string) (*readable.TransactionWithStatus, error)
	GetBlocks(start, end uint64) (*readable.Blocks, error)
	GetBlocksBySeq(seq uint64) (*readable.Block, error)
	GetLastBlocks() (*readable.Block, error)
	Shutdown()
}

// WavesRPCClient rpcclient interface
type WavesRPCClient interface {
	GetTransaction(txid string) (*model.Transactions, error)
	GetBlocks(start, end int64) (*[]model.Blocks, error)
	GetBlocksBySeq(seq int64) (*model.Blocks, error)
	GetLastBlocks() (*model.Blocks, error)
	Shutdown()
}

// EthRPCClient rpcclient interface
type EthRPCClient interface {
	GetBlockVerboseTx(seq uint64) (*types.Block, error)
	GetBlockCount() (int64, error)
	Shutdown()
}

// DepositNote wraps a Deposit with an ack channel
type DepositNote struct {
	Deposit
	ErrC chan error
}

// NewDepositNote returns a DepositNote
func NewDepositNote(dv Deposit) DepositNote {
	return DepositNote{
		Deposit: dv,
		ErrC:    make(chan error, 1),
	}
}

// Deposit struct
type Deposit struct {
	CoinType  string // coin type
	Address   string // deposit address
	Value     int64  // deposit amount. For BTC, measured in satoshis.
	Height    int64  // the block height
	Tx        string // the transaction id
	N         uint32 // the index of vout in the tx [BTC]
	Processed bool   // whether this was received by the exchange and saved
}

// ID returns $tx:$n formatted ID string
func (d Deposit) ID() string {
	return fmt.Sprintf("%s:%d", d.Tx, d.N)
}

// GetCoinTypes returns supported coin types
func GetCoinTypes() []string {
	return []string{CoinTypeBTC, CoinTypeETH, CoinTypeSKY, CoinTypeWAVES, CoinTypeWAVESMDL}
}
