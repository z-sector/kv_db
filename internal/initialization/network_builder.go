package initialization

import (
	"log/slog"
	"time"

	"kv_db/config"
	"kv_db/internal/network"
)

const (
	defaultServerAddress       = "localhost:3223"
	defaultMaxConnectionNumber = 100
	defaultIdleTimeout         = time.Minute * 5
)

func CreateNetwork(cfg config.NetworkConfig, logger *slog.Logger) (*network.TCPServer, error) {
	address := defaultServerAddress
	maxConnectionsNumber := defaultMaxConnectionNumber
	idleTimeout := defaultIdleTimeout

	if cfg.Address != "" {
		address = cfg.Address
	}

	if cfg.MaxConnections != 0 {
		maxConnectionsNumber = cfg.MaxConnections
	}

	return network.NewTCPServer(address, maxConnectionsNumber, idleTimeout, logger)
}
