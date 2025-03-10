package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: go run main.go <host> <port>")
	}

	address := net.JoinHostPort(args[0], args[1])

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		cancel()
		log.Fatalf("connect failed: %v", err)
	}

	defer client.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	running(&wg, client.Send, cancel)
	running(&wg, client.Receive, cancel)

	<-ctx.Done()
	log.Println("connection closed")
	wg.Wait()
}

func running(wg *sync.WaitGroup, task func() error, cancel func()) {
	go func() {
		defer wg.Done()
		if err := task(); err != nil {
			log.Printf("task failed: %v", err)
			cancel()
		}
	}()
}
