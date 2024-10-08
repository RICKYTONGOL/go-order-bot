/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Mongo database for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import "github.com/dongle/go-order-bot/internal/types"

// AddTradeHistory adds specified tradehistory detail into the repository.
func (p *proxy) AddTradeHistory(history *types.TradeHistory) error {
	// add this tradehistory to the database and remember it's been added
	err := p.db.AddTradeHistory(history)

	// if err == nil {
	// 	p.cache.PushTradeHistory(history)
	// }
	return err
}

// AddTradeHistory adds specified tradehistory detail into the repository.
func (p *proxy) AddMultipleTradeHistory(history []types.TradeHistory) error {
	// add this tradehistory to the database and remember it's been added
	err := p.db.AddMultipleTradeHistory(history)

	// if err == nil {
	// 	p.cache.PushTradeHistory(history)
	// }
	return err
}

// EmptyTradeHistory remove all trade history.
func (p *proxy) EmptyTradeHistory() (int64, error) {
	cnt, err := p.db.EmptyTradeHistory()

	return cnt, err
}

// RemoveTradeHistory remove trade history by pair.
func (p *proxy) RemoveTradeHistory(pair string) (int64, error) {
	cnt, err := p.db.RemoveTradeHistory(pair)

	return cnt, err
}
