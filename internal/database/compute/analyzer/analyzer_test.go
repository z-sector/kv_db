package analyzer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database"
	"kv_db/internal/database/comd"
	"kv_db/internal/database/compute"
	"kv_db/pkg/dlog"
)

func TestNewAnalyzer(t *testing.T) {
	t.Parallel()

	analyzer, err := NewAnalyzer(nil)
	require.Error(t, err)
	require.Nil(t, analyzer)

	analyzer, err = NewAnalyzer(dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, analyzer)
}

func TestMustAnalyzer(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustAnalyzer(nil)
	})

	require.NotPanics(t, func() {
		MustAnalyzer(dlog.NewNonSlog())
	})
}

func TestAnalyzeQuery(t *testing.T) {
	ctx := context.Background()
	analyzer, err := NewAnalyzer(dlog.NewNonSlog())
	require.NoError(t, err)

	testcases := map[string]struct {
		tokens   []string
		expQuery database.Query
		expErr   error
	}{
		"empty tokens": {
			tokens: []string{},
			expErr: compute.ErrInvalidCommand,
		},
		"invalid command": {
			tokens: []string{"TEST"},
			expErr: compute.ErrInvalidCommand,
		},
		"invalid number arguments for set query": {
			tokens: []string{"SET", "key"},
			expErr: compute.ErrInvalidArguments,
		},
		"invalid number arguments for get query": {
			tokens: []string{"GET", "key", "value"},
			expErr: compute.ErrInvalidArguments,
		},
		"invalid number arguments for del query": {
			tokens: []string{"DEL", "key", "value"},
			expErr: compute.ErrInvalidArguments,
		},
		"valid set query": {
			tokens:   []string{"SET", "key", "value"},
			expQuery: database.NewQuery(comd.SetCommandID, []string{"key", "value"}),
		},
		"valid get query": {
			tokens:   []string{"GET", "key"},
			expQuery: database.NewQuery(comd.GetCommandID, []string{"key"}),
		},
		"valid del query": {
			tokens:   []string{"DEL", "key"},
			expQuery: database.NewQuery(comd.DelCommandID, []string{"key"}),
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, err := analyzer.AnalyzeQuery(ctx, tc.tokens)
			require.ErrorIs(t, err, tc.expErr)
			require.Equal(t, tc.expQuery, query)
		})
	}
}
