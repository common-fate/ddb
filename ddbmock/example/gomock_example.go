package example

import (
	"context"

	"github.com/common-fate/ddb"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock_gomock_example.go -package=example . GoMockStorage

// GoMockStorage is an example interface to be mocked with GoMock.
type GoMockStorage interface {
	Query(ctx context.Context, qb ddb.QueryBuilder) error
}
