package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database/compute"
	"kv_db/pkg/dlog"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	parser, err := NewParser(nil)
	require.Error(t, err)
	require.Nil(t, parser)

	parser, err = NewParser(dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, parser)
}

func TestMustParser(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustParser(nil)
	})

	require.NotPanics(t, func() {
		MustParser(dlog.NewNonSlog())
	})
}

func TestParse(t *testing.T) {
	ctx := context.Background()
	parser, err := NewParser(dlog.NewNonSlog())
	require.NoError(t, err)

	testcases := map[string]struct {
		query     string
		expTokens []string
		expErr    error
	}{
		"empty query": {
			query: "",
		},
		"query without tokens": {
			query: "   ",
		},
		"query with UTF symbols": {
			query:  "удалить",
			expErr: compute.ErrInvalidSymbol,
		},
		"query with one token": {
			query:     "set",
			expTokens: []string{"set"},
		},
		"query with two tokens": {
			query:     "set key",
			expTokens: []string{"set", "key"},
		},
		"query with two token with digits": {
			query:     "get 1key2",
			expTokens: []string{"get", "1key2"},
		},
		"query with one token with underscores": {
			query:     "_set__",
			expTokens: []string{"_set__"},
		},
		"query with one token with invalid symbols": {
			query:  ".set#",
			expErr: compute.ErrInvalidSymbol,
		},
		"query with two tokens with additional spaces": {
			query:     " get   key  ",
			expTokens: []string{"get", "key"},
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tokens, err := parser.ParseQuery(ctx, tc.query)
			require.ErrorIs(t, err, tc.expErr)
			require.Equal(t, tc.expTokens, tokens)
		})
	}
}
