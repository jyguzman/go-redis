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
		store.Store.Set(lp.key, &store.RedisList{Size: 0})
	}

	_, ok := store.Store.Get(lp.key).(*store.RedisList)
	if !ok {
		return protocol.Err("Value is not a list")
	}

	numItems := store.Store.Get(lp.key).(*store.RedisList).Prepend(lp.item)
	return protocol.Integer{Val: numItems}.Serialize()
}

func Lpush(args ...string) string {
	return NewLPush(args[0], args[1]).Execute()
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
		store.Store.Set(rp.key, &store.RedisList{Size: 0})
	}

	_, ok := store.Store.Get(rp.key).(*store.RedisList)
	if !ok {
		return protocol.Err("Value is not a list")
	}

	numItems := store.Store.Get(rp.key).(*store.RedisList).Append(rp.item)
	return protocol.Integer{Val: numItems}.Serialize()
}

func Rpush(args ...string) string {
	return NewRPush(args[0], args[1]).Execute()
}

type LRange struct {
	Args []string
}

func NewLRange(args ...string) *LRange {
	return &LRange{args}
}

//func (lr *LRange) Execute() string {
//	if lr.Args == nil || len(lr.Args) < 3 {
//		return protocol.Err("Not enough arguments (need key, start, stop)")
//	}
//	key, start, stop := lr.Args[0], lr.Args[1], lr.Args[2]
//	start, err := strconv.Atoi(start)
//	stop, errTwo := strconv.Atoi(stop)
//	if err != nil || errTwo != nil {
//		return protocol.Err("Could not start or stop")
//	}
//
//}

//type LPop struct {
//	key string
//}
//
//func NewLPop(key string) *LPop {
//	return &LPop{key}
//}
//
//func (lp *LPop) Execute() string {
//	if !store.Store.Contains(lp.key) {
//
//	}
//}

func init() {
	CommandRegistry["lpush"] = Lpush
	CommandRegistry["rpush"] = Rpush
}
