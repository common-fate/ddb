package ddbtest

import (
	"context"
	"testing"

	"github.com/common-fate/ddb"
)

func TestTransactWriteItems(t *testing.T) {
	c := getTestClient(t)

	type testcase struct {
		name string
		tx   []ddb.TransactWriteItem
	}

	testcases := []testcase{
		{
			name: "ok",
			tx: []ddb.TransactWriteItem{
				{
					Put: Thing{
						Type: "test",
						ID:   "1",
					},
				},
				{
					Put: Thing{
						Type: "test",
						ID:   "2",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := c.TransactWriteItems(context.Background(), tc.tx)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}
