package ddb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeysJSON(t *testing.T) {
	type testcase struct {
		name string
		give Keys
		want string
	}

	testcases := []testcase{
		{"no gsi", Keys{PK: "primary", SK: "sort"}, `{"PK":"primary","SK":"sort"}`},
		{"with gsi", Keys{PK: "primary", SK: "sort", GSI1PK: "1", GSI2PK: "1"}, `{"PK":"primary","SK":"sort","GSI1PK":"1","GSI2PK":"1"}`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.give)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, string(got))
		})
	}

}
