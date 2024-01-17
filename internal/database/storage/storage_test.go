package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kv_db/internal/database/comd"
	"kv_db/internal/database/storage/wal"
	"kv_db/pkg/dfuture"
	"kv_db/pkg/dlog"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	engine, _ := getMockEngineAndWal(t)

	storage, err := NewStorage(nil, nil, nil)
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil, nil)
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(nil, nil, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, storage)

	storage, err = NewStorage(engine, nil, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestMustStorage(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustStorage(nil, nil, nil)
	})

	require.NotPanics(t, func() {
		engine, _ := getMockEngineAndWal(t)

		MustStorage(engine, nil, dlog.NewNonSlog())
	})
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	key := "set key"
	value := "set value"
	expErr := errors.New("test error")

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		require.NoError(t, err)
		futureRes := make(chan error, 1)
		futureRes <- nil
		walMock.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(dfuture.NewFuture(futureRes))
		engine.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(nil)

		err = storage.Set(ctx, key, value)

		require.NoError(t, err)
	})

	t.Run("wal error", func(t *testing.T) {
		ctx := context.Background()
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		require.NoError(t, err)
		futureRes := make(chan error, 1)
		futureRes <- expErr
		walMock.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(dfuture.NewFuture(futureRes))

		err = storage.Set(ctx, key, value)

		require.ErrorIs(t, err, expErr)
	})

	t.Run("engine error", func(t *testing.T) {
		ctx := context.Background()
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		require.NoError(t, err)
		futureRes := make(chan error, 1)
		futureRes <- nil
		walMock.EXPECT().
			Set(gomock.Eq(ctx), gomock.Eq(key), gomock.Eq(value)).
			Return(dfuture.NewFuture(futureRes))
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
		engine, _ := getMockEngineAndWal(t)
		storage, err := NewStorage(engine, nil, dlog.NewNonSlog())
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
		engine, _ := getMockEngineAndWal(t)
		storage, err := NewStorage(engine, nil, dlog.NewNonSlog())
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
		engine, _ := getMockEngineAndWal(t)
		storage, err := NewStorage(engine, nil, dlog.NewNonSlog())
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
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		futureRes := make(chan error, 1)
		futureRes <- nil
		walMock.EXPECT().
			Del(gomock.Eq(ctx), gomock.Eq(key)).
			Return(dfuture.NewFuture(futureRes))
		require.NoError(t, err)
		engine.EXPECT().
			Delete(gomock.Eq(ctx), gomock.Eq(key)).
			Return(nil)

		err = storage.Delete(ctx, key)

		require.NoError(t, err)
	})

	t.Run("engine error", func(t *testing.T) {
		ctx := context.Background()
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		require.NoError(t, err)
		futureRes := make(chan error, 1)
		futureRes <- expErr
		walMock.EXPECT().
			Del(gomock.Eq(ctx), gomock.Eq(key)).
			Return(dfuture.NewFuture(futureRes))

		err = storage.Delete(ctx, key)

		require.ErrorIs(t, err, expErr)
	})

	t.Run("engine error", func(t *testing.T) {
		ctx := context.Background()
		engine, walMock := getMockEngineAndWal(t)
		walMock.EXPECT().Recover().Return(nil, nil)
		walMock.EXPECT().Start().Return()
		storage, err := NewStorage(engine, walMock, dlog.NewNonSlog())
		require.NoError(t, err)
		futureRes := make(chan error, 1)
		futureRes <- nil
		walMock.EXPECT().
			Del(gomock.Eq(ctx), gomock.Eq(key)).
			Return(dfuture.NewFuture(futureRes))
		engine.EXPECT().
			Delete(gomock.Eq(ctx), gomock.Eq(key)).
			Return(expErr)

		err = storage.Delete(ctx, key)

		require.ErrorIs(t, err, expErr)
	})
}

func TestNewStorage_WalRecover(t *testing.T) {
	t.Parallel()

	engineMock, walMock := getMockEngineAndWal(t)
	logs := []wal.LogData{
		{CommandID: comd.SetCommandID, Arguments: []string{"key", "value"}},
		{CommandID: comd.DelCommandID, Arguments: []string{"key"}},
	}
	walMock.EXPECT().Recover().Return(logs, nil)
	engineMock.EXPECT().Set(gomock.Any(), gomock.Eq("key"), gomock.Eq("value")).Times(1)
	engineMock.EXPECT().Delete(gomock.Any(), gomock.Eq("key")).Times(1)
	walMock.EXPECT().Start().Return()

	storage, err := NewStorage(engineMock, walMock, dlog.NewNonSlog())

	require.NoError(t, err)
	require.NotNil(t, storage)
}

func getMockEngineAndWal(t *testing.T) (*MockEngine, *MockWAL) {
	t.Helper()

	ctrl := gomock.NewController(t)
	return NewMockEngine(ctrl), NewMockWAL(ctrl)
}
