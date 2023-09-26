package main

import (
	"log/slog"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("Chat server started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {}
