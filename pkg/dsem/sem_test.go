package dsem

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSemaphoreLimit1(t *testing.T) {
	t.Parallel()
	limit := 1

	t.Run("channel semaphore", func(t *testing.T) {
		testWithLimit1(t, NewSemaphoreChan(limit))
	})

	t.Run("cond semaphore", func(t *testing.T) {
		testWithLimit1(t, NewSemaphoreCond(limit))
	})
}

func TestSemaphoreLimit2(t *testing.T) {
	t.Parallel()

	limit := 2

	t.Run("channel semaphore", func(t *testing.T) {
		testWithLimit2(t, NewSemaphoreChan(limit))
	})

	t.Run("cond semaphore", func(t *testing.T) {
		testWithLimit2(t, NewSemaphoreCond(limit))
	})
}

func testWithLimit1(t *testing.T, semaphore Semaphore) {
	t.Helper()

	goroutinesCount := 20
	iters := 1000

	wg := sync.WaitGroup{}
	sharedCounter := 0

	testFunc := func(sem Semaphore) {
		defer wg.Done()
		for j := 0; j < iters; j++ {
			sem.WithSemaphore(func() {
				sharedCounter++
			})
		}
	}

	wg.Add(goroutinesCount)
	for i := 0; i < goroutinesCount; i++ {
		go testFunc(semaphore)
	}
	wg.Wait()

	require.Equal(t, iters*goroutinesCount, sharedCounter)
}

func testWithLimit2(t *testing.T, semaphore Semaphore) {
	t.Helper()

	nonBlockingAcquire := func() <-chan struct{} {
		ch := make(chan struct{})
		go func() {
			semaphore.Acquire()
			ch <- struct{}{}
		}()
		return ch
	}

	semaphore.Acquire()
	semaphore.Acquire()

	res := nonBlockingAcquire()

	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()

	var isError bool
	select {
	case <-timer.C:
		isError = false
		semaphore.Release()
	case <-res:
		isError = true
	}
	require.False(t, isError)
}
