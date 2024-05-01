package commands

import (
	"go-redis/protocol"
	"go-redis/store"
	"strconv"
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

func (lr *LRange) Execute() string {
	if lr.Args == nil || len(lr.Args) < 3 {
		return protocol.Err("Not enough arguments for LRANGE (need key, start, stop)")
	}

	key, startStr, stopStr := lr.Args[0], lr.Args[1], lr.Args[2]
	start, err := strconv.Atoi(startStr)
	stop, errTwo := strconv.Atoi(stopStr)
	if err != nil || errTwo != nil {
		return protocol.Err("Could not parse start or stop as integers")
	}

	if !store.Store.Contains(key) {
		return protocol.NewArray([]protocol.RespValue{}).Serialize()
	}

	l := store.Store.Get(key).(*store.RedisList)
	ptr := l.Head
	for i := 0; i < start; i++ {
		ptr = ptr.Next
	}
	array := protocol.Array{Val: []protocol.RespValue{}}
	for i := 0; ptr != nil && i < (stop-start)+1; i++ {
		array.Add(protocol.BulkString{Val: ptr.RedisString.Value})
		ptr = ptr.Next
	}
	return array.Serialize()
}

func Lrange(args ...string) string {
	return NewLRange(args...).Execute()
}

type LPop struct {
	args []string
}

func NewLPop(args ...string) *LPop {
	return &LPop{args}
}

func (lp *LPop) Execute() string {
	key := lp.args[0]
	if !store.Store.Contains(key) {
		return protocol.NewArray([]protocol.RespValue{}).Serialize()
	}

	val := store.Store.Get(key).(*store.RedisList).PopLeft()
	return protocol.BulkString{Val: val}.Serialize()
}

func Lpop(args ...string) string {
	return NewLPop(args...).Execute()
}

type RPop struct {
	args []string
}

func NewRPop(args ...string) *RPop {
	return &RPop{args}
}

func (rp *RPop) Execute() string {
	key := rp.args[0]
	if !store.Store.Contains(key) {
		return protocol.NewArray([]protocol.RespValue{}).Serialize()
	}

	val := store.Store.Get(key).(*store.RedisList).Pop()
	return protocol.BulkString{Val: val}.Serialize()
}

func Rpop(args ...string) string {
	return NewRPop(args...).Execute()
}

func init() {
	CommandRegistry["lpush"] = Lpush
	CommandRegistry["rpush"] = Rpush
	CommandRegistry["lrange"] = Lrange
	CommandRegistry["lpop"] = Lpop
	CommandRegistry["rpop"] = Rpop
}
