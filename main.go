package main

import (
	"flag"
	"log"

	"go-chat/client"
	"go-chat/server"
)

func main() {

	ipAddress := flag.String("ipAddress", "", "* specifies the ip address to create (server) or connect (client)")
	mode := flag.String("mode", "client", "starts as client or server. client as default")
	flag.Parse()

	if *ipAddress == "" {
		log.Fatalln("--ipAddress flag is required")
		return
	}
	if *mode != "server" && *mode != "client" {
		log.Fatalln("--mode flag must be [client, server]. client is default")
		return
	}

	switch *mode {
	case "server":
		server.Server(*ipAddress)
	case "client":
		client.Client(*ipAddress)
	}
}
