package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"kv_db/internal/network"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	address := flag.String("address", "localhost:3223", "Address of the kv_db")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(*address, *idleTimeout)
	if err != nil {
		return err
	}

	for {
		fmt.Print("[kv_db] > ") // nolint:forbidigo
		request, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return fmt.Errorf("connection was closed: %w", err)
			}
			return fmt.Errorf("failed to read user query: %w", err)
		}

		response, err := client.Send([]byte(request))
		if err != nil {
			return err
		}

		fmt.Println(string(response)) // nolint:forbidigo
	}
}
