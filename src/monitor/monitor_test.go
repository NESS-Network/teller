package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MDLlife/teller/src/exchange"
	"github.com/MDLlife/teller/src/scanner"
	"github.com/MDLlife/teller/src/util/testutil"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type dummyBtcAddrMgr struct {
	Num uint64
}
type dummyEthAddrMgr struct {
	Num uint64
}

type dummySkyAddrMgr struct {
	Num uint64
}

type dummyWavesAddrMgr struct {
	Num uint64
}

type dummyWavesMDLAddrMgr struct {
	Num uint64
}

func (db *dummyBtcAddrMgr) Remaining() uint64 {
	return db.Num
}
func (db *dummyEthAddrMgr) Remaining() uint64 {
	return db.Num
}
func (db *dummySkyAddrMgr) Remaining() uint64 {
	return db.Num
}
func (db *dummyWavesAddrMgr) Remaining() uint64 {
	return db.Num
}
func (db *dummyWavesMDLAddrMgr) Remaining() uint64 {
	return db.Num
}

type dummyDepositStatusGetter struct {
	dpis []exchange.DepositInfo
}

func (dps dummyDepositStatusGetter) GetDepositStatusDetail(flt exchange.DepositFilter) ([]exchange.DepositStatusDetail, error) {
	var ds []exchange.DepositStatusDetail
	for _, dpi := range dps.dpis {
		if flt(dpi) {
			ds = append(ds, exchange.DepositStatusDetail{
				Seq:            dpi.Seq,
				DepositAddress: dpi.DepositAddress,
				MDLAddress:     dpi.MDLAddress,
				Status:         dpi.Status.String(),
				UpdatedAt:      dpi.UpdatedAt,
				Txid:           dpi.Txid,
				CoinType:       dpi.CoinType,
			})
		}
	}
	return ds, nil
}

func (dps dummyDepositStatusGetter) GetDepositStats() (*exchange.DepositStats, error) {
	stats := &exchange.DepositStats{
		TotalBTCReceived:      0,
		TotalETHReceived:      0,
		TotalSKYReceived:      0,
		TotalWAVESReceived:    0,
		TotalWAVESMDLReceived: 0,
		TotalMDLSent:          0,
		TotalTransactions:     0,
	}

	for _, dpi := range dps.dpis {
		if dpi.Status == exchange.StatusDone { // count only processed
			switch dpi.CoinType {
			case scanner.CoinTypeBTC:
				stats.TotalBTCReceived += dpi.DepositValue
			case scanner.CoinTypeETH:
				stats.TotalETHReceived += dpi.DepositValue
			case scanner.CoinTypeSKY:
				stats.TotalSKYReceived += dpi.DepositValue
			case scanner.CoinTypeWAVES:
				stats.TotalWAVESReceived += dpi.DepositValue
			case scanner.CoinTypeWAVESMDL:
				stats.TotalWAVESMDLReceived += dpi.DepositValue
			}
			stats.TotalMDLSent += int64(dpi.MDLSent)
			stats.TotalTransactions++
		}
	}

	return stats, nil
}

type dummyScanAddrs struct {
	// addrs []string
}

func (ds dummyScanAddrs) GetScanAddresses() ([]string, error) {
	return []string{}, nil
}

// data for stats tests
var statsDpis = []exchange.DepositInfo{
	{
		CoinType:     scanner.CoinTypeBTC,
		DepositValue: 10000,
		MDLSent:      100,
		Status:       exchange.StatusWaitConfirm,
	},
	{
		CoinType:     scanner.CoinTypeBTC,
		DepositValue: 100000,
		MDLSent:      200,
		Status:       exchange.StatusDone,
	},
	{
		CoinType:     scanner.CoinTypeETH,
		DepositValue: 200000,
		MDLSent:      300,
		Status:       exchange.StatusDone,
	},
	{
		CoinType:     scanner.CoinTypeSKY,
		DepositValue: 30000,
		MDLSent:      400,
		Status:       exchange.StatusDone,
	},
	{
		CoinType:     scanner.CoinTypeWAVES,
		DepositValue: 4000,
		MDLSent:      500,
		Status:       exchange.StatusDone,
	},
	{
		CoinType:     scanner.CoinTypeWAVESMDL,
		DepositValue: 4000,
		MDLSent:      500,
		Status:       exchange.StatusDone,
	},
}
var statsCfg = Config{
	"localhost:1234",
	10, 11, 12, 13, 14, 15, decimal.NewFromFloat(10.5), 10,
}

