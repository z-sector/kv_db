package initialization

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"kv_db/config"
)

func TestInitializer(t *testing.T) {
	t.Parallel()

	initializer, err := NewInitializer(config.Config{}, io.Discard)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err = initializer.Start(ctx)
	require.NoError(t, err)
}

func TestFailedInitializerCreation(t *testing.T) {
	t.Parallel()

	cfg := config.Config{Logging: config.LoggingConfig{Level: "incorrect"}}
	initializer, err := NewInitializer(cfg, io.Discard)
	require.Error(t, err)
	require.Nil(t, initializer)

	cfg = config.Config{Engine: config.EngineConfig{Type: "incorrect"}}
	initializer, err = NewInitializer(cfg, io.Discard)
	require.Error(t, err)
	require.Nil(t, initializer)

	initializer, err = NewInitializer(config.Config{WAL: &config.WALConfig{MaxSegmentSize: "10AB"}}, io.Discard)
	require.Error(t, err)
	require.Nil(t, initializer)

	cfg = config.Config{Network: config.NetworkConfig{MaxConnections: -1}}
	initializer, err = NewInitializer(cfg, io.Discard)
	require.Error(t, err)
	require.Nil(t, initializer)
}
