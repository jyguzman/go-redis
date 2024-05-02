package commands

import (
	"fmt"
	"go-redis/data_types"
	"go-redis/protocol"
	"go-redis/store"
)

type HSetCommand struct {
	args []string
}

func (hs *HSetCommand) Args() []string {
	return hs.args
}

func (hs *HSetCommand) Execute() (string, error) {
	if len(hs.args) < 4 {
		return "", fmt.Errorf("not enough arguments for HSET")
	}

	if len(hs.args[2:])%2 != 0 {
		return "", fmt.Errorf("must provide key-value pairs")
	}

	key := hs.args[1]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisHash())
	}

	_, ok := store.Store.Get(key).(*data_types.RedisHash)
	if !ok {
		return "", fmt.Errorf("value not of type hash")
	}

	numPairs := 0
	for i := 2; i < len(hs.args); i += 2 {
		key, value := hs.args[i], hs.args[i+1]
		store.Store.Get(key).(*data_types.RedisHash).Set(key, data_types.NewRedisString(value))
		numPairs += 1
	}

	return protocol.IntegerResponse(numPairs), nil
}

type HGetCommand struct {
	args []string
}

func (hs *HGetCommand) Args() []string {
	return hs.args
}

func (hs *HGetCommand) Execute() (string, error) {
	if len(hs.args) < 3 {
		return "", fmt.Errorf("not enough arguments for HGET")
	}

	key := hs.args[0]
	if !store.Store.Contains(key) {
		return protocol.NilString(), nil
	}

	rHash, isHash := store.Store.Get(key).(*data_types.RedisHash)
	if !isHash {
		return "", fmt.Errorf("value not of type hash")
	}

	rString, ok := rHash.Get(key)
	if !ok {
		return protocol.NilString(), nil
	}

	return protocol.BulkStringResponse(rString.(data_types.RedisString).Value), nil
}
