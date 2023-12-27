package dlog

import "log/slog"

func NewNonSlog() *slog.Logger {
	return slog.New(discardHandler{})
}
