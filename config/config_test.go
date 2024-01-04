package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("non existent file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "non_test.yaml")

		cfg, err := Load(src)

		require.Error(t, err)
		require.Nil(t, cfg)
	})

	t.Run("empty filename", func(t *testing.T) {
		cfg, err := Load("")

		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("empty file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "empty_test.yaml")
		err := os.WriteFile(src, []byte(""), 0o666)
		require.NoError(t, err)

		cfg, err := Load(src)

		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("invalid file", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "invalid_test.yaml")
		err := os.WriteFile(src, []byte("123"), 0o666)
		require.NoError(t, err)

		cfg, err := Load(src)

		require.Error(t, err)
		require.Nil(t, cfg)
	})

	t.Run("load config from example", func(t *testing.T) {
		src := "example_config.yaml"

		cfg, err := Load(src)

		require.NoError(t, err)

		require.Equal(t, "in_memory", cfg.Engine.Type)

		require.Equal(t, "127.0.0.1:3223", cfg.Network.Address)
		require.Equal(t, 100, cfg.Network.MaxConnections)

		require.Equal(t, "info", cfg.Logging.Level)
		require.Equal(t, "/log/output.log", cfg.Logging.Output)
	})
}
