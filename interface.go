package ddb

import "context"

// Storage defines a common interface to make testing ddb easier.
// Both the real and mock clients meet this interface.
type Storage interface {
	Query(ctx context.Context, qb QueryBuilder) error
	Put(ctx context.Context, item Keyer) error
	PutBatch(ctx context.Context, items ...Keyer) error
}
