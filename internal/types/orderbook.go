package types

import (
	"encoding/json"
)

type Order struct {
	Price  string `bson:"price"`
	Amount string `bson:"amount"`
}

type BinanceOrderBook struct {
	LastUpdateId uint64     `bson:"lastupdateid"`
	Bids         [][]string `bson:"bids"`
	Asks         [][]string `bson:"asks"`
}

func UnmarshalBinanceOrderBook(data []byte) (*BinanceOrderBook, error) {
	var orderBook BinanceOrderBook
	err := json.Unmarshal(data, &orderBook)
	return &orderBook, err
}
