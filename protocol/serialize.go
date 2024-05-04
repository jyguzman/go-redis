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
	Format() string
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

func (i Integer) Format() string {
	return fmt.Sprintf("(integer) %d\n", i.Val)
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

func (ss SimpleString) Format() string {
	return fmt.Sprintf("+%s\n", ss.Val)
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

func (bs BulkString) Format() string {
	return fmt.Sprintf("\"%s\"\n", bs.Val)
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

func (e Error) Format() string {
	return fmt.Sprintf("-%s\n", e.Val)
}

type Array struct {
	Val []RespValue
}

func (ar *Array) Type() rune {
	return RespArray
}

func (ar *Array) Serialize() string {
	if len(ar.Val) == 0 {
		return NilArray()
	}
	var sb strings.Builder
	sb.WriteString("*" + strconv.Itoa(len(ar.Val)) + "\r\n")
	for _, val := range ar.Val {
		sb.WriteString(val.Serialize())
	}
	return sb.String()
}

func (ar *Array) Format() string {
	if len(ar.Val) == 0 {
		return Nil{RespArray}.Format()
	}
	var sb strings.Builder
	for i, val := range ar.Val {
		sb.WriteString(fmt.Sprintf("%d) %s", i+1, val.Format()))
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
	NilType rune
}

func (n Nil) Type() rune {
	return n.NilType
}

func (n Nil) Serialize() string {
	if n.NilType == RespArray {
		return "*-1\r\n"
	}
	return "$-1\r\n"
}

func (n Nil) Format() string {
	return "(nil)\n"
}

func NilString() string {
	return Nil{NilType: RespBulkString}.Serialize()
}

func NilArray() string {
	return Nil{NilType: RespArray}.Serialize()
}

func Ok() string {
	return SimpleString{Val: "OK"}.Serialize()
}

func Err(msg string) string {
	return Error{Val: fmt.Sprintf("ERROR: %s", msg)}.Serialize()
}

func IntegerResponse(i int) string {
	return Integer{Val: i}.Serialize()
}

func BulkStringResponse(s string) string {
	return BulkString{Val: s}.Serialize()
}
