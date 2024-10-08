/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Mongo database for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"github.com/dongle/go-order-bot/internal/types"
)

// Repository interface defines functions the underlying implementation provides to API server.
type Repository interface {

	// TradeHistory
	AddTradeHistory(history *types.TradeHistory) error
	AddMultipleTradeHistory(history []types.TradeHistory) error
	EmptyTradeHistory() (int64, error)
	RemoveTradeHistory(pair string) (int64, error)

	// Close and cleanup the repository.
	Close()
}
