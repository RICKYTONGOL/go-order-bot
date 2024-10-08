// Package v1 implements API server version 1.0.
package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"github.com/dongle/go-order-bot/internal/repository"
	"github.com/dongle/go-order-bot/internal/types"
	"github.com/rs/cors"
)

func validateParms(r *http.Request) (ret types.TradeList) {
	ret.Success = true

	fsym := r.PostFormValue("fsym")
	if fsym == "" {
		ret.Success = false
		ret.Error.Msg = "fsym is require"
	}

	tsym := r.PostFormValue("tsym")
	if tsym == "" {
		ret.Success = false
		ret.Error.Msg = "tsym is require"
	}

	interval := r.PostFormValue("interval")
	if interval == "" {
		ret.Success = false
		ret.Error.Msg = "interval is require"
	}

	startTime := r.PostFormValue("startTime")
	if startTime == "" {
		ret.Success = false
		ret.Error.Msg = "start time is require"
	}

	endTime := r.PostFormValue("endTime")
	if endTime == "" {
		ret.Success = false
		ret.Error.Msg = "end time is require"
	}

	return ret
}

func buildTradeHistory(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ret := validateParms(r)

			w.Header().Set("Content-Type", "application/json")

			if !ret.Success {
				json.NewEncoder(w).Encode(ret)
				return
			}

			fsym := r.PostFormValue("fsym")
			tsym := r.PostFormValue("tsym")
			interval := r.PostFormValue("interval")
			startTime := r.PostFormValue("startTime")
			endTime := r.PostFormValue("endTime")

			reqUrl := "https://api.binance.com/api/v3/klines?symbol=" + fsym + tsym + "&interval=" + interval +
				"&startTime=" + startTime + "&endTime=" + endTime
			fmt.Println(reqUrl)

			// resp, err := http.Get("https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1h&startTime=1635724800000&endTime=1642204800000")
			resp, err := http.Get(reqUrl)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var binance []types.BinanceTrade
			if err := json.Unmarshal([]byte(body), &binance); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", binance[0].Opentime)
			fmt.Printf("%v\n", binance[0].Volume)

			// {
			// 	1635811200000,				Open time
			// 	"60911.12000000",			Open
			// 	"64270.00000000",			High
			// 	"60624.68000000",			Low
			// 	"63219.99000000",			Close
			// 	"46368.28410000",			Volume
			// 	1635897599999,				Close time
			// 	"2909221038.40065000",		Quote asset volume
			// 	1858362,					Number of trades
			// 	"23852.14282000",			Taker buy base asset volume
			// 	"1496448016.59887120",		Taker buy quote asset volume
			// 	"0"							Ignore
			// },

			coinPair := fsym + "/" + tsym
			// if fsym == "SYS" && tsym == "USDT" {
			// 	coinPair = "ZNX/USDT"
			// }
			if fsym == "SYS" && tsym == "BTC" {
				coinPair = "ZNX/BTC"
			}
			if fsym == "TRIBE" && tsym == "BNB" {
				coinPair = "ZNX/BNB"
			}

			if fsym == "SYS" && tsym == "USDT" {
				coinPair = "ZNX/DAI"
			}
			if fsym == "IOTA" && tsym == "ETH" {
				coinPair = "ZNX/ETH"
			}

			for _, bTrade := range binance {
				openPrice, _ := strconv.ParseFloat(bTrade.Open, 64)
				closePrice, _ := strconv.ParseFloat(bTrade.Close, 64)
				highPrice, _ := strconv.ParseFloat(bTrade.High, 64)
				lowPrice, _ := strconv.ParseFloat(bTrade.Low, 64)

				volume, _ := strconv.ParseFloat(bTrade.Volume, 64)
				amount := volume / float64(bTrade.TradeNumber)

				randomAmount := rand.Float64() * amount
				var history = types.TradeHistory{
					UserId:    1,
					Pair:      coinPair,
					Side:      config.ORDER_BUY,
					Price:     openPrice,
					Excuted:   randomAmount,
					Fee:       10,
					Timestamp: bTrade.Opentime / 1000,
				}

				err1 := repository.R().AddTradeHistory(&history)
				if err1 != nil {
					log.Debug(err)
					continue
				}

				randomAmount = rand.Float64() * amount
				history = types.TradeHistory{
					UserId:    1,
					Pair:      coinPair,
					Side:      config.ORDER_SELL,
					Price:     closePrice,
					Excuted:   randomAmount,
					Fee:       10,
					Timestamp: bTrade.Closetime / 1000,
				}

				err1 = repository.R().AddTradeHistory(&history)
				if err1 != nil {
					log.Debug(err)
					continue
				}

				tmin := int(bTrade.Opentime)
				tmax := int(bTrade.Closetime)

				randomAmount = rand.Float64() * amount
				randTime := rand.Intn(tmax-tmin) + tmin
				history = types.TradeHistory{
					UserId:    1,
					Pair:      coinPair,
					Side:      config.ORDER_BUY,
					Price:     highPrice,
					Excuted:   randomAmount,
					Fee:       10,
					Timestamp: uint64(randTime / 1000),
				}
				err1 = repository.R().AddTradeHistory(&history)
				if err1 != nil {
					log.Debug(err)
					continue
				}

				randomAmount = rand.Float64() * amount
				randTime = rand.Intn(tmax-tmin) + tmin
				history = types.TradeHistory{
					UserId:    1,
					Pair:      coinPair,
					Side:      config.ORDER_SELL,
					Price:     lowPrice,
					Excuted:   randomAmount,
					Fee:       10,
					Timestamp: uint64(randTime / 1000),
				}
				err1 = repository.R().AddTradeHistory(&history)
				if err1 != nil {
					log.Debug(err)
					continue
				}
			}

			totalCnt := len(binance)
			ret.TotalCnt = uint64(totalCnt)

			log.Debugf("Build Trade history success")

			json.NewEncoder(w).Encode(ret)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func emptyTradeHistory(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var ret types.TradeList
			ret.Success = true

			w.Header().Set("Content-Type", "application/json")

			delCnt, err := repository.R().EmptyTradeHistory()
			if err != nil {
				ret.Success = false
				ret.Error.Msg = "Empty Tradehistory error"
				json.NewEncoder(w).Encode(ret)

				return
			}

			ret.TotalCnt = uint64(delCnt)
			json.NewEncoder(w).Encode(ret)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func validateRemove(r *http.Request) (ret types.TradeList) {
	ret.Success = true

	pair := r.PostFormValue("pair")
	if pair == "" {
		ret.Success = false
		ret.Error.Msg = "coin pair is require"
	}

	return ret
}

func removeTradeHistory(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ret := validateRemove(r)

			w.Header().Set("Content-Type", "application/json")

			if !ret.Success {
				json.NewEncoder(w).Encode(ret)
				return
			}

			pair := r.PostFormValue("pair")

			delCnt, err := repository.R().RemoveTradeHistory(pair)
			if err != nil {
				ret.Success = false
				ret.Error.Msg = "Remove Tradehistory error"
				json.NewEncoder(w).Encode(ret)

				return
			}

			ret.TotalCnt = uint64(delCnt)
			json.NewEncoder(w).Encode(ret)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func validateAdd(r *http.Request) (ret types.TradeList) {
	ret.Success = true

	userId := r.PostFormValue("user_id")
	if userId == "" {
		ret.Success = false
		ret.Error.Msg = "User Id is require"
	}

	pair := r.PostFormValue("pair")
	if pair == "" {
		ret.Success = false
		ret.Error.Msg = "coin pair is require"
	}

	side := r.PostFormValue("side")
	if side == "" {
		ret.Success = false
		ret.Error.Msg = "side(buy/sell) is require"
	}

	price := r.PostFormValue("price")
	if price == "" {
		ret.Success = false
		ret.Error.Msg = "Price is require"
	}

	excuted := r.PostFormValue("excuted")
	if excuted == "" {
		ret.Success = false
		ret.Error.Msg = "Excuted amount is require"
	}

	fee := r.PostFormValue("fee")
	if fee == "" {
		ret.Success = false
		ret.Error.Msg = "fee is require"
	}

	timestamp := r.PostFormValue("timestamp")
	if timestamp == "" {
		ret.Success = false
		ret.Error.Msg = "timestamp is require"
	}

	return ret
}

func addTradeHistory(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ret := validateAdd(r)

			w.Header().Set("Content-Type", "application/json")

			if !ret.Success {
				json.NewEncoder(w).Encode(ret)
				return
			}

			userId := r.PostFormValue("user_id")
			iUserId, _ := strconv.ParseUint(userId, 10, 64)

			pair := r.PostFormValue("pair")

			side := r.PostFormValue("side")
			iSide, _ := strconv.ParseUint(side, 10, 64)

			price := r.PostFormValue("price")
			fPrice, _ := strconv.ParseFloat(price, 64)

			excuted := r.PostFormValue("excuted")
			fExcuted, _ := strconv.ParseFloat(excuted, 64)

			fee := r.PostFormValue("fee")
			fFee, _ := strconv.ParseFloat(fee, 64)

			timestamp := r.PostFormValue("timestamp")
			iTimestamp, _ := strconv.ParseUint(timestamp, 10, 64)

			var history = types.TradeHistory{
				UserId:    iUserId,
				Pair:      pair,
				Side:      iSide,
				Price:     fPrice,
				Excuted:   fExcuted,
				Fee:       fFee,
				Timestamp: iTimestamp,
			}

			err := repository.R().AddTradeHistory(&history)
			if err != nil {
				ret.Success = false
				ret.Error.Msg = "Add tradehistory error"
				json.NewEncoder(w).Encode(ret)

				return
			}

			var result []types.TradeHistory
			result = append(result, history)
			ret.Data = result

			json.NewEncoder(w).Encode(ret)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func setupTradeHistoryHandlers(mux *http.ServeMux, log logger.Logger, corsHandler *cors.Cors) {
	// No Auth
	mux.Handle(urlPatternStr("/tradehistory/build"), corsHandler.Handler(buildTradeHistory(log)))
	mux.Handle(urlPatternStr("/tradehistory/empty"), corsHandler.Handler(emptyTradeHistory(log)))
	mux.Handle(urlPatternStr("/tradehistory/remove"), corsHandler.Handler(removeTradeHistory(log)))
	mux.Handle(urlPatternStr("/tradehistory/add"), corsHandler.Handler(addTradeHistory(log)))
}
