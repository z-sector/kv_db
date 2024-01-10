package initialization

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"kv_db/config"
	"kv_db/internal/database"
	"kv_db/internal/database/compute"
	"kv_db/internal/database/compute/analyzer"
	"kv_db/internal/database/compute/parser"
	"kv_db/internal/database/storage"
	"kv_db/internal/network"
	"kv_db/pkg/dlog"
)

type Initializer struct {
	engine storage.Engine
	server *network.TCPServer
	logger *slog.Logger
}

func NewInitializer(cfg config.Config, logW io.Writer) (*Initializer, error) {
	logger, err := CreateLogger(cfg.Logging, logW)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	dbEngine, err := CreateEngine(cfg.Engine)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	tcpServer, err := CreateNetwork(
		cfg.Network, logger.With(slog.String("layer", "server")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	initializer := &Initializer{
		engine: dbEngine,
		server: tcpServer,
		logger: logger,
	}

	return initializer, nil
}

func (i *Initializer) Start(ctx context.Context) error {
	computeLayer, err := i.createComputeLayer()
	if err != nil {
		return err
	}

	storageLayer, err := i.createStorageLayer()
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(
		computeLayer,
		storageLayer,
		i.logger.With(slog.String("layer", "database")),
	)
	if err != nil {
		i.logger.Error("failed to start database", dlog.ErrAttr(err))
		return err
	}

	return i.server.HandleQueries(ctx, func(ctx context.Context, query []byte) []byte {
		response := db.HandleQuery(ctx, string(query))
		return []byte(response)
	})
}

func (i *Initializer) createComputeLayer() (*compute.Compute, error) {
	queryParser, err := parser.NewParser(
		i.logger.With(slog.String("layer", "parser")),
	)
	if err != nil {
		i.logger.Error("failed to initialize parser", dlog.ErrAttr(err))
		return nil, err
	}

	queryAnalyzer, err := analyzer.NewAnalyzer(
		i.logger.With(slog.String("layer", "analyzer")),
	)
	if err != nil {
		i.logger.Error("failed to initialize analyzer", dlog.ErrAttr(err))
		return nil, err
	}

	computeLayer, err := compute.NewCompute(
		queryParser, queryAnalyzer, i.logger.With(slog.String("layer", "compute")),
	)
	if err != nil {
		i.logger.Error("failed to initialize compute layer", dlog.ErrAttr(err))
		return nil, err
	}

	return computeLayer, nil
}

func (i *Initializer) createStorageLayer() (*storage.Storage, error) {
	storageLayer, err := storage.NewStorage(
		i.engine, i.logger.With(slog.String("layer", "storage")),
	)
	if err != nil {
		i.logger.Error("failed to initialize storage layer", dlog.ErrAttr(err))
		return nil, err
	}

	return storageLayer, nil
}
