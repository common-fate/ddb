package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Put calls PutItem to create or update an item in DynamoDB.
func (c *Client) Put(ctx context.Context, item Keyer) error {
	attrs, err := marshalItem(item)
	if err != nil {
		return err
	}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      attrs,
		TableName: &c.table,
	})
	return err
}

// PutBatch calls BatchWriteItem to create or update items in DynamoDB.
//
// DynamoDB BatchWriteItem api has a limit of 25 items per batch.
// PutBatch will automatically split the items into batches of 25 by default.
//
// You can override this default batch size using WithBatchSize(n) when you initialize the client.
func (c *Client) PutBatch(ctx context.Context, items ...Keyer) error {
	wr := make([]types.WriteRequest, len(items))
	for i, item := range items {
		dbItem, err := marshalItem(item)
		if err != nil {
			return err
		}
		wr[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: dbItem,
			},
		}
	}
	for i := 0; i < len(wr); i += c.batchSize {
		end := len(wr)
		if i+c.batchSize < end {
			end = i + c.batchSize
		}
		_, err := c.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				c.table: wr[i:end],
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
