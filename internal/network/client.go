package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type TCPClient struct {
	connection  net.Conn
	idleTimeout time.Duration
}

func NewTCPClient(address string, idleTimeout time.Duration) (*TCPClient, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &TCPClient{
		connection:  connection,
		idleTimeout: idleTimeout,
	}, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if err := c.connection.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := c.connection.Write(request); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(c.connection)
	res := make([]byte, 0)

	for scanner.Scan() {
		line := scanner.Bytes()
		if strings.TrimSpace(string(line)) == EndDelim {
			break
		}
		res = append(res, line...)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return res, nil
}

func (c *TCPClient) Close() error {
	return c.connection.Close()
}
