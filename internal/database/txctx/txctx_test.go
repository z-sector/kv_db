package txctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtxWithTx(t *testing.T) {
	txID := int64(100)
	ctx := context.Background()

	newCtx := CtxWithTx(ctx, txID)

	act, ok := newCtx.Value(txnIDCtx{}).(int64)
	require.True(t, ok)
	require.Equal(t, txID, act)
}

func TestGetTxFromCtx(t *testing.T) {
	txID := int64(100)
	ctx := context.WithValue(context.Background(), txnIDCtx{}, txID)

	act := GetTxFromCtx(ctx)

	require.Equal(t, txID, act)
}
