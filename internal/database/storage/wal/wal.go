package wal

import (
	"context"
	"sync"
	"time"

	"kv_db/internal/database/comd"
	"kv_db/internal/database/txctx"
	"kv_db/pkg/dfuture"
	"kv_db/pkg/dlock"
)

type fsWriter interface {
	WriteBatch([]Log)
}

type fsReader interface {
	ReadLogs() ([]LogData, error)
}

type WAL struct {
	fsWriter     fsWriter
	fsReader     fsReader
	flushTimeout time.Duration
	maxBatchSize int

	mutex   sync.Mutex
	batch   []Log
	batches chan []Log

	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewWAL(
	fsWriter fsWriter,
	fsReader fsReader,
	flushTimeout time.Duration,
	maxBatchSize int,
) *WAL {
	return &WAL{
		fsWriter:     fsWriter,
		fsReader:     fsReader,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []Log, 1),
		closeCh:      make(chan struct{}),
		closeDoneCh:  make(chan struct{}),
	}
}

func (w *WAL) Recover() ([]LogData, error) {
	return w.fsReader.ReadLogs()
}

func (w *WAL) Start() {
	go func() {
		defer func() {
			w.closeDoneCh <- struct{}{}
		}()

		tick := time.NewTicker(w.flushTimeout)
		defer tick.Stop()
		for {
			select {
			case <-w.closeCh:
				w.flushBatch()
				return
			case batch := <-w.batches:
				w.fsWriter.WriteBatch(batch)
			case <-tick.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Shutdown() {
	close(w.closeCh)
	<-w.closeDoneCh
}

func (w *WAL) Set(ctx context.Context, key, value string) dfuture.FutureError {
	return w.push(ctx, comd.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) dfuture.FutureError {
	return w.push(ctx, comd.DelCommandID, []string{key})
}

func (w *WAL) flushBatch() {
	var batch []Log
	dlock.WithLock(&w.mutex, func() {
		if len(w.batch) != 0 {
			batch = w.batch
			w.batch = nil
		}
	})

	if len(batch) != 0 {
		w.fsWriter.WriteBatch(batch)
	}
}

func (w *WAL) push(ctx context.Context, commandID comd.CmdID, args []string) dfuture.FutureError {
	txID := txctx.GetTxFromCtx(ctx)
	record := NewLog(txID, commandID, args)

	dlock.WithLock(&w.mutex, func() {
		w.batch = append(w.batch, record)
		if len(w.batch) == w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return record.Result()
}
