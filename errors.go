package ddb

import "errors"

// ErrNoItems is returned when we expect a query result to contain items,
// but it doesn't contain any.
var ErrNoItems error = errors.New("item query returned no items")
