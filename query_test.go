package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindResultsTag(t *testing.T) {
	type qbString struct {
		Result string `ddb:"result"`
	}

	type qbInt struct {
		Result int `ddb:"result"`
	}

	type testcase struct {
		name     string
		give     interface{}
		wantType string
	}

	testcases := []testcase{
		{
			name:     "string",
			give:     &qbString{},
			wantType: "string",
		},
		{
			name:     "int",
			give:     &qbInt{},
			wantType: "int",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := findResultsTag(tc.give)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantType, got.Elem().Type().Name())
		})
	}

}
