package example

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb/ddbmock"
	gomock "github.com/golang/mock/gomock"
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

func TestGoMockQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := testQuery{Result: thing{ID: "test"}}
	m := NewMockGoMockStorage(ctrl)
	ddbmock.GoMockQuery(m, &want)

	var got testQuery
	err := m.Query(context.Background(), &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestGoMockQueryErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := errors.New("test")
	m := NewMockGoMockStorage(ctrl)
	ddbmock.GoMockQueryErr(m, &testQuery{}, want)

	got := m.Query(context.Background(), &testQuery{})
	assert.Equal(t, want, got)
}

func TestGoMockQueryInvalid(t *testing.T) {
	assert.PanicsWithError(t, "could not find an EXPECT() method on goMockStorage. Ensure that you are using GoMockQuery with a mock generated from GoMock.", func() {
		ddbmock.GoMockQuery("other", &testQuery{})
	})
}
