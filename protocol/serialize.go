package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	RespSimpleString = '+'
	RespBulkString   = '$'
	RespArray        = '*'
	RespError        = '-'
	RespInteger      = ':'
)

type RespValue interface {
	Type() rune
	Serialize() string
}

type Integer struct {
	Val int
}

func (i Integer) Type() rune {
	return RespInteger
}

func (i Integer) Serialize() string {
	return fmt.Sprintf(":%d\r\n", i.Val)
}

type SimpleString struct {
	Val string
}

func (ss SimpleString) Type() rune {
	return RespSimpleString
}

func (ss SimpleString) Serialize() string {
	return fmt.Sprintf("+%s\r\n", ss.Val)
}

type BulkString struct {
	Val string
}

func (bs BulkString) Type() rune {
	return RespBulkString
}

func (bs BulkString) Serialize() string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(bs.Val), bs.Val)
}

type Error struct {
	Val string
}

func (e Error) Type() rune {
	return RespError
}

func (e Error) Serialize() string {
	return fmt.Sprintf("-%s\r\n", e.Val)
}

type Array struct {
	Val []RespValue
}

func (ar *Array) Type() rune {
	return RespArray
}

func (ar *Array) Serialize() string {
	if len(ar.Val) == 0 {
		return "*-1\r\n"
	}
	var sb strings.Builder
	sb.WriteString("*" + strconv.Itoa(len(ar.Val)) + "\r\n")
	for _, val := range ar.Val {
		sb.WriteString(val.Serialize())
	}
	return sb.String()
}

func (ar *Array) Add(rv RespValue) {
	ar.Val = append(ar.Val, rv)
}

func NewArray(rvs []RespValue) *Array {
	return &Array{Val: rvs}
}

type Nil struct {
	Val RespValue
}

func (n Nil) Type() rune {
	return RespBulkString
}

func (n Nil) Serialize() string {
	return "$-1\r\n"
}

func Null() string {
	return "$-1\r\n"
}

func Ok() string {
	return SimpleString{Val: "OK"}.Serialize()
}

func Err(msg string) string {
	return Error{Val: fmt.Sprintf("ERROR: %s", msg)}.Serialize()
}
