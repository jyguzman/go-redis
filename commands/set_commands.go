package commands

import (
	"errors"
	"fmt"
	"go-redis/data_types"
	"go-redis/protocol"
	"go-redis/store"
)

type SaddCommand struct {
	args []string
}

func (sc *SaddCommand) Args() []string {
	return sc.args
}

func (sc *SaddCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 3 {
		return "", fmt.Errorf("wrong number of arguments for SADD")
	}

	key, items := sc.args[1], sc.args[2:]
	if !store.Store.Contains(key) {
		store.Store.Set(key, data_types.NewRedisSet())
	}

	_, isSet := store.Store.Get(key).(*data_types.RedisSet)
	if !isSet {
		return "", errors.New("value not of type set")
	}

	numAdded := 0
	for _, item := range items {
		numAdded += store.Store.Get(key).(*data_types.RedisSet).Add(item)
	}

	return protocol.IntegerResponse(numAdded), nil
}

func NewSaddCommand(args []string) *SaddCommand {
	return &SaddCommand{args: args}
}

func SAdd(args ...string) (string, error) {
	return NewSaddCommand(args).Execute()
}

type SremCommand struct {
	args []string
}

func (sc *SremCommand) Args() []string {
	return sc.args
}

func (sc *SremCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 3 {
		return "", fmt.Errorf("wrong number of arguments for SREM")
	}

	key, items := sc.args[1], sc.args[2:]
	if !store.Store.Contains(key) {
		return protocol.IntegerResponse(0), nil
	}

	_, isSet := store.Store.Get(key).(*data_types.RedisSet)
	if !isSet {
		return "", fmt.Errorf("value not of type set")
	}

	numDeleted := 0
	for _, item := range items {
		numDeleted += store.Store.Get(key).(*data_types.RedisSet).Remove(item)
	}

	return protocol.IntegerResponse(numDeleted), nil
}

func NewSremCommand(args []string) *SremCommand {
	return &SremCommand{args: args}
}

func SRem(args ...string) (string, error) {
	return NewSremCommand(args).Execute()
}

type SMembersCommand struct {
	args []string
}

func (sc *SMembersCommand) Args() []string {
	return sc.args
}

func (sc *SMembersCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for SMEMBERS")
	}
	key := sc.args[1]
	if !store.Store.Contains(key) {
		return protocol.NilArray(), nil
	}

	set, isSet := store.Store.Get(key).(*data_types.RedisSet)
	if !isSet {
		return "", errors.New("value not of type set")
	}

	array := protocol.NewArray([]protocol.RespValue{})
	for _, member := range set.Values() {
		array.Add(protocol.BulkString{Val: member.Value})
	}
	return array.Serialize(), nil
}

func NewSMembersCommand(args []string) *SMembersCommand {
	return &SMembersCommand{args: args}
}

func SMembers(args ...string) (string, error) {
	return NewSMembersCommand(args).Execute()
}

type SCardCommand struct {
	args []string
}

func (sc *SCardCommand) Args() []string {
	return sc.args
}

func (sc *SCardCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for SCARD")
	}

	key := sc.args[1]
	if !store.Store.Contains(key) {
		return protocol.IntegerResponse(0), nil
	}

	set, isSet := store.Store.Get(key).(*data_types.RedisSet)
	if !isSet {
		return "", errors.New("value not of type set")
	}

	return protocol.IntegerResponse(set.Size()), nil
}

func NewSCardCommand(args []string) *SCardCommand {
	return &SCardCommand{args: args}
}

func SCard(args ...string) (string, error) {
	return NewSCardCommand(args).Execute()
}

type SIsMemberCommand struct {
	args []string
}

func (sc *SIsMemberCommand) Args() []string {
	return sc.args
}

func (sc *SIsMemberCommand) Execute() (string, error) {
	if sc.args == nil || len(sc.args) < 2 {
		return "", fmt.Errorf("wrong number of arguments for SISMEMBER")
	}

	key, items := sc.args[1], sc.args[2:]
	if !store.Store.Contains(key) {
		return protocol.IntegerResponse(0), nil
	}

	set, isSet := store.Store.Get(key).(*data_types.RedisSet)
	if !isSet {
		return "", errors.New("value not of type set")
	}

	array := protocol.NewArray([]protocol.RespValue{})
	for _, item := range items {
		if set.Contains(item) {
			array.Add(protocol.Integer{Val: 1})
		} else {
			array.Add(protocol.Integer{Val: 0})
		}
	}
	return array.Serialize(), nil
}

func NewSIsMemberCommand(args []string) *SIsMemberCommand {
	return &SIsMemberCommand{args: args}
}

func SIsMember(args ...string) (string, error) {
	return NewSIsMemberCommand(args).Execute()
}

func init() {
	CommandRegistry["sadd"] = SAdd
	CommandRegistry["srem"] = SRem
	CommandRegistry["smembers"] = SMembers
	CommandRegistry["scard"] = SCard
	CommandRegistry["sismember"] = SIsMember
}
