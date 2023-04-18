package ddb

import (
	"context"
	"fmt"
	"reflect"
)

// Query DynamoDB using a given QueryBuilder. Under the hood, this uses Query.
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
// If the query returns a next page, then it will be loaded, until all results have been loaded from dynamodb
// The final aggregated result will be set on the querybuilder field tagged with `ddb:"result"`
func All(ctx context.Context, storage Storage, qb QueryBuilder, opts ...func(*QueryOpts)) error {
	// Get the type of the `qb` value
	qbType := reflect.TypeOf(qb).Elem()

	// Iterate over the fields in the `qb` value and find the field with the `ddb:"result"` tag
	var resultField reflect.Value
	for i := 0; i < qbType.NumField(); i++ {
		field := qbType.Field(i)
		if tag := field.Tag.Get("ddb"); tag == "result" {
			resultField = reflect.ValueOf(qb).Elem().FieldByName(field.Name)
			break
		}
	}

	if !resultField.IsValid() {
		return fmt.Errorf("could not find field with `ddb:\"result\"` tag")
	}

	// Get the type of the result field
	resultType := resultField.Type().Elem()

	// Create a new slice with the same type as the result field
	results := reflect.MakeSlice(reflect.SliceOf(resultType), 0, 0)

	optsWithPage := append([]func(*QueryOpts){}, opts...)

	for {
		// Execute the query with the given QueryBuilder and options.
		queryResult, err := storage.Query(ctx, qb, optsWithPage...)
		if err != nil {
			return err
		}
		// Iterate over the items in the `Result` field and append them to the `items` slice
		for i := 0; i < resultField.Len(); i++ {
			value := resultField.Index(i)
			convertedValue := value.Convert(resultType)
			results = reflect.Append(results, convertedValue)
		}
		// Check if there are more pages of results to fetch.
		if queryResult.NextPage == "" {
			break
		}
		// Set the next page token in the QueryOpts slice to fetch the next page of results.
		optsWithPage = append(optsWithPage, opts...)
		optsWithPage = append(optsWithPage, Page(queryResult.NextPage))
	}

	// this is where I need to set the value of resultField to results
	resultField.Set(results)
	return nil
}
