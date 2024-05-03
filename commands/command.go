package commands

import (
	"go-redis/protocol"
	"go-redis/store"
)

var CommandRegistry = make(map[string]func(...string) (string, error))

type Command interface {
	Args() []string
	Execute() (string, error)
}

type ExistsCommand struct {
	args []string
}

func (ec *ExistsCommand) Args() []string {
	return ec.args
}

func (ec *ExistsCommand) Execute() (string, error) {
	numExists := 0
	for _, key := range ec.args {
		if store.Store.Contains(key) {
			numExists++
		}
	}
	return protocol.IntegerResponse(numExists), nil
}

func NewExistsCommand(keys []string) *ExistsCommand {
	return &ExistsCommand{args: keys}
}

func Exists(keys ...string) (string, error) {
	return NewExistsCommand(keys).Execute()
}

type DelCommand struct {
	args []string
}

func (dc *DelCommand) Args() []string {
	return dc.args
}

func (dc *DelCommand) Execute() (string, error) {
	numDeleted := 0
	for _, key := range dc.args {
		if store.Store.Contains(key) {
			store.Store.Remove(key)
			numDeleted++
		}
	}
	return protocol.IntegerResponse(numDeleted), nil
}

func NewDelCommand(args []string) *DelCommand {
	return &DelCommand{args: args}
}

func Del(keys ...string) (string, error) {
	return NewDelCommand(keys).Execute()
}

type FlushDBCommand struct {
	args []string
}

func (fdc *FlushDBCommand) Args() []string {
	return fdc.args
}

func (fdc *FlushDBCommand) Execute() (string, error) {
	numFlushed := 0
	return protocol.IntegerResponse(numFlushed), nil
}

type SaveCommand struct {
	state store.Cache
}

func init() {
	CommandRegistry["exists"] = Exists
	CommandRegistry["del"] = Del
}
