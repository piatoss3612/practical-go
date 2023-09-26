package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	client, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	quit := make(chan bool)
	defer close(quit)

	go recvMsg(client)
	go sendMsg(client, quit)

	<-quit

	fmt.Println("Bye!")
}

func sendMsg(conn net.Conn, quit chan<- bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for {
		if !scanner.Scan() {
			if errors.Is(scanner.Err(), io.EOF) {
				continue
			}

			slog.Error(scanner.Err().Error())
			os.Exit(2)
		}

		msg := scanner.Bytes()

		if len(msg) == 0 {
			continue
		}

		if string(msg) == "quit" {
			quit <- true
			return
		}

		_, err := conn.Write(msg)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(2)
		}
	}
}

func recvMsg(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {

				return
			}

			slog.Error(err.Error())
			os.Exit(2)
		}

		msg := string(buf[:n])

		fmt.Println(msg)
	}
}
