package ddbtest

import (
	"context"
	"os"
	"testing"

	"github.com/common-fate/ddb"
	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

// QueryTestCase is a test case for running integration tests which call Query().
type QueryTestCase struct {
	Name      string
	Query     ddb.QueryBuilder
	QueryOpts []func(*ddb.QueryOpts)
	Want      ddb.QueryBuilder
	WantErr   error
}

// RunQueryTests runs standardised integration tests to check the behaviour of a QueryBuilder.
func RunQueryTests(t *testing.T, c *ddb.Client, testcases []QueryTestCase) {
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {

			_, err := c.Query(context.Background(), tc.Query, tc.QueryOpts...)
			if err != nil && tc.WantErr == nil {
				t.Fatal(err)
			}

			if tc.WantErr != nil {
				// just compare the errors, as we don't care
				//about what the result would be if an error is returned.
				assert.Equal(t, tc.WantErr, err)
			} else {

				// we don't expect an error here, so compare the results to what we expected.
				changelog, err := diff.Diff(tc.Want, tc.Query)
				assert.NoError(t, err)
				if len(changelog) != 0 {
					// Go doesn't consistently order slices, so just calling assert.Equal
					// causes test cases to fail when the results are out of order
					// compared to what we want.
					// using the changelog length here is a bit of a hack to prevent this,
					// as the diff library ignores the order of slices.
					//
					// If we get here, calling assert.Equal() will definitely fail.
					// This gives us a developer-friendly error message we can use
					// to fix our tests faster.
					assert.Equal(t, tc.Want, tc.Query)
				}
			}
		})
	}
}

// getTestClient returns a test ddb.Client.
// if TESTING_DYNAMODB_TABLE is not set it skips the test.
func getTestClient(t *testing.T, opts ...func(*ddb.Client)) *ddb.Client {
	if os.Getenv("TESTING_DYNAMODB_TABLE") == "" {
		t.Skip("TESTING_DYNAMODB_TABLE is not set")
	}

	c, err := ddb.New(context.Background(), os.Getenv("TESTING_DYNAMODB_TABLE"), opts...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
