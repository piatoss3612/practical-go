package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered", slog.Any("panic", r))
		}
	}()

	args := os.Args

	pingCmd := newPingCmd()
	echoCmd := newEchoCmd()

	if len(args) < 2 {
		printUsage(pingCmd, echoCmd)
		os.Exit(1)
	}

	switch args[1] {
	case "ping":
		ping(pingCmd)
	case "echo":
		echo(echoCmd)
	default:
		printUsage(pingCmd, echoCmd)
	}
}

func ping(flags *flag.FlagSet) {
	cnt := flags.Int("c", 1, "Number of pings")
	flags.Parse(os.Args[2:])

	client := newClient()
	defer client.Close()

	for i := 0; i < *cnt; i++ {
		_, err := client.Write([]byte("ping"))
		if err != nil {
			panic(err)
		}
		slog.Info("Ping sent", slog.Int("count", i+1))

		buf := make([]byte, 4)
		_, err = client.Read(buf)
		if err != nil {
			panic(err)
		}

		slog.Info("Ping received", slog.String("message", string(buf)))

		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Ping finished")
}

func echo(flags *flag.FlagSet) {
	cnt := flags.Int("c", 1, "Number of echoes")
	msg := flags.String("m", "echo", "Message to echo")
	flags.Parse(os.Args[2:])

	client := newClient()
	defer client.Close()

	for i := 0; i < *cnt; i++ {
		_, err := client.Write([]byte(*msg))
		if err != nil {
			panic(err)
		}
		slog.Info("Echo sent", slog.Int("count", i+1))

		buf := make([]byte, len(*msg))
		_, err = client.Read(buf)
		if err != nil {
			panic(err)
		}

		slog.Info("Echo received", slog.String("message", string(buf)))

		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Echo finished")
}

func newClient() net.Conn {
	client, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	slog.Info("TCP connection established", slog.String("address", client.LocalAddr().String()))

	return client
}

func newPingCmd() *flag.FlagSet {
	pingCmd := flag.NewFlagSet("ping", flag.ExitOnError)
	pingCmd.Usage = func() {
		fmt.Println("Usage: client ping [options]")
		fmt.Println("Options:")
		fmt.Println("  -c <count>  Number of pings")
	}
	return pingCmd
}

func newEchoCmd() *flag.FlagSet {
	echoCmd := flag.NewFlagSet("echo", flag.ExitOnError)
	echoCmd.Usage = func() {
		fmt.Println("Usage: client echo [options]")
		fmt.Println("Options:")
		fmt.Println("  -c <count>  Number of echoes")
		fmt.Println("  -m <message>  Message to echo")
	}
	return echoCmd
}

func printUsage(flagSets ...*flag.FlagSet) {
	fmt.Println("Simple TCP client CLI tool for testing TCP server")
	for _, flagSet := range flagSets {
		fmt.Println()
		flagSet.Usage()
	}
}
