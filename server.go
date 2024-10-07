package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const ADDRESS = "192.168.20.65:65500"

func main() {
	connections := []net.Conn{}
	fmt.Println("Starting server...")
	l, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("New connection!")
		connections = append(connections, c)
		go readMessages(c)

	}
}

func readMessages(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			c.Close()
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("User logged out!")
			c.Close()
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}

}
