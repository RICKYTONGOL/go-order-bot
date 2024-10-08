// Package cache implements bridge to fast in-memory object cache.
package cache

import (
	"strings"

	"github.com/dongle/go-order-bot/internal/types"
)

const dwHistoryCacheIdPrefix = "dw_"

func dwHistoryCacheKey(txHash string) string {
	var key strings.Builder
	key.WriteString(dwHistoryCacheIdPrefix)
	key.WriteString(txHash)

	return key.String()
}

// PullDWHistory tries to pull dwhistory from the given txhash from internal in-memory cache.
func (b *MemBridge) PullDWHistory(txHash string) *types.DWHistory {
	data, err := b.cache.Get(dwHistoryCacheKey(txHash))
	if err != nil {
		return nil
	}

	history, err := types.UnmarshalDWHistory(data)
	if err != nil {
		b.log.Criticalf("can not unmarshal history of %s; %s", txHash, err.Error())
		return nil
	}
	return history
}

// PushDWHistory stored the given dwhistory in memory cache.
func (b *MemBridge) PushDWHistory(history *types.DWHistory) {
	if history == nil {
		return
	}

	data, err := history.Marshal()
	if err != nil {
		b.log.Criticalf("can not marshal history of %s; %s", history.TxHash, err.Error())
		return
	}

	if err := b.cache.Set(dwHistoryCacheKey(history.TxHash), data); err != nil {
		b.log.Criticalf("can not marshal history of %s; %s", history.TxHash, err.Error())
	}
}
