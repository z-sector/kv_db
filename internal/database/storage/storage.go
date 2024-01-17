package storage

import (
	"context"
	"errors"
	"log/slog"

	"kv_db/internal/database/comd"
	"kv_db/internal/database/storage/wal"
	"kv_db/pkg/dfuture"
)

type Engine interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, bool, error)
	Delete(context.Context, string) error
}

type WAL interface {
	Start()
	Recover() ([]wal.LogData, error)
	Set(context.Context, string, string) dfuture.FutureError
	Del(context.Context, string) dfuture.FutureError
	Shutdown()
}

type Storage struct {
	engine Engine
	wal    WAL
	logger *slog.Logger
}

func NewStorage(engine Engine, wal WAL, logger *slog.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("storage engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("storage logger is invalid")
	}

	storage := &Storage{engine: engine, wal: wal, logger: logger}

	if wal != nil {
		logs, err := wal.Recover()
		if err != nil {
			logger.Error("failed to recover database from WAL")
		}

		if err := storage.applyLogs(logs); err != nil {
			return nil, err
		}
		wal.Start()
	}

	return storage, nil
}

func MustStorage(engine Engine, wal WAL, logger *slog.Logger) *Storage {
	storage, err := NewStorage(engine, wal, logger)
	if err != nil {
		panic(err)
	}
	return storage
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	if s.wal != nil {
		future := s.wal.Set(ctx, key, value)
		if err := future.Get(); err != nil {
			return err
		}
	}
	return s.engine.Set(ctx, key, value)
}

func (s *Storage) Get(ctx context.Context, key string) (string, bool, error) {
	return s.engine.Get(ctx, key)
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if s.wal != nil {
		future := s.wal.Del(ctx, key)
		if err := future.Get(); err != nil {
			return err
		}
	}
	return s.engine.Delete(ctx, key)
}

func (s *Storage) applyLogs(logs []wal.LogData) error {
	var err error
	for _, log := range logs {
		switch log.CommandID { // nolint: exhaustive
		case comd.SetCommandID:
			if len(log.Arguments) == 2 {
				err = s.engine.Set(context.Background(), log.Arguments[0], log.Arguments[1])
			} else {
				err = errors.New("invalid arguments for set command")
			}
		case comd.DelCommandID:
			if len(log.Arguments) == 1 {
				err = s.engine.Delete(context.Background(), log.Arguments[0])
			} else {
				err = errors.New("invalid arguments for delete command")
			}
		default:
			err = errors.New("incorrect command")
		}
		if err != nil {
			return err
		}
	}
	return nil
}
