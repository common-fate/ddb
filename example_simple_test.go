package ddb_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb"
)

// Apple is an example object which we will
// show how to read from the database with some access patterns.
type Apple struct{}

type ListApples []Apple

func (l *ListApples) BuildQuery() (*dynamodb.QueryInput, error) {
	// the empty QueryInput is just for the example.
	// in a real query this wouldn't be empty.
	return &dynamodb.QueryInput{}, nil
}

// For simple queries you can declare the query as a type alias.
// ddb will unmarshal the results directly into the query struct, as shown below.
func Example_simple() {
	ctx := context.TODO()

	var la ListApples
	c, _ := ddb.New(ctx, "example-table")
	_ = c.Query(ctx, &la)

	// la is now populated with []Apple as fetched from DynamoDB
}

type ListApplesByColor struct {
	Color  string
	Result []Apple `ddb:"result"`
}

func (l *ListApplesByColor) BuildQuery() (*dynamodb.QueryInput, error) {
	// the empty QueryInput is just for the example.
	// in a real query this wouldn't be empty.
	return &dynamodb.QueryInput{}, nil
}
