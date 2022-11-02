package ddbmock

import (
	"context"

	"github.com/common-fate/ddb"
)

type MockTransaction struct {
	// ExecuteError causes Execute() to return with an error if set
	ExecuteError error
}

func (m *MockTransaction) Execute(ctx context.Context) error {
	return m.ExecuteError
}

func (m *MockTransaction) Put(item ddb.Keyer) {}

func (m *MockTransaction) Delete(item ddb.Keyer) {}
