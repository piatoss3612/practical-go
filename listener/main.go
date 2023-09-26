package main

import (
	"io"
	"log/slog"
	"net"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered", slog.Any("panic", r))
		}
	}()

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
				slog.Error(err.Error())
				return
			}

			slog.Info("TCP connection established", slog.String("address", conn.RemoteAddr().String()))

			go func(c net.Conn) {
				defer func() {
					_ = c.Close()
					slog.Info("TCP connection closed", slog.String("address", conn.RemoteAddr().String()))
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							slog.Error(err.Error())
						}
						return
					}

					slog.Info("Received message", slog.String("message", string(buf[:n])))

					_, err = c.Write(buf[:n])
					if err != nil {
						slog.Error(err.Error())
						return
					}

					slog.Info("Echoed message", slog.String("message", string(buf[:n])))
				}
			}(conn)
		}
	}()

	<-done
}
