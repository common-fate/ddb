package ddbmock

import (
	"context"
	"reflect"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb"
)

var _ ddb.Storage = &Client{}

// Client is a mock client which can be used to test ddb queries.
type Client struct {
	t          TestReporter
	mu         *sync.Mutex
	results    map[reflect.Type]mockResult
	getResults map[ddb.GetKey]mockGetResult
	// DeleteErr causes Delete() to return an error if it is set
	DeleteErr error
	// PutErr causes Put() to return an error if it is set
	PutErr error
	// PutBatchErr causes PutBatch() to return an error if it is set
	PutBatchErr error
	// DeleteBatchErr causes DeleteBatch() to return an error if it is set
	DeleteBatchErr error
	// TransactWriteItemsErr causes TransactWriteItems() to return an error if it is set
	TransactWriteItemsErr error
	// TransactionExecuteError causes Transaction objects created from this client
	// to fail with this error if set.
	TransactionExecuteErr error
}

// mockGetResult is the mocked result when Get() is called.
type mockGetResult struct {
	res   *ddb.GetItemResult
	value interface{}
	err   error
}

// mockResult is the mocked result when Query() is called.
type mockResult struct {
	res   *ddb.QueryResult
	value interface{}
	err   error
}

// New creates a new mock client which satisfies the ddb.Storage interface.
func New(t TestReporter) *Client {
	return &Client{
		t:          t,
		mu:         &sync.Mutex{},
		results:    make(map[reflect.Type]mockResult),
		getResults: make(map[ddb.GetKey]mockGetResult),
	}
}

// MockGet mocks a DynamoDB Get operation.
// The contents of the provided query will be used as the results.
//
// For example:
//
//	db := ddbmock.New()
//	db.MockGet(ddb.GetKey{PK: "1", SK: "1"}, Apple{Color: "red"})
//
//	var got Apple
//	db.Get(ctx, ddb.GetKey{PK: "1", SK: "1"}, &got)
//	// got now contains {Result: Apple{Color: "red"}} as defined by MockGet.
func (m *Client) MockGet(key ddb.GetKey, result interface{}) {

	// acquire a mutex lock in case the client is being used across multiple goroutines.
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getResults[key] = mockGetResult{
		value: result,
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
		res:   &ddb.QueryResult{},
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
		res:   &ddb.QueryResult{},
	}
}

// MockQueryWithErrWithResult mocks a DynamoDB query.
// It works the same as MockQueryWithErr, but allows a QueryResult to be set.
// The QueryResult argument can be nil, in which case a nil QueryResult is returned.
func (m *Client) MockQueryWithErrWithResult(qb ddb.QueryBuilder, res *ddb.QueryResult, err error) {
	t := reflect.TypeOf(qb)

	// acquire a mutex lock in case the client is being used across multiple goroutines.
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results[t] = mockResult{
		value: qb,
		err:   err,
		res:   res,
	}
}

// Query returns mock query results based on the type of the 'qb' argument.
func (m *Client) All(ctx context.Context, qb ddb.QueryBuilder, opts ...func(*ddb.QueryOpts)) error {
	_, err := m.Query(ctx, qb, opts...)
	if err != nil {
		return err
	}
	return nil
}

// Query returns mock query results based on the type of the 'qb' argument.
func (m *Client) Query(ctx context.Context, qb ddb.QueryBuilder, opts ...func(*ddb.QueryOpts)) (*ddb.QueryResult, error) {
	t := reflect.TypeOf(qb)
	got, ok := m.results[t]
	if !ok {
		m.t.Fatalf("no mock found for %s - call RegisterQuery(&%s{}) to set a mock response", t, reflect.TypeOf(qb).Elem().Name())
		return nil, nil
	}

	// If we got an error, return it and don't set the results of the query.
	if got.err != nil {
		return nil, got.err
	}

	// set the value of the QueryBuilder to our stored mock result.
	reflect.ValueOf(qb).Elem().Set(reflect.ValueOf(got.value).Elem())

	return got.res, nil
}

// Get returns mock query results based registered mock values.
func (m *Client) Get(ctx context.Context, key ddb.GetKey, item ddb.Keyer, opts ...func(*ddb.GetOpts)) (*ddb.GetItemResult, error) {
	got, ok := m.getResults[key]
	if !ok {
		m.t.Fatalf("no mock found for %+v - call MockGet() to set a mock response", key)
		return nil, nil
	}

	// If we got an error, return it and don't set the results of the query.
	if got.err != nil {
		return nil, got.err
	}

	// set the value of the item to our stored mock result.
	reflect.ValueOf(item).Elem().Set(reflect.ValueOf(got.value).Elem())

	return got.res, nil
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

func (m *Client) DeleteBatch(ctx context.Context, items ...ddb.Keyer) error {
	return m.DeleteBatchErr
}

func (m *Client) NewTransaction() ddb.Transaction {
	return &MockTransaction{ExecuteError: m.TransactionExecuteErr}
}

// Client returns nil. If you're writing tests which use
// DynamoDB implementation details you should probably be
// using integration tests!
func (m *Client) Client() *dynamodb.Client {
	return nil
}

func (m *Client) Table() string {
	return ""
}
