// db_test package implements the test case on db.
package db_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"github.com/dongle/go-order-bot/internal/repository/db"
	"github.com/dongle/go-order-bot/internal/types"
)

func Test_AddTradeHistory(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Error(err)
		return
	}

	// configure logger based on the configuration
	log := logger.New(cfg)

	dbBridge, err := db.New(cfg, log)
	if err != nil {
		t.Error(err)
		return
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		pmin := 52000
		pmax := 57000
		randPrice := rand.Intn(pmax-pmin) + pmin

		tmin := 1627966236
		tmax := 1630644636
		randTime := rand.Intn(tmax-tmin) + tmin

		var history = types.TradeHistory{
			UserId:    2,
			Pair:      "BTC/USDT",
			Side:      0,
			Price:     float64(randPrice),
			Excuted:   2,
			Fee:       27,
			Timestamp: uint64(randTime),
		}

		err1 := dbBridge.AddTradeHistory(&history)
		if err1 != nil {
			t.Error(err)
			return
		}
	}

}
