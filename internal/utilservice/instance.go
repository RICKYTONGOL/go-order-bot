// utilservice package implements functions that process the matched order list.
package utilservice

import (
	"sync"

	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
)

// onceManager is the sync object used to make sure the Blockchain
// is instantiated only once on the first demand.
var onceManager sync.Once

var manager *UtilManager

// config represents the configuration setup used by the blockchain
// to establish and maintain required connectivity to external services
// as needed.
var cfg *config.Config

// log represents the logger to be used by the blockchain.
var log logger.Logger

// SetConfig sets the blockchain configuration to be used to establish
// and maintain external blockchain connections.
func SetConfig(c *config.Config) {
	cfg = c
}

// SetLogger sets the blockchain logger to be used to collect logging info.
func SetLogger(l logger.Logger) {
	log = l
}

// Instance provides access to the singleton instance of the Blockchain.
func Instance() *UtilManager {
	// make sure to instantiate the Blockchain only once
	onceManager.Do(func() {
		manager = newUtilManager()
	})
	return manager
}
