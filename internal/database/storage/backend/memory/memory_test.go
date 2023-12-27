package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHashTable(t *testing.T) {
	t.Parallel()

	table := NewHashTable()
	require.NotNil(t, table.data)
}

func TestHashTable_Set(t *testing.T) {
	t.Parallel()

	t.Run("test set not existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "set key"
		value := "value"

		err := table.Set(context.Background(), key, value)

		require.NoError(t, err)
		actValue, ok := table.data[key]
		require.True(t, ok)
		require.Equal(t, value, actValue)
	})

	t.Run("test set existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "set key"
		value := "value"
		table.data[key] = "old " + value

		err := table.Set(context.Background(), key, value)

		require.NoError(t, err)
		actValue, ok := table.data[key]
		require.True(t, ok)
		require.Equal(t, value, actValue)
	})
}

func TestHashTable_Get(t *testing.T) {
	t.Parallel()

	t.Run("test get existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "get key"
		value := "get value"
		table.data[key] = value

		actValue, ok, err := table.Get(context.Background(), key)

		require.NoError(t, err)
		require.True(t, ok)
		require.NoError(t, err)
		require.Equal(t, value, actValue)
	})

	t.Run("test get not existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "get key"

		actValue, ok, err := table.Get(context.Background(), key)

		require.NoError(t, err)
		require.False(t, ok)
		require.NoError(t, err)
		require.Equal(t, "", actValue)
	})
}

func TestHashTable_Delete(t *testing.T) {
	t.Parallel()

	t.Run("test delete existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "delete key"
		value := "delete value"
		table.data[key] = value

		err := table.Delete(context.Background(), key)

		require.NoError(t, err)
		_, ok := table.data[key]
		require.False(t, ok)
	})

	t.Run("test delete not existing key", func(t *testing.T) {
		table := NewHashTable()
		key := "delete key"

		err := table.Delete(context.Background(), key)

		require.NoError(t, err)
		_, ok := table.data[key]
		require.False(t, ok)
	})
}
