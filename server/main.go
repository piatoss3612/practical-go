package main

import (
	"log/slog"
	"net"
	"strings"
)

var helpMessage = `
Commands:
  create <room> - create a new room
  join <room>   - join a room
  nick <name>   - change nickname
  list          - list all rooms
  quit          - quit chat server
  help          - show help
`

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("Chat server started on port 8080")

	newConnections := make(chan net.Conn)
	defer close(newConnections)

	go handleConnections(newConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		newConnections <- conn
	}
}

type Client struct {
	nickname string
	net.Conn
}

func handleConnections(conns <-chan net.Conn) {
	for conn := range conns {
		go handleUser(&Client{nickname: "unknown", Conn: conn})
	}
}

func handleUser(client *Client) {
	buf := make([]byte, 1024)

	for {
		_, err := client.Write([]byte("Enter command:"))
		if err != nil {
			slog.Error(err.Error())
			return
		}

		n, err := client.Read(buf)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		msg := string(buf[:n])
		fields := strings.Fields(msg)

		switch fields[0] {
		case "create":
			// create room
		case "join":
			// join room
		case "nick":
			// change nickname
		case "list":
			// list rooms
		case "quit":
			// quit chat server
		case "help":
			// show help
			_, _ = client.Write([]byte(helpMessage))
		default:
			_, _ = client.Write([]byte("Unknown command\n"))
		}
	}
}
