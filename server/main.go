package main

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

var (
	rooms = make(map[string]*Room)
	mu    = sync.Mutex{}
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

		go handleClient(&Client{nickname: "unknown", Conn: conn})
	}
}

func handleClient(client *Client) {
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
			if len(fields) < 2 {
				_, _ = client.Write([]byte("Missing room name"))
				continue
			}

			addRoom(fields[1], client)
			joinRoom(fields[1], client)
			return
		case "join":
			if len(fields) < 2 {
				_, _ = client.Write([]byte("Missing room name"))
				continue
			}

			joinRoom(fields[1], client)
			return
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
			_, _ = client.Write([]byte("Unknown command"))
		}
	}
}

func addRoom(name string, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	rooms[name] = NewRoom(name)
}

func joinRoom(name string, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	room, ok := rooms[name]
	if !ok {
		_, _ = client.Write([]byte("Room not found"))
		return
	}

	room.Join(client)
	room.Broadcast(client.nickname + " joined the room")

	go chatInRoom(room, client)
}

func chatInRoom(room *Room, client *Client) {
	buf := make([]byte, 1024)

	for {
		n, err := client.Read(buf)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		msg := string(buf[:n])
		room.Broadcast(client.nickname + ": " + msg)
	}
}
