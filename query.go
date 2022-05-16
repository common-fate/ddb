package ddb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// QueryBuilders build query inputs for DynamoDB access patterns.
// The inputs are passed to the QueryItems DynamoDB API.
//
// When writing a new QueryBuilder access pattern you should always
// implement integration tests for it against a live DynamoDB database.
type QueryBuilder interface {
	BuildQuery() (*dynamodb.QueryInput, error)
}

// QueryOutputUnmarshalers implement custom logic to
// unmarshal the results of a DynamoDB QueryItems call.
type QueryOutputUnmarshaler interface {
	UnmarshalQueryOutput(out *dynamodb.QueryOutput) error
}

// Query DynamoDB using a given QueryBuilder. Under the hood, this uses the
// QueryItems API.
//
// The QueryBuilder 'qb' defines the query, as well as how to unmarshal the
// result back into Go objects. The unmarshaling logic works as follows:
//
// 1. If qb implements UnmarshalQueryOutput, call it and return.
//
// 2. If qb contains a field with a `ddb:"result"` struct tag,
// unmarshal results to that field.
//
// 3. Unmarshal the results directly to qb.
//
// The examples in this package show how to write simple and complex access patterns
// which use each of the three methods above.
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

	var out interface{} = qb

	// check if the QueryBuilder contains a 'ddb:"result"' struct tag
	resultTag, err := findResultsTag(qb)
	if err != nil {
		return err
	}
	if resultTag != nil {
		out = resultTag.Interface()
	}

	// Otherwise, default to the unmarshalling logic provided by the attributevalue package.
	return attributevalue.UnmarshalListOfMaps(got.Items, out)
}

// findResultsTag returns the first struct field with a `ddb:"result"` tag.
func findResultsTag(out interface{}) (*reflect.Value, error) {
	v := reflect.ValueOf(out).Elem()

	if v.Kind() != reflect.Struct {
		// we can't parse this
		return nil, nil
	}

	if !v.CanAddr() {
		return nil, fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		tag, ok := f.Tag.Lookup("ddb")
		if ok && tag == "result" {
			// return nil, nil
			addr := reflect.Indirect(v).Field(i).Addr()
			return &addr, nil
		}
	}
	return nil, nil
}
