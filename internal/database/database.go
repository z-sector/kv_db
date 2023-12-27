package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type ComputeLayer interface {
	HandleQuery(context.Context, string) (Query, error)
}

type StorageLayer interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, bool, error)
	Del(ctx context.Context, key string) error
}

type Database struct {
	computeLayer ComputeLayer
	storageLayer StorageLayer
	logger       *slog.Logger
}

func NewDatabase(computeLayer ComputeLayer, storageLayer StorageLayer, logger *slog.Logger) (*Database, error) {
	if computeLayer == nil {
		return nil, errors.New("database compute is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("database storage is invalid")
	}

	if logger == nil {
		return nil, errors.New("database logger is invalid")
	}

	return &Database{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		logger:       logger,
	}, nil
}

func MustDatabase(computeLayer ComputeLayer, storageLayer StorageLayer, logger *slog.Logger) *Database {
	database, err := NewDatabase(computeLayer, storageLayer, logger)
	if err != nil {
		panic(err)
	}
	return database
}

func (d *Database) HandleQuery(ctx context.Context, queryStr string) string {
	query, err := d.computeLayer.HandleQuery(ctx, queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch query.CommandID() {
	case SetCommandID:
		return d.handleSetQuery(ctx, query)
	case GetCommandID:
		return d.handleGetQuery(ctx, query)
	case DelCommandID:
		return d.handleDelQuery(ctx, query)
	case UnknownCommandID:
		d.logger.Error("compute layer is incorrect")
	}

	return "[error] internal configuration error"
}

func (d *Database) handleSetQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}

func (d *Database) handleGetQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	value, ok, err := d.storageLayer.Get(ctx, arguments[0])
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}
	if !ok {
		return "[nil]"
	}

	return fmt.Sprintf("[ok] %s", value)
}

func (d *Database) handleDelQuery(ctx context.Context, query Query) string {
	arguments := query.Arguments()
	if err := d.storageLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}
