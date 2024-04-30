package main

import (
	"fmt"
	"go-redis/commands"
)

func doCommand(command commands.Command) string {
	return command.Execute()
}

func main() {
	//arrayOne := []protocol.RespValue{
	//	protocol.Integer{Val: 50},
	//	protocol.SimpleString{Val: "He was number one!"},
	//	protocol.BulkString{Val: "Hey Redis?"},
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
	//valOne := store.RedisList{
	//	Values: []store.RedisValue{
	//		store.RedisString{Value: "Hello, "},
	//	},
	//}
	//valOne.Push(store.RedisString{Value: "world!"})
	//lpush := commands.NewLPush("newList", store.RedisString{Value: "Hello, world!"})
	//lpush.Execute()
	fmt.Println(commands.Set("name", "Jordie Guzman"))
	fmt.Println(commands.Get("name"))

	fmt.Println(commands.Set("number", "50"))
	fmt.Println(commands.Incr("number"))
	fmt.Println(commands.Get("number"))

	fmt.Println(doCommand(commands.NewSetCommand("newThing", "oldThing")))
	//respArray := protocol.Array{Val: arrayThree}

	// fmt.Println(store.Store.Get("newList"))
}
