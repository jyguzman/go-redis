package commands

import (
	"go-redis/protocol"
	"go-redis/store"
)

type Command interface {
	Execute() string
}

var CommandRegistry = map[string]func(...string) string{
	"exists": Exists, "set": Set, "get": Get,
}

type ExistsCommand struct {
	keys []string
}

func (ec *ExistsCommand) Execute() string {
	numExists := 0
	for _, key := range ec.keys {
		if store.Store.Contains(key) {
			numExists++
		}
	}
	return protocol.Integer{Val: numExists}.Serialize()
}

func NewExistsCommand(keys []string) *ExistsCommand {
	return &ExistsCommand{keys: keys}
}

func Exists(keys []string) string {
	return NewExistsCommand(keys).Execute()
}

type DelCommand struct {
	keys []string
}

func (dc *DelCommand) Execute() string {
	numDeleted := 0
	for _, key := range dc.keys {
		if store.Store.Contains(key) {
			store.Store.Remove(key)
			numDeleted++
		}
	}
	return protocol.Integer{Val: numDeleted}.Serialize()
}

func NewDelCommand(keys []string) *DelCommand {
	return &DelCommand{keys: keys}
}

func Del(keys []string) string {
	return NewDelCommand(keys).Execute()
}

type SaveCommand struct {
	state store.Cache
}
