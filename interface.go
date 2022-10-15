package ddb

import "context"

// Storage defines a common interface to make testing ddb easier.
// Both the real and mock clients meet this interface.
type Storage interface {
	Query(ctx context.Context, qb QueryBuilder, opts ...func(*QueryOpts)) (*QueryResult, error)
	Put(ctx context.Context, item Keyer) error
	PutBatch(ctx context.Context, items ...Keyer) error
	TransactWriteItems(ctx context.Context, tx []TransactWriteItem) error
	NewTransaction() Transaction
	Delete(ctx context.Context, item Keyer) error
	DeleteBatch(ctx context.Context, items ...Keyer) error
	// Get performs a GetItem call to fetch a single item from DynamoDB.
	// The results are written to the 'item' argument. This argument
	// must be passed by reference to the method.
	//
	// 	var item MyItem
	//	db.Get(ctx, ddb.GetKey{PK: ..., SK: ...}, &item)
	Get(ctx context.Context, key GetKey, item Keyer) (*GetItemResult, error)
}

// Transactions allow atomic write operations to be made to a DynamoDB table.
// DynamoDB transactions support up to 100 operations.
//
// Calling Put() and Delete() on a transaction register items in memory to be
// written to the table. No API calls are performed until Execute() is called.
type Transaction interface {
	// Put adds an item to be written in the transaction.
	Put(item Keyer)
	// Delete adds a item to be delete in the transaction.
	Delete(item Keyer)
	// Execute the transaction.
	// This calls the TransactWriteItems API.
	// See: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
	Execute(ctx context.Context) error
}
