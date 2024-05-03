package main

import (
	"fmt"
	"go-redis/protocol"
)

type Cli struct {
	client Client
}

func (c *Cli) Prompt(cmd string, args ...string) {

}

func displayResponse(value protocol.RespValue) {
	switch value.(type) {
	case protocol.Integer:
		displayInteger(value.(protocol.Integer))
	case protocol.BulkString:
		displayBulkString(value.(protocol.BulkString))
	case protocol.Error:
		displayError(value.(protocol.Error))
	case protocol.SimpleString:
		displaySimpleString(value.(protocol.SimpleString))
	case *protocol.Array:
		displayArrayResponse(value.(*protocol.Array))
	case protocol.Nil:
		displayNilResponse()
	default:
		fmt.Printf("received invalid response\n")
	}
}

func displayInteger(integer protocol.Integer) {
	fmt.Printf("(integer) %d\n", integer.Val)
}

func displaySimpleString(ss protocol.SimpleString) {
	fmt.Printf("+%s\n", ss.Val)
}

func displayBulkString(str protocol.BulkString) {
	fmt.Printf("\"%s\"\n", str.Val)
}

func displayError(err protocol.Error) {
	fmt.Printf("%s\n", err.Val)
}

func displayArrayResponse(array *protocol.Array) {
	for i, item := range array.Val {
		fmt.Printf("%d) ", i+1)
		displayResponse(item)
	}
}

func displayNilResponse() {
	fmt.Printf("(nil)\n")
}
