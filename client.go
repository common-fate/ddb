package ddb

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// Client is a thin wrapper over the native DynamoDB client.
// It has methods which allow access patterns to be written
// in a more ergonomic fashion than the native client.
type Client struct {
	table  string
	client dynamodb.Client
}
