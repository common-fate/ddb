package ddbtest

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
)

type exampleRowPtr struct{}

// DDBKeys implements the Keyer interface with a pointer receiver
func (r *exampleRowPtr) DDBKeys() (ddb.Keys, error) {
	return ddb.Keys{}, nil
}

type exampleRow struct{}

// DDBKeys implements the Keyer interface
func (r exampleRow) DDBKeys() (ddb.Keys, error) {
	return ddb.Keys{}, nil
}

func TestPutFixture(t *testing.T) {
	type testcase struct {
		name string
		give interface{}
	}

	c := ddbmock.New(t)

	testcases := []testcase{
		{"slice", []exampleRow{{}, {}}},
		{"slice of pointer receivers", []exampleRowPtr{{}, {}}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			PutFixtures(t, c, tc.give)
		})
	}

}
