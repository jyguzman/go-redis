package main

import (
	"fmt"
	"go-redis/commands"
	"go-redis/protocol"
	"go-redis/store"
)

func DeserializeRequest(clientRequest string) ([]string, error) {
	deserialized, _, err := protocol.DeserializeMessage(clientRequest)
	respArray, ok := deserialized.(protocol.Array)
	if !ok {
		return nil, fmt.Errorf("response from client must be array")
	}
	bulkStringArray := respArray.Val
	commandArgs := make([]string, len(bulkStringArray))
	for i, val := range bulkStringArray {
		bs, ok := val.(protocol.BulkString)
		if !ok {
			return nil, fmt.Errorf("client request array must be of bulk strings")
		}
		commandArgs[i] = bs.Val
	}
	fmt.Println(commandArgs, err)
	return commandArgs, nil
}

func DoCommand(args ...string) string {
	command, known := commands.CommandRegistry[args[0]]
	if !known {
		return protocol.Err(fmt.Sprintf("Unknown command: %s\n", args[0]))
	}
	return command(args[1:]...)
}

func main() {
	//arrayOne := []protocol.RespValue{
	//	protocol.Integer{Val: 50},
	//	protocol.SimpleString{Val: "He was number one!"},
	//	protocol.BulkString{Val: "Hey Redis?"},
	//	protocol.Error{Val: "ERROR: Invalid type."},
	//}
	//arrayTwo := []protocol.RespValue{
	//	protocol.Integer{Val: 100},
	//	protocol.SimpleString{Val: "I'm in a second array!"},
	//	protocol.BulkString{Val: "Hello world!"},
	//}
	//arrayThree := []protocol.RespValue{
	//	protocol.Array{Val: arrayOne},
	//	protocol.Array{Val: arrayTwo},
	//}
	//arraySer := protocol.Array{Val: arrayThree}.Serialize()
	//val, idx, err := protocol.DeserializeMessage(arraySer)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(val, idx)
	list := store.RedisList{Size: 0, Head: nil, Tail: nil}
	list.Append(store.RedisString{Value: "Append1"})
	list.Append(store.RedisString{Value: "Append2"})
	list.Prepend(store.RedisString{Value: "Prepend1"})
	list.Append(store.RedisString{Value: "Append3"})
	list.Prepend(store.RedisString{Value: "Prepend2"})
	list.Prepend(store.RedisString{Value: "Prepend3"})
	list.Prepend(store.RedisString{Value: "Prepend4"})
	list.Append(store.RedisString{Value: "Append4"})
	ptr := list.Head
	for ptr != nil {
		fmt.Print(ptr.Value.Value + " ")
		ptr = ptr.Next
	}
	fmt.Println(DoCommand("get", "name"))
	fmt.Println(DoCommand("incr", "number"))
	fmt.Println(DoCommand("incr", "number"))
	fmt.Println(DoCommand("incr", "number"))
	fmt.Println(DoCommand("decr", "number"))
	fmt.Println(DoCommand("get", "number"))
	fmt.Println(DoCommand("set", "name", "jordie"))
	fmt.Println(DoCommand("get", "name", "jordie"))
}
