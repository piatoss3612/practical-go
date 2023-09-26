package main

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
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
	<-run()
}

func run() <-chan bool {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("Chat server started on port 8080")

	conns := make(chan net.Conn)

	go handleConnections(conns)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				slog.Error(err.Error())
				continue
			}

			conns <- conn
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	quit := make(chan bool)

	go func() {
		defer func() {
			_ = listener.Close()
			slog.Info("Chat server stopped")

			close(quit)
			close(sig)
			close(conns)
		}()
		<-sig
		quit <- true
	}()

	return quit
}

type Client struct {
	nickname string
	net.Conn
}

func handleConnections(conns <-chan net.Conn) {
	for conn := range conns {
		if conn == nil {
			continue
		}

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
