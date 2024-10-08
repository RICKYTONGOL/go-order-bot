// utilservice package implements functions that process the matched order list.
package utilservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/repository"
	"github.com/dongle/go-order-bot/internal/types"
	"github.com/dongle/go-order-bot/internal/utils"
)

type OrderBot struct {
	sigStop chan bool
	mgr     *UtilManager
}

func (svr *OrderBot) init() {
	// svr.sendDailyTrade()
}

func (svr *OrderBot) start() {
	go svr.doTradeHistoryService()
	go svr.doMakeBotOrderService()

	svr.mgr.started(svr)
}

func (svr *OrderBot) close() {
	if svr.sigStop != nil {
		svr.sigStop <- true
		svr.mgr.wg.Done()
	}
}

func (svr *OrderBot) doTradeHistoryService() {
	for {
		svr.buildDailyTrade()
	}
}

func (svr *OrderBot) buildDailyTrade() {
	time.Sleep(1 * time.Second)

	curTime := time.Now().Unix()
	if curTime%config.ONE_DAY_IN_SEC != 0 {
		return
	}

	svr.generateTradeHistory("BTC/USDT")
	svr.generateTradeHistory("ETH/USDT")
	svr.generateTradeHistory("BNB/USDT")
	svr.generateTradeHistory("ZNX/USDT")

	svr.generateTradeHistory("BNB/ETH")
	svr.generateTradeHistory("ZNX/ETH")

	svr.generateTradeHistory("ZNX/BTC")
	svr.generateTradeHistory("ZNX/BNB")
	svr.generateTradeHistory("ZNX/DAI")
}

func (svr *OrderBot) generateTradeHistory(coinPair string) {
	fsym, tsym := utils.SplitPairString(coinPair)
	interval := "1h"

	curTime := time.Now().Unix()
	startTime := fmt.Sprint((curTime - config.ONE_DAY_IN_SEC) * 1000)
	endTime := fmt.Sprint(curTime * 1000)

	if coinPair == "ZNX/USDT" {
		fsym = "FTM"
		tsym = "USDT"
	}
	if coinPair == "ZNX/BTC" {
		fsym = "FTM"
		tsym = "BTC"
	}
	if coinPair == "ZNX/BNB" {
		fsym = "FTM"
		tsym = "BNB"
	}
	if coinPair == "ZNX/DAI" {
		fsym = "FTM"
		tsym = "USDT"
	}
	if coinPair == "ZNX/ETH" {
		fsym = "FTM"
		tsym = "ETH"
	}

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

	for _, bTrade := range binance {
		openPrice, _ := strconv.ParseFloat(bTrade.Open, 64)
		closePrice, _ := strconv.ParseFloat(bTrade.Close, 64)
		highPrice, _ := strconv.ParseFloat(bTrade.High, 64)
		lowPrice, _ := strconv.ParseFloat(bTrade.Low, 64)

		volume, _ := strconv.ParseFloat(bTrade.Volume, 64)
		amount := volume / float64(bTrade.TradeNumber)

		randomAmount := rand.Float64() * amount
		var history = types.TradeHistory{
			UserId:    config.BOT_USER_ID,
			Pair:      coinPair,
			Side:      config.ORDER_BUY,
			Price:     openPrice,
			Excuted:   randomAmount,
			Fee:       config.BOT_TRADE_FEE,
			Timestamp: bTrade.Opentime / 1000,
		}

		err1 := repository.R().AddTradeHistory(&history)
		if err1 != nil {
			log.Debug(err)
			continue
		}

		randomAmount = rand.Float64() * amount
		history = types.TradeHistory{
			UserId:    config.BOT_USER_ID,
			Pair:      coinPair,
			Side:      config.ORDER_SELL,
			Price:     closePrice,
			Excuted:   randomAmount,
			Fee:       config.BOT_TRADE_FEE,
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
			UserId:    config.BOT_USER_ID,
			Pair:      coinPair,
			Side:      config.ORDER_BUY,
			Price:     highPrice,
			Excuted:   randomAmount,
			Fee:       config.BOT_TRADE_FEE,
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
			UserId:    config.BOT_USER_ID,
			Pair:      coinPair,
			Side:      config.ORDER_SELL,
			Price:     lowPrice,
			Excuted:   randomAmount,
			Fee:       config.BOT_TRADE_FEE,
			Timestamp: uint64(randTime / 1000),
		}
		err1 = repository.R().AddTradeHistory(&history)
		if err1 != nil {
			log.Debug(err)
			continue
		}
	}

	totalCnt := len(binance)

	fmt.Printf("Generate and store trade history - %s : %d \n", coinPair, totalCnt)
}

func (svr *OrderBot) doMakeBotOrderService() {
	for {
		svr.makeBotOrders()
	}
}

func (svr *OrderBot) makeBotOrders() {
	time.Sleep(10 * time.Second)

	svr.createBotOrder("ZNX/USDT")
	svr.createBotOrder("BTC/USDT")
	svr.createBotOrder("ETH/USDT")
	svr.createBotOrder("BNB/USDT")

	svr.createBotOrder("BNB/ETH")
	// svr.createBotOrder("ZNX/ETH")

	svr.createBotOrder("ZNX/BTC")
	// svr.createBotOrder("ZNX/BNB")
	svr.createBotOrder("ZNX/DAI")
}

func (svr *OrderBot) createBotOrder(pair string) {
	fsym, tsym := utils.SplitPairString(pair)
	coinPair := fsym + tsym

	if pair == "ZNX/USDT" {
		coinPair = "FTMUSDT"
	}
	if pair == "ZNX/ETH" {
		coinPair = "FTMETH"
	}
	if pair == "ZNX/BTC" {
		coinPair = "FTMBTC"
	}
	if pair == "ZNX/BNB" {
		coinPair = "FTMBNB"
	}
	if pair == "ZNX/DAI" {
		coinPair = "FTMUSDT"
	}

	// fetch orderbook from binance
	binanceUrl := "https://api.binance.com/api/v3/depth?symbol=" + coinPair + "&limit=2"
	binanceReq, err := http.NewRequest(http.MethodGet, binanceUrl, nil)
	if err != nil {
		log.Error(err)
		return
	}

	clientBsc := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := clientBsc.Do(binanceReq)
	if err != nil {
		log.Error(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("response from server / Status: " + response.Status)
		return
	}

	defer response.Body.Close()

	bodyBsc, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return
	}

	orderBook, err := types.UnmarshalBinanceOrderBook(bodyBsc)
	if err != nil {
		log.Error(err)
		return
	}

	for _, orders := range orderBook.Bids {
		svr.requestBotOrder(config.ORDER_BUY, orders[0], orders[1], pair)
	}

	for _, orders := range orderBook.Asks {
		svr.requestBotOrder(config.ORDER_SELL, orders[0], orders[1], pair)
	}

	fmt.Println("request bot order / Pair: ", pair)
}

func (svr *OrderBot) requestBotOrder(side uint64, price string, amount string, pair string) {
	// request Bot order to dongletrade api
	orderReqUrl := ""
	method := "POST"

	if side == config.ORDER_BUY {
		orderReqUrl = "https://api.dongletrade.com/api/v1/bot/orders/buy"
	}
	if side == config.ORDER_SELL {
		orderReqUrl = "https://api.dongletrade.com/api/v1/bot/orders/sell"
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("price", price)
	writer.WriteField("amount", amount)
	writer.WriteField("user_id", "1")
	writer.WriteField("pair", pair)

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	botReq, err := http.NewRequest(method, orderReqUrl, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	botReq.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(botReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}
