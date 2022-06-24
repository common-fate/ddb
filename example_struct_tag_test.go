package ddb_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb"
)

// Orange is an example object which we will
// show how to read from the database with some access patterns.
type Orange struct{}

type ListOrangesByColor struct {
	Color  string
	Result []Orange `ddb:"result"` // results will be unmarshaled into this field.
}

func (l *ListOrangesByColor) BuildQuery() (*dynamodb.QueryInput, error) {
	// the empty QueryInput is just for the example.
	// in a real query this wouldn't be empty.
	return &dynamodb.QueryInput{}, nil
}

// For queries that take parameters (like an object ID or a status), using the `ddb:"result"`
// struct tag is the simplest way to denote the results field of the query.
// ddb will unmarshal the results into the tagged field.

func Example_structTag() {
	ctx := context.TODO()

	lobc := ListOrangesByColor{Color: "light-orange"}
	c, _ := ddb.New(ctx, "example-table")
	_, _ = c.Query(ctx, &lobc, nil)

	// labc.Result is now populated with []Orange as fetched from DynamoDB
}
