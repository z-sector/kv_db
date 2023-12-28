package storage

import (
	"context"
	"errors"
	"log/slog"
)

type Backend interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, bool, error)
	Delete(context.Context, string) error
}

type Storage struct {
	backend Backend
	logger  *slog.Logger
}

func NewStorage(backend Backend, logger *slog.Logger) (*Storage, error) {
	if backend == nil {
		return nil, errors.New("storage backend is invalid")
	}

	if logger == nil {
		return nil, errors.New("storage logger is invalid")
	}
	return &Storage{backend: backend, logger: logger}, nil
}

func MustStorage(backend Backend, logger *slog.Logger) *Storage {
	storage, err := NewStorage(backend, logger)
	if err != nil {
		panic(err)
	}
	return storage
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	return s.backend.Set(ctx, key, value)
}

func (s *Storage) Get(ctx context.Context, key string) (string, bool, error) {
	return s.backend.Get(ctx, key)
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.backend.Delete(ctx, key)
}
