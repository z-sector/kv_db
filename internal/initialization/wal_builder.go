package initialization

import (
	"errors"
	"log/slog"
	"time"

	"kv_db/config"
	"kv_db/internal/database/storage"
	"kv_db/internal/database/storage/wal"
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20
	defaultWALDataDirectory     = "./data/wal"
)

func CreateWAL(cfg *config.WALConfig, logger *slog.Logger) (storage.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg == nil {
		return nil, nil
	}

	if cfg.FlushingBatchLength != 0 {
		flushingBatchSize = cfg.FlushingBatchLength
	}

	if cfg.FlushingBatchTimeout != 0 {
		flushingBatchTimeout = cfg.FlushingBatchTimeout
	}

	if cfg.MaxSegmentSize != "" {
		size, err := ParseSize(cfg.MaxSegmentSize)
		if err != nil {
			return nil, errors.New("max segment size is incorrect")
		}

		maxSegmentSize = size
	}

	if cfg.DataDirectory != "" {
		dataDirectory = cfg.DataDirectory
	}

	fsReader := wal.NewFSReader(dataDirectory, logger)
	fsWriter := wal.NewFSWriter(dataDirectory, maxSegmentSize, logger)
	return wal.NewWAL(fsWriter, fsReader, flushingBatchTimeout, flushingBatchSize), nil
}
