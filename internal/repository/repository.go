/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Mongo database for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"fmt"
	"sync"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"github.com/dongle/go-order-bot/internal/repository/cache"
	"github.com/dongle/go-order-bot/internal/repository/db"
)

// repo represents an instance of the Repository manager.
var repo Repository

// onceRepo is the sync object used to make sure the Repository
// is instantiated only once on the first demand.
var onceRepo sync.Once

// config represents the configuration setup used by the repository
// to establish and maintain required connectivity to external services
// as needed.
var cfg *config.Config

// log represents the logger to be used by the repository.
var log logger.Logger

// SetConfig sets the repository configuration to be used to establish
// and maintain external repository connections.
func SetConfig(c *config.Config) {
	cfg = c
}

// SetLogger sets the repository logger to be used to collect logging info.
func SetLogger(l logger.Logger) {
	log = l
}

// R provides access to the singleton instance of the Repository.
func R() Repository {
	// make sure to instantiate the Repository only once
	onceRepo.Do(func() {
		repo = newRepository()
	})
	return repo
}

// Proxy represents Repository interface implementation and controls access to data through
// several low level bridges.
type proxy struct {
	log logger.Logger
	cfg *config.Config

	db    *db.MongoDbBridge
	cache *cache.MemBridge
}

func newRepository() Repository {
	if cfg == nil {
		panic(fmt.Errorf("missing configuration"))
	}
	if log == nil {
		panic(fmt.Errorf("missing logger"))
	}

	caBridge, dbBridge, err := connect(cfg, log)
	if err != nil {
		log.Fatal("repository init failed")
		return nil
	}

	p := proxy{
		log:   log,
		cfg:   cfg,
		db:    dbBridge,
		cache: caBridge,
	}

	return &p
}

// connect opens connections to the external sources we need.
func connect(cfg *config.Config, log logger.Logger) (*cache.MemBridge, *db.MongoDbBridge, error) {
	// create new database connection bridge
	dbBridge, err := db.New(cfg, log)
	if err != nil {
		log.Criticalf("can not connect backend persistent storage, %s", err.Error())
		return nil, nil, err
	}

	// create new in-memory cache bridge
	caBridge, err := cache.New(cfg, log)
	if err != nil {
		log.Criticalf("can not create in-memory cache bridge, %s", err.Error())
		return nil, nil, err
	}

	return caBridge, dbBridge, nil
}

// Close with close all connections and clean up the pending work for gracefull termination.
func (p *proxy) Close() {
	p.log.Notice("repository is closing")

	p.db.Close()

	p.log.Notice("repositoy done")
}
