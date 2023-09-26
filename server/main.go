package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
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
	room     *Room
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
		time.Sleep(100 * time.Millisecond)
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
			if len(fields) < 2 {
				_, _ = client.Write([]byte("Missing nickname"))
				continue
			}

			changeNickname(client, fields[1])
		case "list":
			listRooms(client)
		case "quit":
			return
		case "help":
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

	client.room = room
	if !room.Join(client) {
		_, _ = client.Write([]byte("Room is full"))
		return
	}
	room.Broadcast(client.nickname + " joined the room")

	go chatInRoom(client)
}

func chatInRoom(client *Client) {
	buf := make([]byte, 1024)

	for {
		n, err := client.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				client.room.Leave(client)
				client.room.Broadcast(client.nickname + " left the room...")
				return
			}
			slog.Error(err.Error())
			return
		}

		msg := string(buf[:n])

		if msg == ":quit" {
			client.room.Leave(client)
			client.room.Broadcast(client.nickname + " left the room...")
			handleClient(client)
			return
		}

		client.room.Broadcast(client.nickname + ": " + msg)
	}
}

func changeNickname(client *Client, nickname string) {
	client.nickname = nickname
	_, _ = client.Write([]byte("Nickname changed to " + client.nickname))
}

func listRooms(client *Client) {
	mu.Lock()
	defer mu.Unlock()

	if len(rooms) == 0 {
		_, _ = client.Write([]byte("No rooms"))
		return
	}

	var names []string

	for name := range rooms {
		names = append(names, name)
	}

	_, _ = client.Write([]byte(strings.Join(names, "\n")))
}