func TestRunMonitor(t *testing.T) {
	dpis := []exchange.DepositInfo{
		{
			DepositAddress: "b1",
			MDLAddress:     "s1",
			Status:         exchange.StatusWaitDeposit,
		},
		{
			DepositAddress: "b2",
			MDLAddress:     "s2",
			Status:         exchange.StatusWaitSend,
		},
		{
			DepositAddress: "b3",
			MDLAddress:     "s3",
			Status:         exchange.StatusWaitConfirm,
		},
		{
			DepositAddress: "b4",
			MDLAddress:     "s4",
			Status:         exchange.StatusDone,
		},
		{
			DepositAddress: "b5",
			MDLAddress:     "s6",
			Status:         exchange.StatusDone,
		},
	}

	dummyDps := dummyDepositStatusGetter{dpis: dpis}

	log, _ := testutil.NewLogger(t)
	m := New(log, statsCfg, &dummyBtcAddrMgr{10}, &dummyEthAddrMgr{10}, &dummySkyAddrMgr{10}, &dummyWavesAddrMgr{}, &dummyWavesMDLAddrMgr{}, &dummyDps, &dummyScanAddrs{})

	err := setupTestServer(t, m)
	require.NoError(t, err)

	targetServer := fmt.Sprintf("http://%s/api/address", statsCfg.Addr)
	rsp, err := http.Get(targetServer)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var addrUsage addressUsage
	err = json.NewDecoder(rsp.Body).Decode(&addrUsage)
	require.NoError(t, err)
	require.Equal(t, uint64(10), addrUsage.RestAddrNum)
	testutil.CheckError(t, rsp.Body.Close)

	var tt = []struct {
		name        string
		status      string
		expectCode  int //
		expectValue []exchange.DepositInfo
	}{
		{
			"get deposit that are in waiting_deposit status",
			"waiting_deposit",
			http.StatusOK,
			dpis[:1],
		},
		{
			"get deposit that are in waiting_send status",
			"waiting_send",
			http.StatusOK,
			dpis[1:2],
		},
		{
			"get deposit that are in waiting_confirm status",
			"waiting_confirm",
			http.StatusOK,
			dpis[2:3],
		},
		{
			"get deposit that are in waiting_done status",
			"done",
			http.StatusOK,
			dpis[3:5],
		},
		{
			"get unknown status",
			"invalid",
			http.StatusBadRequest,
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			targetServer := fmt.Sprintf("http://%s/api/deposit_status?status=%s", statsCfg.Addr, tc.status)
			rsp, err := http.Get(targetServer)
			require.NoError(t, err)
			defer testutil.CheckError(t, rsp.Body.Close)
			require.Equal(t, tc.expectCode, rsp.StatusCode)

			if rsp.StatusCode == http.StatusOK {
				var st []exchange.DepositStatusDetail
				err := json.NewDecoder(rsp.Body).Decode(&st)
				require.NoError(t, err)

				dss := make([]exchange.DepositInfo, 0, len(st))
				for _, s := range st {
					dss = append(dss, exchange.DepositInfo{
						Seq:            s.Seq,
						UpdatedAt:      s.UpdatedAt,
						Status:         exchange.NewStatusFromStr(s.Status),
						DepositAddress: s.DepositAddress,
						MDLAddress:     s.MDLAddress,
						Txid:           s.Txid,
					})
				}
				require.Equal(t, tc.expectValue, dss)
			}
		})
	}

	defer func() {
		m.Shutdown()
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
	}()
}

func TestMonitorDepositStats(t *testing.T) {
	dummyDps := dummyDepositStatusGetter{dpis: statsDpis}

	log, _ := testutil.NewLogger(t)
	m := New(log, statsCfg, &dummyBtcAddrMgr{10}, &dummyEthAddrMgr{10}, &dummySkyAddrMgr{10}, &dummyWavesAddrMgr{10}, &dummyWavesMDLAddrMgr{10}, &dummyDps, &dummyScanAddrs{})
	err := setupTestServer(t, m)
	require.NoError(t, err)

	targetServer := fmt.Sprintf("http://%s/api/stats", statsCfg.Addr)
	rsp, err := http.Get(targetServer)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var stats exchange.DepositStats
	err = json.NewDecoder(rsp.Body).Decode(&stats)
	require.NoError(t, err)
	require.Equal(t, statsDpis[1].DepositValue, stats.TotalBTCReceived)
	require.Equal(t, statsDpis[2].DepositValue, stats.TotalETHReceived)
	require.Equal(t, statsDpis[3].DepositValue, stats.TotalSKYReceived)
	require.Equal(t, statsDpis[4].DepositValue, stats.TotalWAVESReceived)
	require.Equal(t, statsDpis[5].DepositValue, stats.TotalWAVESMDLReceived)
	require.Equal(t, int64(5), stats.TotalTransactions)

	mdlTotal := statsDpis[1].MDLSent + statsDpis[2].MDLSent + statsDpis[3].MDLSent + statsDpis[4].MDLSent + statsDpis[5].MDLSent
	require.Equal(t, mdlTotal, uint64(stats.TotalMDLSent))

	defer func() {
		testutil.CheckError(t, rsp.Body.Close)
		m.Shutdown()
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
	}()

}

