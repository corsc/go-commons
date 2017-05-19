// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmap

import (
	"sync"
)

// New returns a initialized concurrent map
func New(manager ...ShardManager) *Map {
	out := &Map{}

	if len(manager) > 0 {
		out.manager = manager[0]
	} else {
		out.manager = &ShardManagerFNV{}
	}

	totalShards := out.manager.GetTotalShards()

	out.shards = make([]*mapShard, totalShards)
	for shardNo := int64(0); shardNo < totalShards; shardNo++ {
		out.shards[shardNo] = &mapShard{
			items: make(map[string]interface{}),
		}
	}

	return out
}

// Map is a concurrent map
type Map struct {
	// controls how many shards exist and how keys are hashed
	manager ShardManager

	shards []*mapShard
}

type mapShard struct {
	sync.RWMutex
	items map[string]interface{}
}

// ShardManager controls how many shards exist and how keys are hashed
type ShardManager interface {
	// Return the total number of shards in this concurrent map
	GetTotalShards() int64

	// Return the shard number for the supplied key
	GetShardNo(key string) (int64, error)
}

// Tuple is 1 key/value pair from the map
type Tuple struct {
	Key   string
	Value interface{}
}

// Get will attempt to return the requested key or an error.
// `ErrNoSuchItem` indicates the item does not exist
func (c *Map) Get(key string) (interface{}, error) {
	shard, err := c.getShard(key)
	if err != nil {
		return nil, err
	}

	shard.RLock()
	defer shard.RUnlock()

	val, found := shard.items[key]
	if !found {
		return nil, ErrNoSuchItem
	}
	return val, nil
}

// GetElseSet will return the existing value in the map or will set the value using `newValue`.
// Regardless, this method will return the map item value or an error.
func (c *Map) GetElseSet(key string, newValue interface{}) (interface{}, error) {
	shard, err := c.getShard(key)
	if err != nil {
		return nil, err
	}

	shard.Lock()
	defer shard.Unlock()

	val, found := shard.items[key]
	if !found {
		shard.items[key] = newValue
		return newValue, nil
	}
	return val, nil
}

// Set will set the supplied value into the map
func (c *Map) Set(key string, newValue interface{}) error {
	shard, err := c.getShard(key)
	if err != nil {
		return err
	}

	shard.Lock()
	defer shard.Unlock()

	shard.items[key] = newValue
	return nil
}

// Count will return the total number of items in the map
func (c *Map) Count() int64 {
	total := int64(0)

	for _, thisShard := range c.shards {
		thisShard.RLock()
		total += int64(len(thisShard.items))
		thisShard.RUnlock()
	}

	return total
}

// Has will return true if the key exists in the map or false
//
// Note: this method will silently fail on errors
func (c *Map) Has(key string) bool {
	_, err := c.Get(key)
	return err == nil
}

// Remove will remove the key from the map (if exists)
//
// Note: this method will silently fail on errors
func (c *Map) Remove(key string) {
	shard, err := c.getShard(key)
	if err != nil {
		return
	}

	shard.Lock()
	defer shard.Unlock()

	delete(shard.items, key)
}

// Iterator will return a iterator (snapshot) of the map
func (c *Map) Iterator() chan Tuple {
	outputCh := make(chan Tuple)

	go func() {
		defer close(outputCh)

		for _, thisShard := range c.shards {
			thisShard.RLock()
			for key, value := range thisShard.items {
				outputCh <- Tuple{
					Key:   key,
					Value: value,
				}
			}
			thisShard.RUnlock()
		}
	}()

	return outputCh
}

// return the map shard that contains the supplied key
func (c *Map) getShard(key string) (*mapShard, error) {
	shardNo, err := c.manager.GetShardNo(key)
	if err != nil {
		return nil, err
	}

	return c.shards[shardNo], nil
}
