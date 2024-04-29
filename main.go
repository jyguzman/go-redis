package main

import (
	"fmt"
	"go-redis/protocol"
)

func main() {
	arrayOne := []protocol.RespValue{
		protocol.Integer{Val: 50},
		protocol.SimpleString{Val: "He was number one!"},
		protocol.BulkString{Val: "Hey Redis?"},
	}
	arrayTwo := []protocol.RespValue{
		protocol.Integer{Val: 100},
		protocol.SimpleString{Val: "I'm in a second array!"},
		protocol.BulkString{Val: "Hello world!"},
	}
	arrayThree := []protocol.RespValue{
		protocol.Array{Val: arrayOne},
		protocol.Array{Val: arrayTwo},
	}
	respArray := protocol.Array{Val: arrayThree}
	fmt.Println(respArray.Serialize().Deserialize())
}
