package parser

import (
	"context"
	"errors"
	"log/slog"
)

type QueryParser struct {
	logger *slog.Logger
}

func NewParser(logger *slog.Logger) (*QueryParser, error) {
	if logger == nil {
		return nil, errors.New("queryParser logger is invalid")
	}

	return &QueryParser{logger: logger}, nil
}

func MustParser(logger *slog.Logger) *QueryParser {
	parser, err := NewParser(logger)
	if err != nil {
		panic(err)
	}
	return parser
}

func (p *QueryParser) ParseQuery(_ context.Context, query string) ([]string, error) {
	machine := NewStateMachine()
	tokens, err := machine.Parse(query)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
