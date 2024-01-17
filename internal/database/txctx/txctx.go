package txctx

import "context"

type txnIDCtx struct{}

func CtxWithTx(ctx context.Context, txID int64) context.Context {
	return context.WithValue(ctx, txnIDCtx{}, txID)
}

func GetTxFromCtx(ctx context.Context) int64 {
	return ctx.Value(txnIDCtx{}).(int64)
}
