package main

import (
	"fmt"
	"go-redis/protocol"
	"strings"
)

func serializeRequest(request string) string {
	var bulkStrings []protocol.RespValue
	for _, s := range strings.Fields(request) {
		bulkStrings = append(bulkStrings, protocol.BulkString{Val: s})
	}
	return protocol.Array{Val: bulkStrings}.Serialize()
}

func main() {
	const (
		Port = 6379
		Host = "localhost"
	)
	thing := "$13\r\nHello, world!\r\n"
	fmt.Println(strings.Split(thing, "\r\n"))
}
