package commands

import (
	"go-redis/protocol"
	"go-redis/store"
	"strconv"
)

type SetCommand struct {
	key   string
	value store.RedisString
}

func (sc *SetCommand) Execute() string {
	store.Store.Set(sc.key, sc.value)
	return protocol.Ok()
}

func NewSetCommand(key string, value string) *SetCommand {
	return &SetCommand{key: key, value: store.RedisString{Value: value}}
}

func Set(args []string) string {
	key, value := args[0], args[1]
	return NewSetCommand(key, value).Execute()
}

type GetCommand struct {
	key string
}

func (gc *GetCommand) Execute() string {
	val := store.Store.Get(gc.key)
	if v, ok := val.(store.RedisString); ok {
		return protocol.BulkString{Val: v.Value}.Serialize()
	}
	return protocol.Null()
}

func NewGetCommand(key string) *GetCommand {
	return &GetCommand{key: key}
}

func Get(key string) string {
	return NewGetCommand(key).Execute()
}

type IncrCommand struct {
	key string
}

func (ic *IncrCommand) Execute() string {
	if !store.Store.Contains(ic.key) {
		store.Store.Set(ic.key, store.RedisString{Value: "0"})
	}
	val := store.Store.Get(ic.key)
	rString, ok := val.(store.RedisString)
	if !ok {
		return protocol.Err("Value not of type string")
	}
	num, err := strconv.Atoi(rString.Value)
	if err != nil {
		return protocol.Err("Failed to parse value as integer")
	}
	store.Store.Set(ic.key, store.RedisString{Value: strconv.Itoa(num + 1)})
	return protocol.Ok()
}

func NewIncrCommand(key string) *IncrCommand {
	return &IncrCommand{key: key}
}

func Incr(key string) string {
	return NewIncrCommand(key).Execute()
}

type DecrCommand struct {
	key string
}

func (dc *DecrCommand) Execute() string {
	if !store.Store.Contains(dc.key) {
		store.Store.Set(dc.key, store.RedisString{Value: "0"})
	}
	val := store.Store.Get(dc.key)
	rString, ok := val.(store.RedisString)
	if !ok {
		return protocol.Err("Value not of type string")
	}
	num, err := strconv.Atoi(rString.Value)
	if err != nil {
		return protocol.Err("Failed to parse value as integer.")
	}
	store.Store.Set(dc.key, store.RedisString{Value: strconv.Itoa(num - 1)})
	return protocol.Ok()
}

func NewDecrCommand(key string) *DecrCommand {
	return &DecrCommand{key: key}
}

func Decr(key string) string {
	return NewDecrCommand(key).Execute()
}
