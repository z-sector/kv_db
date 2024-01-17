package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("non existent file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "non_test.yaml")

		_, err := Load(src)

		require.Error(t, err)
	})

	t.Run("empty filename", func(t *testing.T) {
		_, err := Load("")

		require.NoError(t, err)
	})

	t.Run("empty file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "empty_test.yaml")
		err := os.WriteFile(src, []byte(""), 0o666)
		require.NoError(t, err)

		_, err = Load(src)

		require.NoError(t, err)
	})

	t.Run("invalid file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "invalid_test.yaml")
		err := os.WriteFile(src, []byte("123"), 0o666)
		require.NoError(t, err)

		_, err = Load(src)

		require.Error(t, err)
	})

	t.Run("load config from example", func(t *testing.T) {
		src := "example_config.yaml"

		cfg, err := Load(src)

		require.NoError(t, err)

		require.Equal(t, "in_memory", cfg.Engine.Type)

		require.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
		require.Equal(t, 100, cfg.Network.MaxConnections)

		require.Equal(t, false, *cfg.Logging.JSON)
		require.Equal(t, "info", cfg.Logging.Level)
		require.Equal(t, "/log/output.log", cfg.Logging.Output)

		require.Equal(t, 100, cfg.WAL.FlushingBatchLength)
		require.Equal(t, 10*time.Millisecond, cfg.WAL.FlushingBatchTimeout)
		require.Equal(t, "10MB", cfg.WAL.MaxSegmentSize)
		require.Equal(t, "/data/wal", cfg.WAL.DataDirectory)
	})
}
