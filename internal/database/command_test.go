package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCommandIDByName(t *testing.T) {
	t.Parallel()

	require.Equal(t, UnknownCommandID, GetCommandIDByName("TEST"))
	require.Equal(t, SetCommandID, GetCommandIDByName("SET"))
	require.Equal(t, GetCommandID, GetCommandIDByName("GET"))
	require.Equal(t, DelCommandID, GetCommandIDByName("DEL"))
}
