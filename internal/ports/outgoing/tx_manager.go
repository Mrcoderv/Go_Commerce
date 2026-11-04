package outgoing

import "context"

// TxManager runs fn inside one database transaction.
// error, everything is rolled back.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
