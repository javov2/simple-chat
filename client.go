package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddress := os.Args[len(os.Args)-1]
	fmt.Println("Attempting connection to: ", serverAddress)
	c, err := net.Dial("tcp", serverAddress)
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
		_, err = c.Write([]byte(text))
		if err != nil {
			return
		}

	}
}
