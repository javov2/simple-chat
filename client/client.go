package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	globals "go-chat/config"
)

const UserPrompt = "[%s][%s] >> "

func Client(serverAddress string) {
	fmt.Println("Starting application as client in ", serverAddress, "...")
	status := "X"
	fmt.Print("Connecting to: ", serverAddress+"... ")
	c, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Print("[ERROR]")
		return
	}
	defer c.Close()
	fmt.Print("[OK]")
	fmt.Println()
	status = "CONNECTED"
	// login
	sessionId := login(c)
	go printServerMessages(c)
	for {
		// reading user input
		fmt.Printf(UserPrompt, status, sessionId)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		_, err = c.Write([]byte(text))
		if err != nil {
			return
		}
	}
}

func login(conn net.Conn) string {
	fmt.Print("Logging in... \n")
	conn.Write([]byte(globals.Commands.Login + "\n"))
	reader := bufio.NewReader(conn)
	msg, _ := reader.ReadString('\n')
	fmt.Printf("Your Session Id: %s", msg)
	fmt.Println()
	return strings.TrimRight(msg, "\r\n")
}

func printServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Print(msg)
	}
}
