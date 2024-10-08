// Package db implements bridge to persistent storage represented by Mongo database.
package db

import (
	"context"
	"sync"
	"time"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDbBridge represents Mongo DB abstraction layer.
type MongoDbBridge struct {
	client *mongo.Client
	log    logger.Logger
	dbName string

	// Mutex
	updateBalanceMutex *sync.RWMutex
	userSignUpMutex    *sync.RWMutex
	storeOrderMutex    *sync.RWMutex

	// init state marks
	initTradeHistory *sync.Once
}

// New creates a new Mongo Db connection bridge.
func New(cfg *config.Config, log logger.Logger) (*MongoDbBridge, error) {
	// log what we do
	log.Debugf("connecting database at %s/%s", cfg.Db.Url, cfg.Db.DbName)

	// open the database connection
	con, err := connectDb(&cfg.Db)
	if err != nil {
		log.Criticalf("can not contact the database; %s", err.Error())
		return nil, err
	}

	// log the event
	log.Notice("database connection established")

	// return the bridge
	db := &MongoDbBridge{
		client:             con,
		log:                log,
		dbName:             cfg.Db.DbName,
		updateBalanceMutex: &sync.RWMutex{},
		userSignUpMutex:    &sync.RWMutex{},
		storeOrderMutex:    &sync.RWMutex{},
	}

	// check the state
	db.CheckDatabaseInitState()
	return db, nil
}

// connectDb opens Mongo database connection
func connectDb(cfg *config.Database) (*mongo.Client, error) {
	// get empty unrestricted context
	ctx := context.Background()

	// create new Mongo client
	credential := options.Credential{
		Username:   "dongleapi",
		Password:   "dongleapi123@",
		AuthSource: "dongleapi",
	}
	clientOpts := options.Client().ApplyURI(cfg.Url).SetAuth(credential)

	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Url))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// validate the connection was indeed established
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CheckDatabaseInitState verifies if database collections have been
// already initialized and marks the empty collections so they can be properly
// configured when created.
func (db *MongoDbBridge) CheckDatabaseInitState() {
	// log what we do
	db.log.Debugf("checking database init state")

	db.collectionNeedInit("tradehistory", db.TradeHistoryCount, &db.initTradeHistory)
}

// checkAccountCollectionState checks the Accounts collection state.
func (db *MongoDbBridge) collectionNeedInit(name string, counter func() (uint64, error), init **sync.Once) {
	// use the counter to get the collection size
	count, err := counter()
	if err != nil {
		db.log.Errorf("can not check %s count; %s", name, err.Error())
		return
	}

	// collection not empty,
	if count != 0 {
		db.log.Debugf("found %d %s", count, name)
		return
	}

	// collection init needed, create the init control
	db.log.Noticef("%s collection empty", name)
	var once sync.Once
	*init = &once
}

// EstimateCount calculates an estimated number of documents in the given collection.
func (db *MongoDbBridge) EstimateCount(col *mongo.Collection) (uint64, error) {
	// do the counting
	val, err := col.EstimatedDocumentCount(context.Background())
	if err != nil {
		db.log.Errorf("can not count documents in rewards collection; %s", err.Error())
		return 0, err
	}
	return uint64(val), nil
}

// Close will terminate or finish all operations and close the connection to Mongo database.
func (db *MongoDbBridge) Close() {
	if db.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// try to disconnect
		err := db.client.Disconnect(ctx)
		if err != nil {
			db.log.Errorf("error on closing database connection; %s", err.Error())
		}

		// inform
		db.log.Info("database connection is closed")
		cancel()
	}
}
