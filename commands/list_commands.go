package commands

import (
	"fmt"
	"go-redis/data_types"
	"go-redis/protocol"
	"go-redis/store"
	"strconv"
)

type LPush struct {
	args []string
}

func (lp *LPush) Args() []string {
	return lp.args
}

func (lp *LPush) Execute() (string, error) {
	if lp.args == nil || len(lp.args) < 2 {
		return "", fmt.Errorf("not enough arguments for LPUSH")
	}

	key, item := lp.args[1], lp.args[2]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisList())
	}

	_, ok := store.Store.Get(key).(*data_types.RedisList)
	if !ok {
		return "", fmt.Errorf("value is not a list")
	}

	numItems := store.Store.Get(key).(*data_types.RedisList).Prepend(data_types.NewRedisString(item))
	return protocol.IntegerResponse(numItems), nil
}

func NewLPush(args []string) *LPush {
	return &LPush{args: args}
}

func Lpush(args ...string) (string, error) {
	return NewLPush(args).Execute()
}

type RPush struct {
	args []string
}

func (rp *RPush) Args() []string {
	return rp.args
}

func (rp *RPush) Execute() (string, error) {
	if rp.args == nil || len(rp.args) < 2 {
		return "", fmt.Errorf("not enough arguments for RPUSH")
	}

	key, item := rp.args[1], rp.args[2]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisList())
	}

	_, ok := store.Store.Get(key).(*data_types.RedisList)
	if !ok {
		return "", fmt.Errorf("value is not a list")
	}

	numItems := store.Store.Get(key).(*data_types.RedisList).Append(data_types.NewRedisString(item))
	return protocol.IntegerResponse(numItems), nil
}

func NewRPush(args []string) *RPush {
	return &RPush{args: args}
}

func Rpush(args ...string) (string, error) {
	return NewRPush(args).Execute()
}

type LRange struct {
	args []string
}

func (lr *LRange) Args() []string {
	return lr.args
}

func (lr *LRange) Execute() (string, error) {
	if lr.args == nil || len(lr.args) < 3 {
		return "", fmt.Errorf("not enough arguments for LRANGE (need key, start, stop)")
	}

	key, startStr, stopStr := lr.args[0], lr.args[1], lr.args[2]
	start, err := strconv.Atoi(startStr)
	stop, errTwo := strconv.Atoi(stopStr)
	if err != nil || errTwo != nil {
		return "", fmt.Errorf("could not parse start or stop as integers")
	}

	if !store.Store.Contains(key) {
		return protocol.NilArray(), nil
	}

	ptr := store.Store.Get(key).(*data_types.RedisList).Head
	for i := 0; i < start; i++ {
		ptr = ptr.Next
	}
	array := &protocol.Array{Val: []protocol.RespValue{}}
	for i := 0; ptr != nil && i < (stop-start)+1; i++ {
		array.Add(protocol.BulkString{Val: ptr.RedisString.Value})
		ptr = ptr.Next
	}
	return array.Serialize(), nil
}

func NewLRange(args []string) *LRange {
	return &LRange{args: args}
}

func Lrange(args ...string) (string, error) {
	return NewLRange(args).Execute()
}

type LPop struct {
	args []string
}

func NewLPop(args []string) *LPop {
	return &LPop{args: args}
}

func (lp *LPop) Execute() (string, error) {
	if lp.args == nil || len(lp.args) < 2 {
		return "", fmt.Errorf("not enough arguments for LPOP")
	}

	key := lp.args[0]
	if !store.Store.Contains(key) {
		return protocol.NilArray(), nil
	}

	val := store.Store.Get(key).(*data_types.RedisList).PopLeft()
	return protocol.BulkStringResponse(val), nil
}

func Lpop(args ...string) (string, error) {
	return NewLPop(args).Execute()
}

type RPop struct {
	args []string
}

func NewRPop(args []string) *RPop {
	return &RPop{args: args}
}

func (rp *RPop) Execute() (string, error) {
	if rp.args == nil || len(rp.args) < 2 {
		return "", fmt.Errorf("not enough arguments for RPOP")
	}

	key := rp.args[0]
	if !store.Store.Contains(key) {
		return protocol.NilArray(), nil
	}

	val := store.Store.Get(key).(*data_types.RedisList).Pop()
	return protocol.BulkStringResponse(val), nil
}

func Rpop(args ...string) (string, error) {
	return NewRPop(args).Execute()
}

type LIndex struct {
	args []string
}

func (li *LIndex) Args() []string {
	return li.args
}

func (li *LIndex) Execute() (string, error) {
	if li.args == nil || len(li.args) < 3 {
		return "", fmt.Errorf("not enough arguments for LINDEX")
	}

	key := li.args[1]
	if !store.Store.Contains(key) {
		return protocol.NilString(), nil
	}

	indexStr := li.args[2]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("could not parse index as integer")
	}

	list, isList := store.Store.Get(key).(*data_types.RedisList)
	if !isList {
		return "", fmt.Errorf("value is not of type list")
	}

	if index < 0 || index >= list.Size {
		return protocol.NilString(), nil
	}

	ptr := list.Head
	for i := 0; i < index; i++ {
		ptr = ptr.Next
	}

	return protocol.BulkStringResponse(ptr.RedisString.Value), nil
}

func NewLIndex(args []string) *LIndex {
	return &LIndex{args: args}
}

func Lindex(args ...string) (string, error) {
	return NewLIndex(args).Execute()
}

func init() {
	CommandRegistry["lpush"] = Lpush
	CommandRegistry["lpop"] = Lpop
	CommandRegistry["rpush"] = Rpush
	CommandRegistry["rpop"] = Rpop
	CommandRegistry["lrange"] = Lrange
	CommandRegistry["lindex"] = Lindex
}
