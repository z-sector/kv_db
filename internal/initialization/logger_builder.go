package initialization

import (
	"errors"
	"io"
	"log/slog"
	"os"

	"kv_db/config"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

var supportedLoggingLevels = map[string]slog.Level{
	debugLevel: slog.LevelDebug,
	infoLevel:  slog.LevelInfo,
	warnLevel:  slog.LevelWarn,
	errorLevel: slog.LevelError,
}

const (
	defaultLevel      = slog.LevelInfo
	defaultOutputPath = "output.log"
)

func CreateLogger(cfg config.LoggingConfig, w io.Writer) (*slog.Logger, error) {
	level := defaultLevel

	if cfg.Level != "" {
		var found bool
		if level, found = supportedLoggingLevels[cfg.Level]; !found {
			return nil, errors.New("logging level is incorrect")
		}
	}
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if cfg.JSON == nil || *cfg.JSON {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	log := slog.New(handler)

	return log, nil
}

func CreateLogFile(path string) (*os.File, error) {
	filePath := defaultOutputPath
	if path != "" {
		filePath = path
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
