package ddbmock

import (
	"context"
	"reflect"
	"sync"

	"github.com/common-fate/ddb"
)

var _ ddb.Storage = &Client{}

// Client is a mock client which can be used to test ddb queries.
type Client struct {
	t       TestReporter
	mu      *sync.Mutex
	results map[reflect.Type]mockResult
	// DeleteErr causes Delete() to return an error if it is set
	DeleteErr error
	// PutErr causes Put() to return an error if it is set
	PutErr error
	// PutBatchErr causes PutBatch() to return an error if it is set
	PutBatchErr error
	// TransactWriteItemsErr causes TransactWriteItems() to return an error if it is set
	TransactWriteItemsErr error
}

// mockResult is the mocked result when Query() is called.
type mockResult struct {
	value interface{}
	err   error
}

// New creates a new mock client which satisfies the ddb.Storage interface.
func New(t TestReporter) *Client {
	return &Client{
		t:       t,
		mu:      &sync.Mutex{},
		results: make(map[reflect.Type]mockResult),
	}
}

// MockQuery mocks a DynamoDB query.
// The contents of the provided query will be used as the results.
//
// For example:
//
//	db := ddbmock.New()
//	db.MockQuery(&getApple{Result: Apple{Color: "red"}})
//
//	var got getApple
//	db.Query(ctx, &got)
//	// got now contains {Result: Apple{Color: "red"}} as defined by MockQuery.
func (m *Client) MockQuery(qb ddb.QueryBuilder) {
	t := reflect.TypeOf(qb)

	// acquire a mutex lock in case the client is being used across multiple goroutines.
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results[t] = mockResult{
		value: qb,
	}
}

// MockQueryWithErr mocks a DynamoDB query.
// It works the same as MockQuery, but allows an error response to be set.
// The err argument can be nil, in which case a nil error is returned.
// The contents of the provided query will be used as the results.
//
// For example:
//
//	db := ddbmock.New()
//	db.MockQueryWithErr(&getApple{}, ddb.ErrNoItems)
//
//	var got getApple
//	err := db.Query(ctx, &got)
//	// err is equal to ddb.ErrNoItems.
func (m *Client) MockQueryWithErr(qb ddb.QueryBuilder, err error) {
	t := reflect.TypeOf(qb)

	// acquire a mutex lock in case the client is being used across multiple goroutines.
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results[t] = mockResult{
		value: qb,
		err:   err,
	}
}

// Query returns mock query results based on the type of the 'qb' argument.
func (m *Client) Query(ctx context.Context, qb ddb.QueryBuilder) error {
	t := reflect.TypeOf(qb)
	got, ok := m.results[t]
	if !ok {
		m.t.Fatalf("no mock found for %s - call RegisterQuery(&%s{}) to set a mock response", t, reflect.TypeOf(qb).Elem().Name())
		return nil
	}

	// If we got an error, return it and don't set the results of the query.
	if got.err != nil {
		return got.err
	}

	// set the value of the QueryBuilder to our stored mock result.
	reflect.ValueOf(qb).Elem().Set(reflect.ValueOf(got.value).Elem())

	return nil
}

func (m *Client) Put(ctx context.Context, item ddb.Keyer) error {
	return m.PutErr
}

func (m *Client) PutBatch(ctx context.Context, items ...ddb.Keyer) error {
	return m.PutBatchErr
}

func (m *Client) TransactWriteItems(ctx context.Context, tx []ddb.TransactWriteItem) error {
	return m.TransactWriteItemsErr
}

func (m *Client) Delete(ctx context.Context, item ddb.Keyer) error {
	return m.DeleteErr
}
