package ddbmock

import (
	"context"
	"fmt"
	"reflect"

	"github.com/common-fate/ddb"
)

var _ ddb.Storage = &Client{}

type Client struct {
	results map[string]interface{}
	// QueryErr causes Query() to return an error if it is set
	QueryErr error
	// PutErr causes Put() to return an error if it is set
	PutErr error
	// PutBatchErr causes PutBatch() to return an error if it is set
	PutBatchErr error
}

func New() *Client {
	return &Client{
		results: make(map[string]interface{}),
	}
}

func (m *Client) OnQuery(ap ddb.QueryBuilder) {
	name := reflect.TypeOf(ap).String()
	m.results[name] = ap
}

func (m *Client) Query(ctx context.Context, qb ddb.QueryBuilder) error {
	if m.QueryErr != nil {
		return m.QueryErr
	}

	name := reflect.TypeOf(qb).String()
	got, ok := m.results[name]
	if !ok {
		return fmt.Errorf("no mock found for %s - call OnQuery(&%s{}) to set a mock response", name, reflect.TypeOf(qb).Elem().Name())
	}

	// set the value of the QueryBuilder to our stored mock result.
	reflect.ValueOf(qb).Elem().Set(reflect.ValueOf(got).Elem())

	return nil
}

func (m *Client) Put(ctx context.Context, item ddb.Keyer) error {
	return m.PutErr
}

func (m *Client) PutBatch(ctx context.Context, items ...ddb.Keyer) error {
	return m.PutBatchErr
}
