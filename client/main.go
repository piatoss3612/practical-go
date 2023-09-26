package main

import (
	"io"
	"log/slog"
	"net"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered", slog.Any("panic", r))
		}
	}()

	client, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("TCP connection established", slog.String("address", client.LocalAddr().String()))

	buf := make([]byte, 1024)

	for i := 0; i < 10; i++ {
		_, err := client.Write([]byte("Hello, world!"))
		if err != nil {
			slog.Error(err.Error())
			return
		}
		slog.Info("Sent message", slog.String("message", "Hello, world!"))

		n, err := client.Read(buf)
		if err != nil {
			if err != io.EOF {
				slog.Error(err.Error())
			}
			return
		}

		slog.Info("Received message", slog.String("message", string(buf[:n])))

		time.Sleep(time.Second)
	}

	client.Close()
	slog.Info("TCP connection closed", slog.String("address", client.LocalAddr().String()))
}
