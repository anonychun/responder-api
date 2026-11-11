package current

import (
	"context"

	"github.com/uptrace/bun"
)

type key int

const txKey key = iota

func Tx(ctx context.Context) *bun.Tx {
	tx, _ := ctx.Value(txKey).(*bun.Tx)
	return tx
}

func SetTx(ctx context.Context, tx *bun.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
