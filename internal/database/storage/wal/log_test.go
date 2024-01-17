package wal

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database/comd"
)

func TestLogFuture(t *testing.T) {
	args := []string{"test"}
	log := NewLog(1, comd.SetCommandID, args)

	require.Equal(t, int64(1), log.LSN())
	require.Equal(t, comd.SetCommandID, log.CommandID())
	require.Equal(t, args, log.Arguments())
	require.Equal(t, comd.SetCommandID, log.Data().CommandID)

	future := log.Result()

	wg := sync.WaitGroup{}

	nonBlockingFutureGet := func() <-chan error {
		defer wg.Done()
		ch := make(chan error, 1)
		go func() {
			ch <- future.Get()
		}()
		return ch
	}

	wg.Add(1)
	resCh := nonBlockingFutureGet()
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-resCh:
	}
	require.Equal(t, 0, len(resCh))

	expErr := errors.New("test")
	log.SetResult(expErr)
	wg.Wait()
	actErr := <-resCh

	require.ErrorIs(t, actErr, expErr)
}
