package ddbtest

import (
	"context"
	"reflect"
	"testing"

	"github.com/common-fate/ddb"
)

// PutFixtures inserts fixture data into the database.
// It's useful for provisioning data to be used in integration tests.
//
// The provided fixtures may be a single item or a slice of items.
//
// Each provided item must implement the ddb.Keyer interface.
// The test will fail if you try and insert an item which doesn't implement it.
//
// To insert multiple items:
//	ddbtest.PutFixtures(t, client, []MyData{{ID: "1"}, {ID: "2"}})
//
// To insert a single item:
//	ddbtest.PutFixtures(t, client, MyData{ID: "1"})
func PutFixtures(t *testing.T, c ddb.Storage, fixtures interface{}) {
	keyers := toKeyers(t, fixtures)
	err := c.PutBatch(context.Background(), keyers...)
	if err != nil {
		t.Fatal(err)
	}
}

// toKeyers tries to convert a provided object into Keyers, so that we can insert them in the database.
func toKeyers(t *testing.T, in interface{}) []ddb.Keyer {
	var keyers []ddb.Keyer
	if k, ok := in.(ddb.Keyer); ok {
		return []ddb.Keyer{k}
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(in)

		for i := 0; i < s.Len(); i++ {
			k, ok := s.Index(i).Interface().(ddb.Keyer)
			if !ok {
				// try converting as a pointer
				k, ok = s.Index(i).Addr().Interface().(ddb.Keyer)
			}
			if !ok {
				t.Fatalf("fixture %s must implement ddb.Keyer interface", reflect.TypeOf(s.Index(i)))
			}

			keyers = append(keyers, k)
		}
		return keyers
	}
	t.Fatalf("fixture %s must implement ddb.Keyer interface", reflect.TypeOf(in))
	return nil
}
