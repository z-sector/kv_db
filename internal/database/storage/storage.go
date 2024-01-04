package storage

import (
	"context"
	"errors"
	"log/slog"
)

type Engine interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, bool, error)
	Delete(context.Context, string) error
}

type Storage struct {
	engine Engine
	logger *slog.Logger
}

func NewStorage(engine Engine, logger *slog.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("storage engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("storage logger is invalid")
	}
	return &Storage{engine: engine, logger: logger}, nil
}

func MustStorage(engine Engine, logger *slog.Logger) *Storage {
	storage, err := NewStorage(engine, logger)
	if err != nil {
		panic(err)
	}
	return storage
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	return s.engine.Set(ctx, key, value)
}

func (s *Storage) Get(ctx context.Context, key string) (string, bool, error) {
	return s.engine.Get(ctx, key)
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.engine.Delete(ctx, key)
}
