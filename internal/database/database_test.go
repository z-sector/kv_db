package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"kv_db/internal/database/comd"
	"kv_db/pkg/dlog"
)

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	compute, storage := getMockComputeAndStorage(t)

	database, err := NewDatabase(compute, storage, nil)
	require.Error(t, err)
	require.Nil(t, database)

	database, err = NewDatabase(compute, nil, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, database)

	database, err = NewDatabase(nil, storage, dlog.NewNonSlog())
	require.Error(t, err)
	require.Nil(t, database)

	database, err = NewDatabase(compute, storage, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, database)
}

func TestMustDatabase(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		MustDatabase(nil, nil, nil)
	})

	require.NotPanics(t, func() {
		compute, storage := getMockComputeAndStorage(t)

		MustDatabase(compute, storage, dlog.NewNonSlog())
	})
}

func TestDatabase_SetCommand(t *testing.T) {
	t.Parallel()

	inputQuery := "SET keyS valueS"
	arg0 := "keyS"
	arg1 := "valueS"
	query := Query{
		commandID: comd.SetCommandID,
		arguments: []string{arg0, arg1},
	}

	t.Run("success set command", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)

		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Set(gomock.Any(), arg0, arg1).Return(nil)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, "[ok]", res)
	})

	t.Run("error set command", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		expError := errors.New("test error")
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Set(gomock.Any(), arg0, arg1).Return(expError)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, fmt.Sprintf("[error] %s", expError), res)
	})
}

func TestDatabase_GetCommand(t *testing.T) {
	t.Parallel()

	inputQuery := "GET keyG"
	arg0 := "keyG"
	expValue := "valueG"
	query := Query{
		commandID: comd.GetCommandID,
		arguments: []string{arg0},
	}

	t.Run("success get command with value", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Get(gomock.Any(), arg0).Return(expValue, true, nil)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, fmt.Sprintf("[ok] %s", expValue), res)
	})

	t.Run("success get command with nil", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Get(gomock.Any(), arg0).Return("", false, nil)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, "[nil]", res)
	})

	t.Run("error get command", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Get(gomock.Any(), arg0).Return("", false, errors.New("test error"))

		res := database.HandleQuery(ctx, inputQuery)

		require.True(t, strings.HasPrefix(res, "[error]"))
	})
}

func TestDatabase_DelCommand(t *testing.T) {
	t.Parallel()

	inputQuery := "DEL keyD"
	arg0 := "keyD"
	query := Query{
		commandID: comd.DelCommandID,
		arguments: []string{arg0},
	}

	t.Run("success del command", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Delete(gomock.Any(), arg0).Return(nil)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, "[ok]", res)
	})

	t.Run("error del command", func(t *testing.T) {
		ctx := context.Background()
		compute, storage := getMockComputeAndStorage(t)
		database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
		require.NoError(t, err)
		expError := errors.New("test error")
		compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)
		storage.EXPECT().Delete(gomock.Any(), arg0).Return(expError)

		res := database.HandleQuery(ctx, inputQuery)

		require.Equal(t, fmt.Sprintf("[error] %s", expError), res)
	})
}

func TestDatabase_UnknownCommand(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	compute, storage := getMockComputeAndStorage(t)
	database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
	require.NoError(t, err)
	inputQuery := "Unknown"
	query := Query{
		commandID: comd.UnknownCommandID,
		arguments: nil,
	}
	compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(query, nil)

	res := database.HandleQuery(ctx, inputQuery)

	require.Equal(t, "[error] internal configuration error", res)
}

func TestDatabase_ComputeError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	compute, storage := getMockComputeAndStorage(t)
	database, err := NewDatabase(compute, storage, dlog.NewNonSlog())
	require.NoError(t, err)
	inputQuery := "TEST"
	expError := errors.New("test error")
	compute.EXPECT().HandleQuery(gomock.Any(), gomock.Eq(inputQuery)).Return(Query{}, expError)

	res := database.HandleQuery(ctx, inputQuery)

	require.Equal(t, fmt.Sprintf("[error] %s", expError), res)
}

func getMockComputeAndStorage(t *testing.T) (*MockComputeLayer, *MockStorageLayer) {
	t.Helper()

	ctrl := gomock.NewController(t)
	return NewMockComputeLayer(ctrl), NewMockStorageLayer(ctrl)
}
