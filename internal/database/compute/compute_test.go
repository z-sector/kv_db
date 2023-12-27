package compute

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kv_db/internal/database"
	"kv_db/pkg/dlog"
)

func TestNewCompute(t *testing.T) {
	t.Parallel()

	parser, analyzer := getMockParserAndAnalyzer(t)

	compute, err := NewCompute(parser, analyzer, nil)
	require.Error(t, err)
	require.Nil(t, compute)

	compute, err = NewCompute(parser, nil, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, compute)

	compute, err = NewCompute(nil, analyzer, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, compute)

	compute, err = NewCompute(parser, analyzer, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, compute)
}

func TestMustCompute(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustCompute(nil, nil, nil)
	})

	require.NotPanics(t, func() {
		parser, analyzer := getMockParserAndAnalyzer(t)

		MustCompute(parser, analyzer, dlog.NewNonSlog())
	})
}

func TestCompute_HandleQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		parser, analyzer := getMockParserAndAnalyzer(t)
		compute, err := NewCompute(parser, analyzer, dlog.NewNonSlog())
		require.NoError(t, err)
		parser.EXPECT().
			ParseQuery(ctx, gomock.Eq("GET key")).
			Return([]string{"GET", "key"}, nil)
		analyzer.EXPECT().
			AnalyzeQuery(ctx, gomock.Eq([]string{"GET", "key"})).
			Return(database.NewQuery(database.GetCommandID, []string{"key"}), nil)

		query, err := compute.HandleQuery(ctx, "GET key")

		require.NoError(t, err)
		require.Equal(t, database.NewQuery(database.GetCommandID, []string{"key"}), query)
	})

	t.Run("parser error", func(t *testing.T) {
		ctx := context.Background()
		parser, analyzer := getMockParserAndAnalyzer(t)
		compute, err := NewCompute(parser, analyzer, dlog.NewNonSlog())
		require.NoError(t, err)
		rawQuery := "SET key value"
		parser.EXPECT().
			ParseQuery(ctx, gomock.Eq(rawQuery)).
			Return(nil, ErrInvalidSymbol)

		query, err := compute.HandleQuery(ctx, rawQuery)

		require.ErrorIs(t, err, ErrInvalidSymbol)
		require.Empty(t, query)
	})

	t.Run("analyzer error", func(t *testing.T) {
		ctx := context.Background()
		parser, analyzer := getMockParserAndAnalyzer(t)
		compute, err := NewCompute(parser, analyzer, dlog.NewNonSlog())
		require.NoError(t, err)
		rawQuery := "SET key value"
		parser.EXPECT().
			ParseQuery(ctx, gomock.Eq(rawQuery)).
			Return([]string{"SET", "key", "value"}, nil)
		analyzer.EXPECT().
			AnalyzeQuery(ctx, gomock.Any()).
			Return(database.Query{}, ErrInvalidCommand)

		query, err := compute.HandleQuery(ctx, rawQuery)

		require.ErrorIs(t, err, ErrInvalidCommand)
		require.Empty(t, query)
	})
}

func getMockParserAndAnalyzer(t *testing.T) (*MockParser, *MockAnalyzer) {
	t.Helper()

	ctrl := gomock.NewController(t)
	return NewMockParser(ctrl), NewMockAnalyzer(ctrl)
}
