package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
)

var ListenerAddr = "127.0.0.1:8080"

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered", slog.Any("panic", r))
		}
	}()

	listener, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			return
		}

		slog.Info("Accepted connection", slog.Any("remote_addr", conn.RemoteAddr()))

		go func(from net.Conn) {
			to, err := net.Dial("tcp", ListenerAddr)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			defer to.Close()

			if err := proxy(from, to); err != nil {
				if !errors.Is(err, io.EOF) {
					slog.Error(err.Error())
				}
			}
		}(conn)
	}
}

func proxy(from, to net.Conn) error {
	defer func() {
		_ = from.Close()
		_ = to.Close()
	}()

	go func() {
		_, _ = io.Copy(from, to)
	}()

	_, err := io.Copy(to, from)
	return err
}
