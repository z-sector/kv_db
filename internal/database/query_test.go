package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQuery(t *testing.T) {
	args := []string{"arg1", "arg2"}
	query := NewQuery(GetCommandID, args)
	require.Equal(t, GetCommandID, query.CommandID())
	require.Equal(t, args, query.Arguments())
}
