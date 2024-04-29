package commands

import (
	"go-redis/protocol"
)

type Command interface {
	Execute() error
}

type Set struct {
	Args protocol.Message
}
