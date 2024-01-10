package initialization

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/config"
)

func TestCreateLoggerWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	logger, err := CreateLogger(config.LoggingConfig{}, io.Discard)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateLoggerWithIncorrectLevel(t *testing.T) {
	t.Parallel()

	logger, err := CreateLogger(config.LoggingConfig{Level: "incorrect"}, io.Discard)
	require.Error(t, err)
	require.Nil(t, logger)
}

func TestCreateLogger(t *testing.T) {
	t.Parallel()

	cfg := config.LoggingConfig{
		Level:  "debug",
		Output: "test_output.log",
	}

	logger, err := CreateLogger(cfg, io.Discard)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestCreateJSONLogger(t *testing.T) {
	t.Parallel()

	cfg := config.LoggingConfig{
		JSON: true,
	}

	logger, err := CreateLogger(cfg, io.Discard)
	require.NoError(t, err)
	require.NotNil(t, logger)
	require.IsType(t, &slog.JSONHandler{}, logger.Handler())
}

func TestCreateLogFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "test_log.txt")

	file, err := CreateLogFile(path)

	require.NoError(t, err)
	require.NotNil(t, file)
}
