package analyzer

import (
	"context"
	"errors"
	"log/slog"

	"kv_db/internal/database"
	"kv_db/internal/database/compute"
)

type QueryAnalyzer struct {
	validators []func(database.Query) error
	logger     *slog.Logger
}

func NewAnalyzer(logger *slog.Logger) (*QueryAnalyzer, error) {
	if logger == nil {
		return nil, errors.New("queryAnalyzer logger is invalid")
	}

	analyser := &QueryAnalyzer{
		logger: logger,
	}

	analyser.validators = []func(database.Query) error{
		database.SetCommandID: validateArgsCount(2),
		database.GetCommandID: validateArgsCount(1),
		database.DelCommandID: validateArgsCount(1),
	}

	return analyser, nil
}

func MustAnalyzer(logger *slog.Logger) *QueryAnalyzer {
	analyser, err := NewAnalyzer(logger)
	if err != nil {
		panic(err)
	}
	return analyser
}

func (a *QueryAnalyzer) AnalyzeQuery(_ context.Context, tokens []string) (database.Query, error) {
	if len(tokens) == 0 {
		return database.Query{}, compute.ErrInvalidCommand
	}

	command := tokens[0]
	commandID := database.GetCommandIDByName(command)
	if commandID == database.UnknownCommandID {
		return database.Query{}, compute.ErrInvalidCommand
	}

	query := database.NewQuery(commandID, tokens[1:])
	validator := a.validators[commandID]
	if validator != nil {
		if err := validator(query); err != nil {
			return database.Query{}, err
		}
	}

	return query, nil
}

func validateArgsCount(expected int) func(database.Query) error {
	return func(query database.Query) error {
		if len(query.Arguments()) != expected {
			return compute.ErrInvalidArguments
		}
		return nil
	}
}
