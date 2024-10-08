// Package db implements bridge to persistent storage represented by Mongo database.
package db

import (
	"context"
	"fmt"
	"log"

	"github.com/dongle/go-order-bot/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// coTradeHistory is the name of the off-chain database collection storing trade history.
	coTradeHistory = "tradehistory"

	// tradeOrderIdField is the order id field of TradeHistory collection.
	tradeOrderIdField = "orderid"

	// tradeUserIdField is the user_id field of TradeHistory collection.
	tradeUserIdField = "userid"

	// tradePairField is the pair field of TradeHistory collection.
	tradePairField = "pair"

	// tradeSideField is the side field of TradeHistory collection
	tradeSideField = "side"

	// tradePriceField is the price field of TradeHistory collection
	tradePriceField = "price"

	// tradeExcutedField is the excuted field of TradeHistory collection
	tradeExcutedField = "excuted"

	// tradeFeeField is the to fee field of TradeHistory collection
	tradeFeeField = "fee"

	// tradeTimestampField is the timestamp field of TradeHistory collection
	tradeTimestampField = "timestamp"
)

type TradeHistoryRow struct {
	OrderId   uint64 `bson:"orderid"`
	UserId    uint64 `bson:"userid"`
	Pair      string `bson:"pair"`
	Side      string `bson:"side"`
	Price     string `bson:"price"`
	Excuted   uint64 `bson:"excuted"`
	Fee       string `bson:"fee"`
	Timestamp uint64 `bson:"timestamp"`
}

func (db *MongoDbBridge) initTradeHisotryCollection() {
	db.log.Debugf("tradehistory collection initialized")
}

func (db *MongoDbBridge) AddTradeHistory(history *types.TradeHistory) error {
	if history == nil {
		return fmt.Errorf("can not add empty history")
	}

	// get the collection for tradehistory transactions
	col := db.client.Database(db.dbName).Collection(coTradeHistory)

	_, err := col.InsertOne(context.Background(), bson.D{
		{Key: tradeOrderIdField, Value: history.OrderId},
		{Key: tradeUserIdField, Value: history.UserId},
		{Key: tradePairField, Value: history.Pair},
		{Key: tradeSideField, Value: history.Side},
		{Key: tradePriceField, Value: history.Price},
		{Key: tradeExcutedField, Value: history.Excuted},
		{Key: tradeFeeField, Value: history.Fee},
		{Key: tradeTimestampField, Value: history.Timestamp},
	})

	if err != nil {
		db.log.Error("can not insert new trade history")

		return err
	}

	// check init state
	// make sure transactions collection is initialized
	if db.initTradeHistory != nil {
		db.initTradeHistory.Do(func() { db.initTradeHisotryCollection(); db.initTradeHistory = nil })
	}

	db.log.Debugf("added history :Pair %s", history.Pair)

	return nil
}

// TradeHistoryCount calculates total number of tradehistories in the database.
func (db *MongoDbBridge) TradeHistoryCount() (uint64, error) {
	return db.EstimateCount(db.client.Database(db.dbName).Collection(coTradeHistory))
}

func (db *MongoDbBridge) AddMultipleTradeHistory(history []types.TradeHistory) error {
	if history == nil {
		return fmt.Errorf("can not add empty history")
	}

	// get the collection for tradehistory transactions
	col := db.client.Database(db.dbName).Collection(coTradeHistory)

	for _, trade := range history {
		_, err := col.InsertOne(context.Background(), bson.D{
			{Key: tradeOrderIdField, Value: trade.OrderId},
			{Key: tradeUserIdField, Value: trade.UserId},
			{Key: tradePairField, Value: trade.Pair},
			{Key: tradeSideField, Value: trade.Side},
			{Key: tradePriceField, Value: trade.Price},
			{Key: tradeExcutedField, Value: trade.Excuted},
			{Key: tradeFeeField, Value: trade.Fee},
			{Key: tradeTimestampField, Value: trade.Timestamp},
		})

		if err != nil {
			db.log.Error("can not insert new trade history")
			return err
		}
	}

	// check init state
	// make sure transactions collection is initialized
	if db.initTradeHistory != nil {
		db.initTradeHistory.Do(func() { db.initTradeHisotryCollection(); db.initTradeHistory = nil })
	}

	db.log.Debugf("added multiple tradehistory")

	return nil
}

func (db *MongoDbBridge) EmptyTradeHistory() (int64, error) {
	// get the collection for tradehistory transactions
	col := db.client.Database(db.dbName).Collection(coTradeHistory)
	opts := options.Delete()

	res, err := col.DeleteMany(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	return res.DeletedCount, nil
}

func (db *MongoDbBridge) RemoveTradeHistory(pair string) (int64, error) {
	// get the collection for tradehistory transactions
	col := db.client.Database(db.dbName).Collection(coTradeHistory)
	opts := options.Delete()
	filter := bson.D{
		{Key: tradePairField, Value: pair},
	}

	res, err := col.DeleteMany(context.Background(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	return res.DeletedCount, nil
}
