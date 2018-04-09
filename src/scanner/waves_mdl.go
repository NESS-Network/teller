// Package scanner scans waves blockchain and check all transactions
// to see if there are addresses in vout that can match our deposit addresses.
// If found, then generate an event and push to deposit event channel
//
package scanner

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/modeneis/waves-go-client/model"
)

// WAVESMDLScanner blockchain scanner to check if there're deposit coins
type WAVESMDLScanner struct {
	log            logrus.FieldLogger
	Base           CommonScanner
	wavesRPCClient WavesRPCClient
}

// NewWavesMDLcoinScanner creates scanner instance
func NewWavesMDLcoinScanner(log logrus.FieldLogger, store Storer, client WavesRPCClient, cfg Config) (*WAVESMDLScanner, error) {
	bs := NewBaseScanner(store, log.WithField("prefix", "scanner.wavesMDL"), CoinTypeWAVESMDL, cfg)

	return &WAVESMDLScanner{
		wavesRPCClient: client,
		log:            log.WithField("prefix", "scanner.wavesMDL"),
		Base:           bs,
	}, nil
}

// Run starts the scanner
func (s *WAVESMDLScanner) Run() error {
	return s.Base.Run(s.GetBlockCount, s.getBlockAtHeight, s.waitForNextBlock, s.scanBlock)
}

// Shutdown shutdown the scanner
func (s *WAVESMDLScanner) Shutdown() {
	s.log.Info("Closing WAVESMDL scanner")
	s.wavesRPCClient.Shutdown()
	s.Base.Shutdown()
	s.log.Info("Waiting for WAVESMDL scanner to stop")
	s.log.Info("WAVESMDL scanner stopped")
}

// scanBlock scans for a new WAVESMDL block every ScanPeriod.
// When a new block is found, it compares the block against our scanning
// deposit addresses. If a matching deposit is found, it saves it to the DB.
func (s *WAVESMDLScanner) scanBlock(block *CommonBlock) (int, error) {
	log := s.log.WithField("hash", block.Hash)
	log = log.WithField("height", block.Height)

	log.Debug("WAVESMDLScanner, Scanning block")

	dvs, err := s.Base.GetStorer().ScanBlock(block, CoinTypeWAVESMDL)
	if err != nil {
		log.WithError(err).Error("WAVESMDLScanner.ScanBlock failed")
		return 0, err
	}

	log = log.WithField("scannedDeposits", len(dvs))
	log.Infof("WAVESMDLScanner, Counted %d deposits from block", len(dvs))

	n := 0
	for _, dv := range dvs {
		select {
		case s.Base.GetScannedDepositChan() <- dv:
			n++
		case <-s.Base.GetQuitChan():
			return n, errQuit
		}
	}

	return n, nil
}

// wavesMDLBlock2CommonBlock convert wavescoin block to common block
func wavesMDLBlock2CommonBlock(block *model.Blocks) (*CommonBlock, error) {
	if block == nil {
		return nil, ErrEmptyBlock
	}
	cb := CommonBlock{}

	cb.Hash = block.Signature
	cb.Height = block.Height
	//cb.RawTx = make([]CommonTx, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		if tx.Recipient == "" {
			continue
		}
		cbTx := CommonTx{}
		cbTx.Txid = tx.ID
		//cbTx.Vout = make([]CommonVout, 0, 1)

		cv := CommonVout{}
		cv.N = uint32(0)
		cv.Value = tx.Amount
		cv.Addresses = []string{tx.Recipient}

		cbTx.Vout = append(cbTx.Vout, cv)
		cb.RawTx = append(cb.RawTx, cbTx)
	}

	return &cb, nil
}

// GetBlockCount returns the hash and height of the block in the longest (best) chain.
func (s *WAVESMDLScanner) GetBlockCount() (int64, error) {
	rb, err := s.wavesRPCClient.GetLastBlocks()
	if err != nil {
		return 0, err
	}

	return rb.Height, nil
}

// getBlock returns block of given hash
func (s *WAVESMDLScanner) getBlock(seq int64) (*CommonBlock, error) {
	rb, err := s.wavesRPCClient.GetBlocksBySeq(seq)
	if err != nil {
		return nil, err
	}

	return wavesMDLBlock2CommonBlock(rb)
}

// getBlockAtHeight returns that block at a specific height
func (s *WAVESMDLScanner) getBlockAtHeight(seq int64) (*CommonBlock, error) {
	b, err := s.getBlock(seq)
	return b, err
}

// getNextBlock returns the next block of given hash, return nil if next block does not exist
func (s *WAVESMDLScanner) getNextBlock(seq int64) (*CommonBlock, error) {
	b, err := s.wavesRPCClient.GetBlocksBySeq(seq + 1)
	if err != nil {
		return nil, err
	}
	return wavesMDLBlock2CommonBlock(b)
}

// waitForNextBlock scans for the next block until it is available
func (s *WAVESMDLScanner) waitForNextBlock(block *CommonBlock) (*CommonBlock, error) {
	log := s.log.WithField("blockHash", block.Hash)
	log = log.WithField("blockHeight", block.Height)
	log.Debug("Waiting for the next block")

	for {
		nextBlock, err := s.getNextBlock(block.Height)
		if err != nil {
			if err == ErrEmptyBlock {
				log.WithError(err).Debug("getNextBlock empty")
			} else {
				log.WithError(err).Error("getNextBlock failed")
			}
		}
		if nextBlock == nil {
			log.Debug("No new block yet")
		}
		if err != nil || nextBlock == nil {
			select {
			case <-s.Base.GetQuitChan():
				return nil, errQuit
			case <-time.After(s.Base.GetScanPeriod()):
				continue
			}
		}

		log.WithFields(logrus.Fields{
			"hash":   nextBlock.Hash,
			"height": nextBlock.Height,
		}).Debug("Found nextBlock")

		return nextBlock, nil
	}
}

// AddScanAddress adds new scan address
func (s *WAVESMDLScanner) AddScanAddress(addr, coinType string) error {
	return s.Base.GetStorer().AddScanAddress(addr, coinType)
}

// GetScanAddresses returns the deposit addresses that need to scan
func (s *WAVESMDLScanner) GetScanAddresses() ([]string, error) {
	return s.Base.GetStorer().GetScanAddresses(CoinTypeWAVESMDL)
}

// GetDeposit returns deposit value channel.
func (s *WAVESMDLScanner) GetDeposit() <-chan DepositNote {
	return s.Base.GetDeposit()
}
