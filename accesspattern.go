package ddb

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// AccessPatterns build and parse DynamoDB queries into Go types.
type AccessPattern interface {
	Build() (*dynamodb.QueryInput, error)
	ParseResults(out *dynamodb.QueryOutput) error
}
