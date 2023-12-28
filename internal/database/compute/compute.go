package compute

import (
	"context"
	"errors"
	"log/slog"

	"kv_db/internal/database"
)

type Parser interface {
	ParseQuery(context.Context, string) ([]string, error)
}

type Analyzer interface {
	AnalyzeQuery(context.Context, []string) (database.Query, error)
}

type Compute struct {
	parser   Parser
	analyzer Analyzer
	logger   *slog.Logger
}

func NewCompute(parser Parser, analyzer Analyzer, logger *slog.Logger) (*Compute, error) {
	if parser == nil {
		return nil, errors.New("compute parser is invalid")
	}

	if analyzer == nil {
		return nil, errors.New("compute analyzer is invalid")
	}

	if logger == nil {
		return nil, errors.New("compute logger is invalid")
	}

	return &Compute{
		parser:   parser,
		analyzer: analyzer,
		logger:   logger,
	}, nil
}

func MustCompute(parser Parser, analyzer Analyzer, logger *slog.Logger) *Compute {
	compute, err := NewCompute(parser, analyzer, logger)
	if err != nil {
		panic(err)
	}
	return compute
}

func (d *Compute) HandleQuery(ctx context.Context, queryStr string) (database.Query, error) {
	tokens, err := d.parser.ParseQuery(ctx, queryStr)
	if err != nil {
		return database.Query{}, err
	}

	query, err := d.analyzer.AnalyzeQuery(ctx, tokens)
	if err != nil {
		return database.Query{}, err
	}

	return query, nil
}
