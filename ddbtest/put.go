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
//
// Ideally, we would like to define PutFixtures as
//	PutFixtures(t *testing.T, c ddb.Storage, fixtures ...ddb.Keyer)
// However, Go does not automatically cast slices of structs to a slice of an interface, even if the struct
// meets the interface.
// So if we were to call PutFixtures with
//	PutFixtures(t, c, []MyObject{{ID: "1"}, {ID: "2"}})
// This is invalid.
//
// To save lots of repetive conversion code when writing tests, we allow users to provide *any* object as
// the 'in' argument and toKeyers will try to convert it to a slice of Keyers so that they can be inserted.
//
// To do this, toKeyers uses reflection to inspect the type of the provided interface.
func toKeyers(t *testing.T, in interface{}) []ddb.Keyer {
	var keyers []ddb.Keyer

	// check if we can cast the object to a Keyer.
	// If we can, 'in' was a single item, so return immediately with a slice containing the item.
	if k, ok := in.(ddb.Keyer); ok {
		return []ddb.Keyer{k}
	}

	// if the provided value isn't a slice, we can't iterate through it, so return an error.
	if reflect.TypeOf(in).Kind() != reflect.Slice {
		t.Fatalf("fixture %s must implement ddb.Keyer interface", reflect.TypeOf(in))
		return nil
	}

	s := reflect.ValueOf(in)

	// loop through the items in the slice to convert them all to Keyers.
	for i := 0; i < s.Len(); i++ {
		// check if the item can be converted to a Keyer
		k, ok := s.Index(i).Interface().(ddb.Keyer)

		if !ok {
			// try taking the address of the item and converting it to a Keyer
			k, ok = s.Index(i).Addr().Interface().(ddb.Keyer)
		}

		if !ok {
			// we can't convert it, so return an error.
			t.Fatalf("fixture %s must implement ddb.Keyer interface", reflect.TypeOf(s.Index(i)))
		}

		// if we get here, we successfully converted the item, so add it to our results slice.
		keyers = append(keyers, k)
	}
	return keyers
}
