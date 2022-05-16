package ddbmock

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

type testQuery struct {
	Result thing
}

func (g *testQuery) BuildQuery() (*dynamodb.QueryInput, error) {
	return &dynamodb.QueryInput{}, nil
}

type thing struct {
	ID string
}

func TestMock(t *testing.T) {
	type testcase struct {
		name    string
		want    thing
		mock    *testQuery
		wantErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			want: thing{ID: "hello"},
			mock: &testQuery{thing{ID: "hello"}},
		},
		{
			name:    "no mock provided",
			wantErr: errors.New("no mock found for *ddbmock.testQuery - call OnQuery(&testQuery{}) to set a mock response"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := New()
			if tc.mock != nil {
				m.Mock(tc.mock)
			}

			var q testQuery
			err := m.Query(context.Background(), &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			assert.Equal(t, tc.want, q.Result)
		})
	}

}
