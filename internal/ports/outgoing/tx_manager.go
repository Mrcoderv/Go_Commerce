package outgoing

import "context"

// TxManager runs fn inside one database transaction.
// If fn returns nil, the transaction is committed.
// If fn returns an error or panics, it is rolled back.
// Repositories called with the ctx that fn receives take part in that same transaction.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}