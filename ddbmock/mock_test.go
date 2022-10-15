package ddbmock

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb"
	"github.com/stretchr/testify/assert"
)

type testQuery struct {
	Result thing
}

func (g *testQuery) BuildQuery() (*dynamodb.QueryInput, error) {
	return &dynamodb.QueryInput{}, nil
}

type secondQuery struct {
	Result thing
}

func (g *secondQuery) BuildQuery() (*dynamodb.QueryInput, error) {
	return &dynamodb.QueryInput{}, nil
}

type thing struct {
	ID string
}

func (t thing) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: "PK",
		SK: "SK",
	}

	return k, nil
}

func TestMockQuery(t *testing.T) {
	type testcase struct {
		name string
		want thing
		mock *testQuery
	}

	testcases := []testcase{
		{
			name: "ok",
			want: thing{ID: "hello"},
			mock: &testQuery{thing{ID: "hello"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := New(&mockTestReporter{})
			if tc.mock != nil {
				m.MockQuery(tc.mock)
			}

			var q testQuery
			_, err := m.Query(context.Background(), &q)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, q.Result)
		})
	}
}

func TestMockQueryFailure(t *testing.T) {
	type testcase struct {
		name string
		want []string // logs from the test failure
		mock ddb.QueryBuilder
	}

	testcases := []testcase{
		{
			name: "wrong query type",
			want: []string{"no mock found for *ddbmock.testQuery - call RegisterQuery(&testQuery{}) to set a mock response"},
			mock: &secondQuery{thing{ID: "hello"}},
		},
		{
			name: "no mock set",
			want: []string{"no mock found for *ddbmock.testQuery - call RegisterQuery(&testQuery{}) to set a mock response"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tt := &mockTestReporter{}

			m := New(tt)
			if tc.mock != nil {
				m.MockQuery(tc.mock)
			}

			var q testQuery
			_, _ = m.Query(context.Background(), &q)
			assert.Equal(t, tc.want, tt.Logs)
		})
	}
}

func TestMockQueryWithErr(t *testing.T) {
	type testcase struct {
		name    string
		want    thing
		mock    *testQuery
		mockErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			want: thing{ID: "hello"},
			mock: &testQuery{thing{ID: "hello"}},
		},
		{
			name:    "error",
			want:    thing{}, // we shouldn't get a result back if we got an error
			mock:    &testQuery{thing{ID: "hello"}},
			mockErr: ddb.ErrNoItems,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := New(&mockTestReporter{})
			if tc.mock != nil {
				m.MockQueryWithErr(tc.mock, tc.mockErr)
			}

			var q testQuery
			_, err := m.Query(context.Background(), &q)
			assert.Equal(t, tc.mockErr, err)
			assert.Equal(t, tc.want, q.Result)
		})
	}
}

func TestMockGet(t *testing.T) {
	type testcase struct {
		name      string
		want      thing
		mockThing *thing
		mockErr   error
	}

	testcases := []testcase{
		{
			name:      "ok",
			want:      thing{ID: "hello"},
			mockThing: &thing{ID: "hello"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := New(&mockTestReporter{})
			if tc.mockThing != nil {
				m.MockGet(ddb.GetKey{PK: "1", SK: "1"}, tc.mockThing)
			}

			var got thing
			_, err := m.Get(context.Background(), ddb.GetKey{PK: "1", SK: "1"}, &got)
			assert.Equal(t, tc.mockErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
