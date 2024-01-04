package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kv_db/pkg/dlog"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	engine := getMockEngine(t)

	storage, err := NewStorage(nil, nil)
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil)
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(nil, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(engine, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestMustStorage(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustStorage(nil, nil)
	})

	require.NotPanics(t, func() {
		engine := getMockEngine(t)

		MustStorage(engine, dlog.NewNonSlog())
	})
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	key := "set key"
	value := "set value"
	expErr := errors.New("test error")

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(nil)

		err = storage.Set(ctx, key, value)

		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(expErr)

		err = storage.Set(ctx, key, value)

		require.ErrorIs(t, err, expErr)
	})
}

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	key := "get key"
	value := "get value"
	expErr := errors.New("test error")

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Get(gomock.Eq(ctx), gomock.Eq(key)).
			Return(value, true, nil)

		result, found, err := storage.Get(ctx, key)

		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, value, result)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Get(gomock.Eq(ctx), gomock.Eq(key)).
			Return("", false, nil)

		result, found, err := storage.Get(ctx, key)

		require.NoError(t, err)
		require.False(t, found)
		require.Empty(t, result)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Get(gomock.Eq(ctx), gomock.Eq(key)).
			Return("", false, expErr)

		result, found, err := storage.Get(ctx, key)

		require.ErrorIs(t, err, expErr)
		require.False(t, found)
		require.Empty(t, result)
	})
}

func TestStorage_Delete(t *testing.T) {
	t.Parallel()

	key := "delete key"
	expErr := errors.New("test error")

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Delete(gomock.Eq(ctx), gomock.Eq(key)).
			Return(nil)

		err = storage.Delete(ctx, key)

		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		engine := getMockEngine(t)
		storage, err := NewStorage(engine, dlog.NewNonSlog())
		require.NoError(t, err)
		engine.EXPECT().
			Delete(gomock.Eq(ctx), gomock.Eq(key)).
			Return(expErr)

		err = storage.Delete(ctx, key)

		require.ErrorIs(t, err, expErr)
	})
}

func getMockEngine(t *testing.T) *MockEngine {
	t.Helper()

	ctrl := gomock.NewController(t)
	return NewMockEngine(ctrl)
}
