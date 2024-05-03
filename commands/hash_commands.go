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

func NewHSetCommand(args []string) *HSetCommand {
	return &HSetCommand{args}
}

func Hset(args ...string) (string, error) {
	return NewHSetCommand(args).Execute()
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

func NewHGetCommand(args []string) *HGetCommand {
	return &HGetCommand{args}
}

func HGet(args ...string) (string, error) {
	return NewHGetCommand(args).Execute()
}

type HDelCommand struct {
	args []string
}

func (hs *HDelCommand) Args() []string {
	return hs.args
}

func (hs *HDelCommand) Execute() (string, error) {
	if len(hs.args) < 4 {
		return "", fmt.Errorf("not enough arguments for HDEL")
	}
	key, fields := hs.args[1], hs.args[2:]
	if !store.Store.Contains(key) {
		return protocol.IntegerResponse(0), nil
	}

	_, isHash := store.Store.Get(key).(*data_types.RedisHash)
	if !isHash {
		return "", fmt.Errorf("value not of type hash")
	}

	numDeleted := 0
	for _, field := range fields {
		numDeleted += store.Store.Get(key).(*data_types.RedisHash).Delete(field)
	}

	return protocol.IntegerResponse(numDeleted), nil
}

func NewHDelCommand(args []string) *HDelCommand {
	return &HDelCommand{args}
}

func HDel(args ...string) (string, error) {
	return NewHDelCommand(args).Execute()
}

func init() {
	CommandRegistry["hset"] = Hset
	CommandRegistry["hget"] = HGet
	CommandRegistry["hdel"] = HDel
}
