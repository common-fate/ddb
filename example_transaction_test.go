package ddb_test

import (
	"context"

	"github.com/common-fate/ddb"
)

type Item struct {
	ID string
}

func (i Item) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: i.ID,
		SK: "SK",
	}
	return k, nil
}

// ddb exposes a transactions API.
func Example_transaction() {
	ctx := context.TODO()

	c, _ := ddb.New(ctx, "example-table")
	tx := c.NewTransaction()

	tx.Put(Item{ID: "1"})
	tx.Put(Item{ID: "2"})
	tx.Delete(Item{ID: "3"})
	_ = tx.Execute(ctx)
}
