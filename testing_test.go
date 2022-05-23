package ddb

import "testing"

type exampleRowPtr struct{}

// DDBKeys implements the Keyer interface with a pointer receiver
func (r *exampleRowPtr) DDBKeys() (Keys, error) {
	return Keys{}, nil
}

type exampleRow struct{}

// DDBKeys implements the Keyer interface
func (r exampleRow) DDBKeys() (Keys, error) {
	return Keys{}, nil
}

func TestToKeyers(t *testing.T) {
	type testcase struct {
		name string
		give interface{}
	}

	testcases := []testcase{
		{"slice", []exampleRow{{}, {}}},
		{"slice of pointer receivers", []exampleRowPtr{{}, {}}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			toKeyers(t, tc.give)
		})
	}

}
