package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const SERVER_ADDRESS = "192.168.20.65:65500"


func main() {
	c, err := net.Dial("tcp", SERVER_ADDRESS)
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to: ", c.LocalAddr().String())
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		_, err = c.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Server time out!")
			c.Close()
		}

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
