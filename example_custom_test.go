package ddb_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
)

type Car struct{}

type Wheel struct{}

// Result is an example complex query result for the example access pattern.
// In this example, we fetch both a Car object as well as it's associated Wheels
// in the same query
type Result struct {
	Car    Car
	Wheels []Wheel
}

type ListCarAndWheelsByColor struct {
	Color  string
	Result Result
}

func (l *ListCarAndWheelsByColor) BuildQuery() (*dynamodb.QueryInput, error) {
	// the empty QueryInput is just for the example.
	// in a real query this wouldn't be empty.
	return &dynamodb.QueryInput{}, nil
}

func (l *ListCarAndWheelsByColor) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	// an example of custom unmarshalling logic for complex queries which return multiple item types
	for _, item := range out.Items {
		typeField, ok := item["type"].(*types.AttributeValueMemberS)
		if !ok {
			return fmt.Errorf("couldn't unmarshal: %+v", item)
		}

		if typeField.Value == "car" {
			err := attributevalue.UnmarshalMap(item, &l.Result.Car)
			if err != nil {
				return err
			}
		} else {
			var wheel Wheel
			err := attributevalue.UnmarshalMap(item, &wheel)
			if err != nil {
				return err
			}
			l.Result.Wheels = append(l.Result.Wheels, wheel)
		}
	}
	return nil
}

// For complex queries you can implement UnmarshalQueryOutput to control how
// the DynamoDB query results are unmarshaled.
func Example_customUnmarshalling() {
	ctx := context.TODO()

	q := ListCarAndWheelsByColor{Color: "light-orange"}
	c, _ := ddb.New(ctx, "example-table")
	_ = c.Query(ctx, &q, nil)

	// q.Result.Car and q.Result.Wheels are now populated with data as fetched from DynamoDB
}
