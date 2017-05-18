package cmap

import (
	"hash/fnv"
)

// ShardManagerFNV implements manager using `hash/fnv.Hash32` to hash the keys
type ShardManagerFNV struct {
	TotalShards int64
}

// GetTotalShards implements manager
func (sm *ShardManagerFNV) GetTotalShards() int64 {
	if sm.TotalShards == 0 {
		return 32
	}

	return sm.TotalShards
}

// GetShardNo implements manager
func (sm *ShardManagerFNV) GetShardNo(key string) (int64, error) {
	myHash := fnv.New32()
	_, err := myHash.Write([]byte(key))
	if err != nil {
		return 0, err
	}

	sum := myHash.Sum32()
	totalShards := uint32(sm.GetTotalShards())

	return int64(sum % totalShards), nil
}
