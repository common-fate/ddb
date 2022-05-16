package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Client is a thin wrapper over the native DynamoDB client.
// It has methods which allow access patterns to be written
// in a more ergonomic fashion than the native client.
type Client struct {
	table  string
	client *dynamodb.Client
}

// New creates a new DynamoDB Client.
func New(ctx context.Context, table string, opts ...func(*Client)) (*Client, error) {
	c := &Client{
		table: table,
	}

	for _, o := range opts {
		o(c)
	}

	if c.client == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, err
		}
		c.client = dynamodb.NewFromConfig(cfg)
	}

	return c, nil
}

// WithDynamoDBClient allows a custom dynamodb.Client to be provided.
func WithDynamoDBClient(d *dynamodb.Client) func(*Client) {
	return func(c *Client) {
		c.client = d
	}
}
