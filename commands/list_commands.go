package commands

import (
	"go-redis/protocol"
	"go-redis/store"
)

type LPush struct {
	key  string
	item store.RedisString
}

func NewLPush(key string, item string) *LPush {
	return &LPush{key: key, item: store.RedisString{Value: item}}
}

func (lp *LPush) Execute() string {
	if !store.Store.Contains(lp.key) {
		store.Store.Set(lp.key, &store.RedisList{Values: []store.RedisString{}})
	}
	val, ok := store.Store.Get(lp.key).(*store.RedisList)
	if !ok {
		return protocol.Err("Value is not a list")
	}
	numItems := len(val.Values)
	store.Store.Get(lp.key).(*store.RedisList).Prepend(lp.item)
	return protocol.Integer{Val: numItems + 1}.Serialize()
}

func Lpush(key string, item string) string {
	return NewLPush(key, item).Execute()
}

type RPush struct {
	key  string
	item store.RedisString
}

func NewRPush(key string, item string) *RPush {
	return &RPush{key: key, item: store.RedisString{Value: item}}
}

func (rp *RPush) Execute() string {
	if !store.Store.Contains(rp.key) {
		store.Store.Set(rp.key, &store.RedisList{Values: []store.RedisString{}})
	}
	val, ok := store.Store.Get(rp.key).(*store.RedisList)
	if !ok {
		return protocol.Err("Value is not a list")
	}
	numItems := len(val.Values)
	store.Store.Get(rp.key).(*store.RedisList).Push(rp.item)
	return protocol.Integer{Val: numItems + 1}.Serialize()
}

func Rpush(key string, item string) string {
	return NewRPush(key, item).Execute()
}
