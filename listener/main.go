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
	defer listener.Close()

	slog.Info("TCP server started", slog.String("address", listener.Addr().String()))

	done := make(chan struct{})

	go func() {
		defer func() { close(done) }()

		for {
			conn, err := listener.Accept()
			if err != nil {
				slog.Error(err.Error(), slog.String("from", "listener/main.go line:25"))
				return
			}

			go func(c net.Conn) {
				defer c.Close()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						slog.Error(err.Error(), slog.String("from", "listener/main.go line:36"))
						return
					}

					slog.Info("Received message", slog.String("message", string(buf[:n])))
				}
			}(conn)
		}
	}()

	<-done
}
