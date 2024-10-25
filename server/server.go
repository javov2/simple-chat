package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type ConnectionEvent struct {
	id      string
	conn    net.Conn
	command string
}

const ServerStatusPrompt = "Number of active connections [%d]"

func Server(serverAddress string) {
	conns := sync.Map{}
	connsChannel := make(chan ConnectionEvent)
	l, err := net.Listen("tcp", serverAddress+":0")

	fmt.Println("Starting application as server in ", l.Addr().String(), "...")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conns.Clear()
	defer l.Close()

	go manageConnections(&conns, connsChannel)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleSession(c, connsChannel)
	}
}

func manageConnections(conns *sync.Map, connsChannel chan ConnectionEvent) {
	for {
		r := <-connsChannel
		if r.command == "/login" {
			fmt.Println("[LOGIN] " + r.id)
			conns.Store(r.id, r.conn)
		}
		if r.command == "/logout" {
			fmt.Println("[LOGOUT] " + r.id)
			tmpC, _ := conns.Load(r.id)
			tmpC.(net.Conn).Close()
			conns.Delete(r.id)
		}
	}
}

func handleSession(c net.Conn, connsChannel chan ConnectionEvent) {
	isLoggedIn := false
	connUUID := uuid.NewString()
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[ERROR] %s %s", connUUID, err)
			c.Close()
			return
		}
		if !isLoggedIn && strings.TrimSpace(string(netData)) == "/login" {
			isLoggedIn = true
			sendMessage(c, connUUID)
			connsChannel <- ConnectionEvent{
				id:      connUUID,
				command: "/login",
				conn:    c,
			}
		}
		if isLoggedIn && strings.TrimSpace(string(netData)) == "/logout" {
			sendMessage(c, "[OK]")
			connsChannel <- ConnectionEvent{
				id:      connUUID,
				command: "/logout",
				conn:    c,
			}
			return
		}

		fmt.Print(connUUID+" -> ", string(netData))
	}
}

func sendMessage(c net.Conn, message string) {
	c.Write([]byte(message + "\n"))
}
