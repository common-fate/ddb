package ddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testitem struct{}

func (t testitem) DDBKeys() (Keys, error) {
	k := Keys{
		PK: "PK",
		SK: "SK",
	}
	return k, nil
}

func TestTransaction_buildTransactWriteItemsPayload(t *testing.T) {

	tests := []struct {
		name        string
		putItems    []Keyer
		deleteItems []Keyer
		want        []TransactWriteItem
	}{
		{
			name: "ok",
			putItems: []Keyer{
				testitem{},
			},
			want: []TransactWriteItem{
				{
					Put: testitem{},
				},
			},
		},
		{
			name: "put and delete",
			putItems: []Keyer{
				testitem{},
				testitem{},
			},
			deleteItems: []Keyer{
				testitem{},
			},
			want: []TransactWriteItem{
				{
					Put: testitem{},
				},
				{
					Put: testitem{},
				},
				{
					Delete: testitem{},
				},
			},
		},
		{
			name:     "just delete",
			putItems: []Keyer{},
			deleteItems: []Keyer{
				testitem{},
				testitem{},
			},
			want: []TransactWriteItem{
				{
					Delete: testitem{},
				},
				{
					Delete: testitem{},
				},
			},
		},
		{
			name: "nothing",
			want: []TransactWriteItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &DBTransaction{
				putItems:    tt.putItems,
				deleteItems: tt.deleteItems,
			}
			got := tr.buildTransactWriteItemsPayload()
			assert.Equal(t, tt.want, got)
		})
	}
}
