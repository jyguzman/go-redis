package commands

import (
	"fmt"
	"go-redis/data_types"
	"go-redis/protocol"
	"go-redis/store"
	"strconv"
)

type SetCommand struct {
	args []string
}

func (sc *SetCommand) Args() []string {
	return sc.args
}

func (sc *SetCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 3 {
		return "", fmt.Errorf("wrong number of arguments for SET")
	}
	store.Store.Set(sc.args[1], data_types.NewRedisString(sc.args[2]))
	return protocol.Ok(), nil
}

func NewSetCommand(args []string) *SetCommand {
	return &SetCommand{args: args}
}

func Set(args ...string) (string, error) {
	return NewSetCommand(args).Execute()
}

type MSetCommand struct {
	args []string
}

func (ms *MSetCommand) Args() []string {
	return ms.args
}

func (ms *MSetCommand) Execute() (string, error) {
	if ms.args == nil || len(ms.args) < 3 {
		return "", fmt.Errorf("wrong number of arguments for MSET")
	}

	key, items := ms.args[1], ms.args[2:]
	_, isString := store.Store.Get(key).(*data_types.RedisString)
	if !isString {
		return "", fmt.Errorf("value not of type")
	}

	for i := 2; i < len(items); i += 2 {
		store.Store.Set(items[i], data_types.NewRedisString(items[i+1]))
	}

	return protocol.Ok(), nil
}

func NewMSetCommand(args []string) *MSetCommand {
	return &MSetCommand{args: args}
}

func MSet(args ...string) (string, error) {
	return NewMSetCommand(args).Execute()
}

type GetCommand struct {
	args []string
}

func (gc *GetCommand) Args() []string {
	return gc.args
}

func (gc *GetCommand) Execute() (string, error) {
	if gc.args == nil || len(gc.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for GET")
	}

	key := gc.args[1]
	if !store.Store.Contains(key) {
		fmt.Println("im here")
		return protocol.NilString(), nil
	}

	rString, isString := store.Store.Get(key).(data_types.RedisString)
	if !isString {
		return "", fmt.Errorf("value not of type string")
	}

	return protocol.BulkStringResponse(rString.Value), nil
}

func NewGetCommand(args []string) *GetCommand {
	return &GetCommand{args: args}
}

func Get(args ...string) (string, error) {
	return NewGetCommand(args).Execute()
}

type IncrCommand struct {
	args []string
}

func (ic *IncrCommand) Args() []string {
	return ic.args
}

func (ic *IncrCommand) Execute() (string, error) {
	if ic.args == nil || len(ic.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for INCR")
	}

	key := ic.args[1]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisString("0"))
	}

	val := store.Store.Get(key)
	rString, ok := val.(data_types.RedisString)
	if !ok {
		return "", fmt.Errorf("value not of type string")
	}

	num, err := strconv.Atoi(rString.Value)
	if err != nil {
		return "", fmt.Errorf("failed to parse value as integer")
	}

	num += 1
	store.Store.Set(key, data_types.NewRedisString(strconv.Itoa(num)))
	return protocol.IntegerResponse(num), nil
}

func NewIncrCommand(args []string) *IncrCommand {
	return &IncrCommand{args: args}
}

func Incr(args ...string) (string, error) {
	return NewIncrCommand(args).Execute()
}

type DecrCommand struct {
	args []string
}

func (dc *DecrCommand) Args() []string {
	return dc.args
}

func (dc *DecrCommand) Execute() (string, error) {
	if dc.args == nil || len(dc.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for DECR")
	}

	key := dc.args[1]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisString("0"))
	}

	rString, isString := store.Store.Get(key).(data_types.RedisString)
	if !isString {
		return "", fmt.Errorf("value not of type string")
	}

	num, err := strconv.Atoi(rString.Value)
	if err != nil {
		return "", fmt.Errorf("failed to parse value as integer")
	}

	num -= 1
	store.Store.Set(key, data_types.NewRedisString(strconv.Itoa(num)))
	return protocol.IntegerResponse(num), nil
}

func NewDecrCommand(args []string) *DecrCommand {
	return &DecrCommand{args: args}
}

func Decr(args ...string) (string, error) {
	return NewDecrCommand(args).Execute()
}

func init() {
	CommandRegistry["set"] = Set
	CommandRegistry["get"] = Get
	CommandRegistry["incr"] = Incr
	CommandRegistry["decr"] = Decr
	CommandRegistry["mset"] = MSet
}
