package main

import (
	"bufio"
	"fmt"
	"go-redis/protocol"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func serializeRequest(request string) string {
	var bulkStrings []protocol.RespValue
	for _, s := range strings.Fields(request) {
		bulkStrings = append(bulkStrings, protocol.BulkString{Val: s})
	}
	return protocol.NewArray(bulkStrings).Serialize()
}

type Client struct {
	history []string
	host    string
	port    int
	conn    net.Conn
}

func (c *Client) History() []string {
	return c.history
}

func (c *Client) Connect(host string, port int) error {
	c.host = host
	c.port = port
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func NewClient(host string, port int) *Client {
	return &Client{host: host, port: port}
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) Communicate() {
	var input string
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Printf("127.0.0.1:6379> ")
	for stdin.Scan() {
		fmt.Printf("127.0.0.1:6379> ")
		input = stdin.Text()
		if input == "exit" {
			break
		}
		_, writeErr := c.conn.Write([]byte(serializeRequest(input)))
		if writeErr != nil {
			log.Fatal(writeErr)
		}
		c.history = append(c.history, input)
		buffer := make([]byte, 1024)
		n, err := bufio.NewReader(c.conn).Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		res, _, desErr := protocol.DeserializeMessage(string(buffer[:n]))
		if desErr != nil {
			log.Println(desErr)
		}
		displayResponse(res)
	}
}

func main() {
	const (
		Port = 6379
		Host = "localhost"
	)
	c := NewClient(Host, Port)
	err := c.Connect(Host, Port)
	if err != nil {
		log.Fatal(err)
	}
	c.Communicate()
}