func TestMonitorWebReadyDepositStats(t *testing.T) {
	dummyDps := dummyDepositStatusGetter{dpis: statsDpis}

	log, _ := testutil.NewLogger(t)
	m := New(log, statsCfg, &dummyBtcAddrMgr{10}, &dummyEthAddrMgr{10}, &dummySkyAddrMgr{10}, &dummyWavesAddrMgr{10}, &dummyWavesMDLAddrMgr{10}, &dummyDps, &dummyScanAddrs{})
	err := setupTestServer(t, m)

	require.NoError(t, err)

	targetServer := fmt.Sprintf("http://%s/api/web-stats", statsCfg.Addr)
	rsp, err := http.Get(targetServer)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var webStats WebReadyStats
	err = json.NewDecoder(rsp.Body).Decode(&webStats)
	require.NoError(t, err)

	require.Equal(t, "0.0010001", webStats.TotalBTCReceived)
	require.Equal(t, "0.000200011", webStats.TotalETHReceived)
	require.Equal(t, "0.030012", webStats.TotalSKYReceived)
	require.Equal(t, "0.00004013", webStats.TotalWAVESReceived)
	require.Equal(t, "0.00004014", webStats.TotalWAVESMDLReceived)
	require.Equal(t, "0.001915", webStats.TotalMDLSent)
	require.Equal(t, "10.50009575", webStats.TotalUSDReceived)
	require.Equal(t, int64(15), webStats.TotalTransactions)

	defer func() {
		testutil.CheckError(t, rsp.Body.Close)
		m.Shutdown()
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
	}()

}

func TestMonitorUpdateEthToUSDCourse(t *testing.T) {
	cryptocompareETHtoUSDcourse = float32(0)

	log, _ := testutil.NewLogger(t)

	updateEthToUSDCourse(log)

	require.NotEqual(t, float32(0), cryptocompareETHtoUSDcourse)
}

// func TestMonitorUpdateAdditionalEth(t *testing.T) {
// 	ethApiValue = decimal.New(0,0)
//
// 	log, _ := testutil.NewLogger(t)
//
// 	updateAdditionalEth = func(log logrus.FieldLogger) {}
// 	updateAdditionalEth(log)
//
// 	log.Print(ethApiValue)
// 	require.False(t, ethApiValue.Equal(decimal.New(0,0)))
// }

func TestMonitorEthTotalStatsHandler(t *testing.T) {
	ethAPIValue = decimal.New(0, 0)
	updateEthToUSDCourse = func(log logrus.FieldLogger) {}
	cryptocompareETHtoUSDcourse = 100

	dummyDps := dummyDepositStatusGetter{dpis: statsDpis}

	log, _ := testutil.NewLogger(t)
	m := New(log, statsCfg, &dummyBtcAddrMgr{10}, &dummyEthAddrMgr{10}, &dummySkyAddrMgr{10}, &dummyWavesAddrMgr{10}, &dummyWavesMDLAddrMgr{10}, &dummyDps, &dummyScanAddrs{})
	err := setupTestServer(t, m)

	require.NoError(t, err)

	targetServer := fmt.Sprintf("http://%s/api/eth-total-stats", statsCfg.Addr)
	rsp, err := http.Get(targetServer)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var j map[string]string
	err = json.NewDecoder(rsp.Body).Decode(&j)
	require.NoError(t, err)

	require.Equal(t, map[string]string{"eth": "0.1050009575"}, j)

	defer func() {
		testutil.CheckError(t, rsp.Body.Close)
		m.Shutdown()
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
	}()
}

func setupTestServer(t *testing.T, m *Monitor) error {
	mux := m.setupMux()

	m.ln = &http.Server{
		Addr:         m.cfg.Addr,
		Handler:      mux,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	go func() {
		if err := m.Run(); err != nil {
			t.Logf("m.Run(), %v", err)
			return
		}
	}()

	timer := time.NewTimer(time.Second * 1)
	<-timer.C

	return nil

}
