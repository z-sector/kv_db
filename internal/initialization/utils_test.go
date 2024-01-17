package initialization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSize(t *testing.T) {
	t.Parallel()

	t.Run("with bytes", func(t *testing.T) {
		size, err := ParseSize("20B")
		require.NoError(t, err)
		require.Equal(t, 20, size)

		size, err = ParseSize("20b")
		require.NoError(t, err)
		require.Equal(t, 20, size)

		size, err = ParseSize("20")
		require.NoError(t, err)
		require.Equal(t, 20, size)
	})

	t.Run("with kilo bytes", func(t *testing.T) {
		size, err := ParseSize("20B")
		require.NoError(t, err)
		require.Equal(t, 20, size)

		size, err = ParseSize("20b")
		require.NoError(t, err)
		require.Equal(t, 20, size)

		size, err = ParseSize("20")
		require.NoError(t, err)
		require.Equal(t, 20, size)
	})

	t.Run("with mega bytes", func(t *testing.T) {
		size, err := ParseSize("20MB")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024, size)

		size, err = ParseSize("20Mb")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024, size)

		size, err = ParseSize("20mb")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024, size)
	})

	t.Run("with giga bytes", func(t *testing.T) {
		size, err := ParseSize("20GB")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024*1024, size)

		size, err = ParseSize("20Gb")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024*1024, size)

		size, err = ParseSize("20gb")
		require.NoError(t, err)
		require.Equal(t, 20*1024*1024*1024, size)
	})

	t.Run("incorrect size", func(t *testing.T) {
		_, err := ParseSize("-20")
		require.Error(t, err)

		_, err = ParseSize("20AB")
		require.Error(t, err)

		_, err = ParseSize("GB")
		require.Error(t, err)
	})
}
