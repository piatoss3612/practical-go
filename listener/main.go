package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	slog.Info("TCP server started", slog.String("address", listener.Addr().String()))
	defer func() {
		slog.Info("TCP server stopped", slog.String("address", listener.Addr().String()))
	}()

	<-run(listener)
}

func run(listener net.Listener) <-chan struct{} {
	pingChan := make(chan net.Conn)
	echoChan := make(chan echoConn)

	close := func() {
		close(pingChan)
		close(echoChan)
		_ = listener.Close()
	}

	go pingHandler(pingChan)
	go echoHandler(echoChan)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					slog.Error(err.Error())
				}
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

					switch string(buf[:n]) {
					case "ping":
						pingChan <- c
					default:
						echoChan <- echoConn{c, buf[:n]}
					}
				}
			}(conn)
		}
	}()

	return gracefulShutdown(close)
}

func gracefulShutdown(fn func()) <-chan struct{} {
	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	defer func() {
		<-sig
		fn()
		close(sig)
		close(done)
	}()

	return done
}

func pingHandler(conns <-chan net.Conn) {
	for conn := range conns {
		_, err := conn.Write([]byte("pong"))
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

type echoConn struct {
	net.Conn
	msg []byte
}

func echoHandler(conns <-chan echoConn) {
	for conn := range conns {
		_, err := conn.Write(conn.msg)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}
