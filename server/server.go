package main

import (
	"bufio"
	"fmt"
	"go-redis/commands"
	"go-redis/protocol"
	"log"
	"net"
	"strings"
)

func DeserializeRequest(clientRequest string) ([]string, error) {
	deserialized, _, err := protocol.DeserializeMessage(clientRequest)
	if err != nil {
		return nil, err
	}
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := bufio.NewReader(conn).Read(buffer)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("got " + string(buffer[:n]))
		desReq, reqErr := DeserializeRequest(string(buffer[:n]))
		if reqErr != nil {
			log.Println(reqErr)
		}
		res := DoCommand(desReq...)
		fmt.Printf("req: %v", desReq)
		_, writeErr := conn.Write([]byte(res))
		if writeErr != nil {
			log.Println(writeErr)
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}
