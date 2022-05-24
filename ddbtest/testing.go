package ddbtest

import (
	"context"
	"os"
	"testing"

	"github.com/common-fate/ddb"
	"github.com/stretchr/testify/assert"
)

// QueryTestCase is a test case for running integration tests which call Query().
type QueryTestCase struct {
	Name    string
	Query   ddb.QueryBuilder
	Want    ddb.QueryBuilder
	WantErr error
}

// RunQueryTests runs standardised integration tests to check the behaviour of a QueryBuilder.
func RunQueryTests(t *testing.T, c *ddb.Client, testcases []QueryTestCase) {
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err := c.Query(context.Background(), tc.Query)
			if err != nil && tc.WantErr == nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.WantErr, err)
			assert.Equal(t, tc.Want, tc.Query)
		})
	}
}

// getTestClient returns a test ddb.Client.
// if TESTING_DYNAMODB_TABLE is not set it skips the test.
func getTestClient(t *testing.T) *ddb.Client {
	if os.Getenv("TESTING_DYNAMODB_TABLE") == "" {
		t.Skip("TESTING_DYNAMODB_TABLE is not set")
	}

	c, err := ddb.New(context.Background(), os.Getenv("TESTING_DYNAMODB_TABLE"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}
