package main

import (
	"fmt"
	"go-redis/commands"
	"go-redis/data_types"
	"go-redis/protocol"
	"go-redis/store"
	"strings"
)

func DeserializeRequest(clientRequest string) ([]string, error) {
	deserialized, _, err := protocol.DeserializeMessage(clientRequest)
	respArray, ok := deserialized.(*protocol.Array)
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
	name := args[0]
	command, registered := commands.CommandRegistry[strings.ToLower(name)]
	if !registered {
		return protocol.Err(fmt.Sprintf("Unknown command: %s\n", strings.ToUpper(name)))
	}

	response, err := command(args...)
	if err != nil {
		return protocol.Err(err.Error())
	}

	return response
}

func main() {
	DoCommand("rpush", "new-list", "Append1")
	DoCommand("lpush", "new-list", "Prepend1")
	DoCommand("rpush", "new-list", "Append2")
	DoCommand("lpush", "new-list", "Prepend2")
	list := store.Store.Get("new-list").(*data_types.RedisList).Head
	ptr := list
	for ptr != nil {
		fmt.Print(ptr.RedisString.Value + " ")
		ptr = ptr.Next
	}
	//fmt.Println()
	// fmt.Println(store.Store.Get("new-list"))
	//fmt.Println(DoCommand("lrange", "new-list", "1", "7"))
	//fmt.Println(DoCommand("lpop", "new-list"))
	//fmt.Println(DoCommand("rpop", "new-list"))
	//fmt.Println(DoCommand("rpop", "new-list"))
	//fmt.Println(DoCommand("lrange", "new-list", "0", "7"))
	//listTwo := store.Store.Get("new-list").(*store.RedisList).Head
	//ptrTwo := listTwo
	//for ptrTwo != nil {
	//	fmt.Print(ptrTwo.RedisString.Value + " ")
	//	ptrTwo = ptrTwo.Next
	//}
}
