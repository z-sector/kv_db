package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"kv_db/pkg/dlog"
	"kv_db/pkg/dsem"
)

type TCPHandlerFunc = func(context.Context, []byte) []byte

type TCPServer struct {
	address     string
	semaphore   dsem.Semaphore
	idleTimeout time.Duration
	logger      *slog.Logger
}

func NewTCPServer(
	address string, maxConnectionsNumber int, idleTimeout time.Duration, logger *slog.Logger,
) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("tcp server logger is invalid")
	}

	if maxConnectionsNumber <= 0 {
		return nil, errors.New("invalid number of max connections for tcp server")
	}

	return &TCPServer{
		address:     address,
		semaphore:   dsem.NewSemaphoreChan(maxConnectionsNumber),
		idleTimeout: idleTimeout,
		logger:      logger,
	}, nil
}

func (s *TCPServer) HandleQueries(ctx context.Context, handler TCPHandlerFunc) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			connection, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("failed to accept", dlog.ErrAttr(err))
				continue
			}

			s.semaphore.Acquire()
			wg.Add(1)
			go func(connection net.Conn) {
				defer func() {
					s.semaphore.Release()
					wg.Done()
				}()

				s.handleConnection(ctx, connection, handler)
			}(connection)
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		if err := listener.Close(); err != nil {
			s.logger.Warn("failed to close listener", dlog.ErrAttr(err))
		}
	}()

	wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler TCPHandlerFunc) {
	defer func() {
		if err := connection.Close(); err != nil {
			s.logger.Warn("failed to close connection", dlog.ErrAttr(err))
		}
	}()

	if err := connection.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
		s.logger.Warn("failed to set deadline", dlog.ErrAttr(err))
		return
	}

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		query := scanner.Bytes()

		response := handler(ctx, query)
		response = append(response, fmt.Sprintf("\n%s\n", EndDelim)...)
		if _, err := connection.Write(response); err != nil {
			s.logger.Warn("failed to write", dlog.ErrAttr(err))
			break
		}

		if err := connection.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			s.logger.Warn("failed to set deadline", dlog.ErrAttr(err))
			return
		}
	}
	if err := scanner.Err(); err != nil {
		s.logger.Warn("failed to scan", dlog.ErrAttr(err))
	}
}
