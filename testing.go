package ddb

import (
	"context"
	"os"
	"testing"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

// QueryTestCase is a test case for running integration tests which call Query().
type QueryTestCase struct {
	Name    string
	Query   QueryBuilder
	Want    interface{}
	WantErr error
}

// RunQueryTests runs standardised integration tests to check the behaviour of a QueryBuilder.
func RunQueryTests(t *testing.T, c *Client, testcases []QueryTestCase) {
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err := c.Query(context.Background(), tc.Query)
			if err != nil && tc.WantErr == nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.WantErr, err)
			changelog, err := diff.Diff(tc.Want, tc.Query)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

// getTestClient returns a test ddb.Client.
// if TESTING_DYNAMODB_TABLE is not set it skips the test.
func getTestClient(t *testing.T) *Client {
	if os.Getenv("TESTING_DYNAMODB_TABLE") == "" {
		t.Skip("TESTING_DYNAMODB_TABLE is not set")
	}

	c, err := New(context.Background(), os.Getenv("TESTING_DYNAMODB_TABLE"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}
