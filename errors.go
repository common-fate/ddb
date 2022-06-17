package ddb

import "errors"

// ErrNoItems is returned when we expect a query result to contain items,
// but it doesn't contain any.
var ErrNoItems error = errors.New("item query returned no items")

// ErrInvalidBatchSize is returned if an invalid batch size is specified when creating a ddb instance.
var ErrInvalidBatchSize error = errors.New("batch size must be greater than 0 and must not be greater than 25")
