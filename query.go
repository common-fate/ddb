package ddb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
)

// QueryBuilders build query inputs for DynamoDB access patterns.
// The inputs are passed to the QueryItems DynamoDB API.
//
// When writing a new QueryBuilder access pattern you should always
// implement integration tests for it against a live DynamoDB database.
type QueryBuilder interface {
	BuildQuery() (*dynamodb.QueryInput, error)
}

type QueryOpts struct {
	PageToken string
	Limit     int32
}

// QueryOutputUnmarshalers implement custom logic to
// unmarshal the results of a DynamoDB QueryItems call.
type QueryOutputUnmarshaler interface {
	UnmarshalQueryOutput(out *dynamodb.QueryOutput) error
}

// Page sets the pagination token to provide an offset for the query.
// It is mapped to the 'ExclusiveStartKey' argument in the dynamodb.Query method.
func Page(pageToken string) func(*QueryOpts) {
	return func(qo *QueryOpts) {
		qo.PageToken = pageToken
	}
}

// Limit overrides the amount of items returned from the query.
// It is mapped to the 'Limit' argument in the dynamodb.Query method.
func Limit(limit int32) func(*QueryOpts) {
	return func(qo *QueryOpts) {
		qo.Limit = limit
	}
}

type QueryResult struct {
	// RawOutput is the DynamoDB API response. Usually you won't need this,
	// as results are parsed onto the QueryBuilder argument.
	RawOutput *dynamodb.QueryOutput

	// NextPage is the next page token. If empty, there is no next page.
	NextPage string
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
func (c *Client) Query(ctx context.Context, qb QueryBuilder, opts ...func(*QueryOpts)) (*QueryResult, error) {
	q, err := qb.BuildQuery()
	if err != nil {
		return nil, err
	}

	qo := QueryOpts{}
	for _, o := range opts {
		o(&qo)
	}

	// ensure that we have a tokenizer so that we don't nil panic
	if c.tokenizer == nil {
		return nil, errors.New("a page encoder must be set up to use pagination (call ddb.WithPageEncoder when setting up the client to fix)")
	}

	// set up query pagination if it's provided
	if qo.PageToken != "" {
		startKey, err := c.tokenizer.UnmarshalToken(ctx, qo.PageToken)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling page start key")
		}
		q.ExclusiveStartKey = startKey
	}

	// set the page size if it's provided
	if qo.Limit > 0 {
		q.Limit = &qo.Limit
	}

	// query builders don't necessarily know which table the client uses,
	// so update the query input to override the table name.
	q.TableName = &c.table

	got, err := c.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	result := &QueryResult{
		RawOutput: got,
	}

	// marshal the LastEvaluatedKey into a pagination token if pagination is enabled.
	if got.LastEvaluatedKey != nil {
		s, err := c.tokenizer.MarshalToken(ctx, got.LastEvaluatedKey)
		if err != nil {
			return nil, errors.Wrap(err, "marshalling LastEvaluatedKey to page token")
		}
		result.NextPage = s
	}

	// call the custom unmarshalling logic if the QueryBuilder implements it.
	if rp, ok := qb.(QueryOutputUnmarshaler); ok {
		err = rp.UnmarshalQueryOutput(got)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	var out interface{} = qb

	// check if the QueryBuilder contains a 'ddb:"result"' struct tag
	resultTag, err := findResultsTag(qb)
	if err != nil {
		return nil, err
	}
	if resultTag != nil {
		out = resultTag.Interface()
	}

	// Otherwise, default to the unmarshalling logic provided by the attributevalue package.
	err = attributevalue.UnmarshalListOfMaps(got.Items, out)
	if err != nil {
		return nil, err
	}
	return result, nil
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
