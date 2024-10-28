package server

import (
	"bufio"

	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"

	"go-chat/config"
)

const CLIENT_MAX_IDLE_TIME = 4 * 60 * 1000000000

type ConnectionEvent struct {
	connection Connection
	command    string
}

type Connection struct {
	id   string
	conn net.Conn
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
		command := r.command
		connection := r.connection.conn
		connectionId := r.connection.id

		switch command {
		case config.Commands.Login:
			fmt.Println("[LOGIN] " + connectionId)
			conns.Store(connectionId, connection)
		case config.Commands.Logout:
			fmt.Println("[LOGOUT] " + connectionId)
			tmpC, _ := conns.Load(connectionId)
			tmpC.(net.Conn).Close()
			conns.Delete(connectionId)
		case config.Commands.SysDisconnect:
			fmt.Println("[DISCONNECTED] " + connectionId)
			tmpC, _ := conns.Load(connectionId)
			tmpC.(net.Conn).Close()
			conns.Delete(connectionId)
		case config.Commands.Subscribe:
			fmt.Println("[SUBSCRIBE][<TOPIC>] " + connectionId)
			conns.Store(connectionId, connectionId)
		}

	}
}

func handleSession(c net.Conn, connsChannel chan ConnectionEvent) {
	isLoggedIn := false
	connUUID := uuid.NewString()
	connection := Connection{id: connUUID, conn: c}
	reader := bufio.NewReader(c)
	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			connsChannel <- ConnectionEvent{
				command:    config.Commands.SysDisconnect,
				connection: connection,
			}
			return
		}
		if !isLoggedIn && strings.TrimSpace(string(netData)) == config.Commands.Login {
			isLoggedIn = true
			sendMessage(c, connUUID)
			connsChannel <- ConnectionEvent{
				command:    config.Commands.Login,
				connection: connection,
			}
		}
		if isLoggedIn && strings.TrimSpace(string(netData)) == config.Commands.Logout {
			sendMessage(c, "[OK]")
			connsChannel <- ConnectionEvent{
				command:    config.Commands.Logout,
				connection: connection,
			}
			return
		}

		fmt.Print(connUUID+" -> ", string(netData))
	}
}

func sendMessage(c net.Conn, message string) {
	c.Write([]byte(message + "\n"))
}
