package ddbmock

import (
	"fmt"
	"reflect"

	"github.com/common-fate/ddb"
)

type Client struct {
	results map[string]interface{}
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

func (m *Client) Query(qb ddb.QueryBuilder) error {
	name := reflect.TypeOf(qb).String()
	got, ok := m.results[name]
	if !ok {
		return fmt.Errorf("no mock found for %s - call OnQuery(&%s{}) to set a mock response", name, reflect.TypeOf(qb).Elem().Name())
	}

	// set the value of the QueryBuilder to our stored mock result.
	reflect.ValueOf(qb).Elem().Set(reflect.ValueOf(got).Elem())

	return nil
}
