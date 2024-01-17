package initialization

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"kv_db/config"
	"kv_db/pkg/dlog"
)

func TestCreateWALWithoutConfig(t *testing.T) {
	t.Parallel()

	wal, err := CreateWAL(nil, dlog.NewNonSlog())
	require.NoError(t, err)
	require.Nil(t, wal)
}

func TestCreateWALWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	wal, err := CreateWAL(&config.WALConfig{}, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, wal)
}

func TestCreateWALWithIncorrectSegmentSize(t *testing.T) {
	t.Parallel()

	wal, err := CreateWAL(&config.WALConfig{MaxSegmentSize: "100AB"}, dlog.NewNonSlog())
	require.Error(t, err, "max segment size is incorrect")
	require.Nil(t, wal)
}

func TestCreateWAL(t *testing.T) {
	t.Parallel()

	cfg := &config.WALConfig{
		FlushingBatchLength:  200,
		FlushingBatchTimeout: 20 * time.Millisecond,
		MaxSegmentSize:       "20MB",
		DataDirectory:        "/data/wal",
	}

	wal, err := CreateWAL(cfg, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, wal)
}
