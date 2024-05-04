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
	host string
	port int
	id   string
	conn net.Conn
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
	fmt.Printf("%s:%d> ", c.host, c.port)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		input = stdin.Text()
		if input == "exit" {
			break
		}
		_, writeErr := c.conn.Write([]byte(serializeRequest(input)))
		if writeErr != nil {
			log.Fatal(writeErr)
		}
		buffer := make([]byte, 1024)
		n, err := bufio.NewReader(c.conn).Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		res, _, desErr := protocol.DeserializeMessage(string(buffer[:n]))
		if desErr != nil {
			log.Println(desErr)
		}
		fmt.Print(res.Format())
		fmt.Printf("%s:%d> ", c.host, c.port)
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
