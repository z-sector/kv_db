package wal

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database/comd"
	"kv_db/pkg/dlog"
)

const testWALDirectory = "temp_test_data"

func TestMain(m *testing.M) {
	if err := os.Mkdir(testWALDirectory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	if err := os.RemoveAll(testWALDirectory); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestBatchWritingToWALSegment(t *testing.T) {
	maxSegmentSize := 100 << 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, dlog.NewNonSlog())

	batch := []Log{
		NewLog(1, comd.SetCommandID, []string{"key_1", "value_1"}),
		NewLog(2, comd.SetCommandID, []string{"key_2", "value_2"}),
		NewLog(3, comd.SetCommandID, []string{"key_3", "value_3"}),
	}

	now = func() time.Time {
		return time.Unix(1, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/wal_1000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}

func TestWALSegmentsRotation(t *testing.T) {
	maxSegmentSize := 10
	fsWriter := NewFSWriter(testWALDirectory, maxSegmentSize, dlog.NewNonSlog())

	batch := []Log{
		NewLog(4, comd.SetCommandID, []string{"key_4", "value_4"}),
		NewLog(5, comd.SetCommandID, []string{"key_5", "value_5"}),
		NewLog(6, comd.SetCommandID, []string{"key_6", "value_6"}),
	}

	now = func() time.Time {
		return time.Unix(2, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	batch = []Log{
		NewLog(7, comd.SetCommandID, []string{"key_7", "value_7"}),
		NewLog(8, comd.SetCommandID, []string{"key_8", "value_8"}),
		NewLog(9, comd.SetCommandID, []string{"key_9", "value_9"}),
	}

	now = func() time.Time {
		return time.Unix(3, 0)
	}

	fsWriter.WriteBatch(batch)
	for _, record := range batch {
		err := record.Result()
		require.NoError(t, err.Get())
	}

	stat, err := os.Stat(testWALDirectory + "/wal_2000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())

	stat, err = os.Stat(testWALDirectory + "/wal_3000.log")
	require.NoError(t, err)
	require.NotZero(t, stat.Size())
}
