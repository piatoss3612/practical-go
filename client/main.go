package main

import (
	"log/slog"
	"net"
	"time"
)

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("TCP connection established", slog.String("address", client.LocalAddr().String()))

	for i := 0; i < 10; i++ {
		client.Write([]byte("Hello, world!"))
		time.Sleep(time.Second)
	}

	client.Close()
	slog.Info("TCP connection closed", slog.String("address", client.LocalAddr().String()))
}
