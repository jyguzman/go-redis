package protocol

import (
	"fmt"
	"strconv"
)

const (
	RespSimpleString = '+'
	RespBulkString   = '$'
	RespArray        = '*'
	RespError        = '-'
	RespInteger      = ':'
	RespMap          = '%'
	RespSet          = '~'
)

type RespValue interface {
	Type() rune
	Serialize() Message
}

type Integer struct {
	Val int
}

func (i Integer) Type() rune {
	return RespInteger
}

func (i Integer) Serialize() Message {
	return &IntegerMessage{
		IntegerMessage: fmt.Sprintf(":%d\r\n", i.Val),
		RespValue:      i,
	}
}

type SimpleString struct {
	Val string
}

func (ss SimpleString) Type() rune {
	return RespSimpleString
}

func (ss SimpleString) Serialize() Message {
	return &SimpleStringMessage{
		SimpleStringMessage: fmt.Sprintf("+%s\r\n", ss.Val),
		RespValue:           ss,
	}
}

type BulkString struct {
	Val string
}

func (bs BulkString) Type() rune {
	return RespBulkString
}

func (bs BulkString) Serialize() Message {
	return &BulkStringMessage{
		BulkStringMessage: fmt.Sprintf("$%d\r\n%s\r\n", len(bs.Val), bs.Val),
		RespValue:         bs,
	}
}

type Error struct {
	Val string
}

func (e Error) Type() rune {
	return RespError
}

func (e Error) Serialize() Message {
	return &ErrorMessage{
		ErrorMessage: fmt.Sprintf("-%s\r\n", e.Val),
		RespValue:    e,
	}
}

type Array struct {
	Val []RespValue
}

func (ar Array) Type() rune {
	return RespArray
}

func (ar Array) Serialize() Message {
	message := "*" + strconv.Itoa(len(ar.Val)) + "\r\n"
	for _, val := range ar.Val {
		message += val.Serialize().Message()
	}
	return &ArrayMessage{ArrayMessage: message, RespValue: ar}
}
