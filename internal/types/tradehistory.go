// Package types implements different core types of the API.
package types

import (
	"encoding/json"
	"fmt"
)

// TradeHistory is the tradehistory base row
type TradeHistory struct {
	OrderId   interface{} `bson:"orderid"`
	UserId    uint64      `bson:"userid"`
	Pair      string      `bson:"pair"`
	Side      uint64      `bson:"side"`
	Price     float64     `bson:"price"`
	Excuted   float64     `bson:"excuted"`
	Fee       float64     `bson:"fee"`
	Timestamp uint64      `bson:"timestamp"`
}

type TradeList struct {
	Success  bool               `bson:"success"`
	TotalCnt uint64             `bson:"totalcnt"`
	Data     []TradeHistory     `bson:"data"`
	Error    TradeResponseError `bson:"error"`
}

type TradeResponseError struct {
	UnknownError uint64 `bson:"unknown"`
	Msg          string `bson:"msg"`
}

type BinanceTrade struct {
	Opentime         uint64 `bson:"opentime"`
	Open             string `bson:"open"`
	High             string `bson:"high"`
	Low              string `bson:"low"`
	Close            string `bson:"close"`
	Volume           string `bson:"volume"`
	Closetime        uint64 `bson:"closetime"`
	QuoteVolume      string `bson:"quotevolume"`
	TradeNumber      uint64 `bson:"tradenumber"`
	TakerBaseVolume  string `bson:"takerbasevolume"`
	TakerQuoteVolume string `bson:"takerquotevolume"`
	Ignore           string `bson:"ignore"`
}

// UnmarshalTradeHistory parses the JSON-encoded tradehistory data.
func UnmarshalTradeHistory(data []byte) (*TradeHistory, error) {
	var history TradeHistory
	err := json.Unmarshal(data, &history)
	return &history, err
}

// Marshal returns the JSON encoding of tradehistory.
func (acc *TradeHistory) Marshal() ([]byte, error) {
	return json.Marshal(acc)
}

func (n *BinanceTrade) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Opentime, &n.Open, &n.High, &n.Low, &n.Close, &n.Volume,
		&n.Closetime, &n.QuoteVolume, &n.TradeNumber, &n.TakerBaseVolume, &n.TakerQuoteVolume, &n.Ignore}

	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}
