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
	proxyCmd := newProxyCmd()

	if len(args) < 2 {
		printUsage(pingCmd, echoCmd)
		os.Exit(1)
	}

	switch args[1] {
	case "ping":
		runPing(pingCmd)
	case "echo":
		runEcho(echoCmd)
	case "proxy":
		runProxy(proxyCmd)
	default:
		printUsage(pingCmd, echoCmd)
	}
}

func runPing(flags *flag.FlagSet) {
	cnt := flags.Int("c", 1, "Number of pings")
	flags.Parse(os.Args[2:])

	ping("127.0.0.1:8080", *cnt)
}

func runEcho(flags *flag.FlagSet) {
	cnt := flags.Int("c", 1, "Number of echoes")
	msg := flags.String("m", "echo", "Message to echo")
	flags.Parse(os.Args[2:])

	echo("127.0.0.1:8080", *msg, *cnt)
}

func runProxy(flags *flag.FlagSet) {
	cmd := flags.String("cmd", "", "Command to execute")
	cnt := flags.Int("c", 1, "Number of echoes")
	msg := flags.String("m", "echo", "Message to echo")
	flags.Parse(os.Args[2:])

	addr := "127.0.0.1:7070"

	switch *cmd {
	case "ping":
		ping(addr, *cnt)
	case "echo":
		echo(addr, *msg, *cnt)
	default:
		printUsage(flags)
	}
}

func ping(addr string, cnt int) {
	client := newClient(addr)
	defer client.Close()

	for i := 0; i < cnt; i++ {
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

func echo(addr, msg string, cnt int) {
	client := newClient(addr)
	defer client.Close()

	for i := 0; i < cnt; i++ {
		_, err := client.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
		slog.Info("Echo sent", slog.Int("count", i+1))

		buf := make([]byte, len(msg))
		_, err = client.Read(buf)
		if err != nil {
			panic(err)
		}

		slog.Info("Echo received", slog.String("message", string(buf)))

		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Echo finished")
}

func newClient(addr string) net.Conn {
	client, err := net.Dial("tcp", addr)
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

func newProxyCmd() *flag.FlagSet {
	proxyCmd := flag.NewFlagSet("proxy", flag.ExitOnError)
	proxyCmd.Usage = func() {
		fmt.Println("Usage: client proxy [options]")
		fmt.Println("Options:")
		fmt.Println("  -cmd <command>  Command to execute")
		fmt.Println("  -c <count>  Number of echoes")
		fmt.Println("  -m <message>  Message to echo")
	}
	return proxyCmd
}

func printUsage(flagSets ...*flag.FlagSet) {
	fmt.Println("Simple TCP client CLI tool for testing TCP server")
	for _, flagSet := range flagSets {
		fmt.Println()
		flagSet.Usage()
	}
}
