package ddb

import (
	"context"
	"os"
	"reflect"
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

// PutFixtures inserts fixture data into the database.
// It's useful for provisioning data to be used in integration tests.
func PutFixtures(t *testing.T, c *Client, fixtures interface{}) {
	keyers := toKeyers(t, fixtures)
	err := c.PutBatch(context.Background(), keyers...)
	if err != nil {
		t.Fatal(err)
	}
}

// toKeyers tries to convert a provided object into Keyers, so that we can insert them in the database.
func toKeyers(t *testing.T, in interface{}) []Keyer {
	var keyers []Keyer
	if k, ok := in.(Keyer); ok {
		return []Keyer{k}
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(in)

		for i := 0; i < s.Len(); i++ {
			k, ok := s.Index(i).Interface().(Keyer)
			if !ok {
				t.Fatalf("fixture %s must implement Keyer interface", reflect.TypeOf(s.Index(i)))
			}
			keyers = append(keyers, k)
		}
		return keyers
	}
	t.Fatalf("fixture %s must implement Keyer interface", reflect.TypeOf(in))
	return nil
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
