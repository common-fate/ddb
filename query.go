package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// QueryBuilders build query inputs for DynamoDB access patterns.
// The inputs are passed to the QueryItems DynamoDB API.
type QueryBuilder interface {
	BuildQuery() (*dynamodb.QueryInput, error)
}

// QueryOutputUnmarshalers implement custom logic to
// unmarshal the results of a DynamoDB QueryItems call.
type QueryOutputUnmarshaler interface {
	UnmarshalQueryOutput(out *dynamodb.QueryOutput) error
}

// Query DynamoDB using a given QueryBuilder.
// The results are unmarshaled into the QueryBuilder.
// To implement custom unmarshaling logic, implement the QueryOutputUnmarshaler
// interface on your QueryBuilder struct.
func (c *Client) Query(ctx context.Context, qb QueryBuilder) error {
	q, err := qb.BuildQuery()
	if err != nil {
		return err
	}

	// query builders don't necessarily know which table the client uses,
	// so update the query input to override the table name.
	q.TableName = &c.table

	got, err := c.client.Query(ctx, q)
	if err != nil {
		return err
	}

	// call the custom unmarshalling logic if the QueryBuilder implements it.
	if rp, ok := qb.(QueryOutputUnmarshaler); ok {
		return rp.UnmarshalQueryOutput(got)
	}

	// Otherwise, default to the unmarshalling logic provided by the attributevalue package.
	return attributevalue.UnmarshalListOfMaps(got.Items, qb)
}
