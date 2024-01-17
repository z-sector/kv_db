package database

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database/comd"
)

func TestNewQuery(t *testing.T) {
	args := []string{"arg1", "arg2"}
	query := NewQuery(comd.GetCommandID, args)
	require.Equal(t, comd.GetCommandID, query.CommandID())
	require.Equal(t, args, query.Arguments())
}
