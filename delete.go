package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Delete calls DeleteItem to delete an item in DynamoDB.
func (c *Client) Delete(ctx context.Context, item Keyer) error {
	keys, err := item.DDBKeys()
	if err != nil {
		return err
	}

	keyAttrs, err := attributevalue.MarshalMap(keys)
	if err != nil {
		return err
	}

	_, err = c.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"PK": keyAttrs["PK"],
			"SK": keyAttrs["SK"],
		},
		TableName: &c.table,
	})
	return err
}

// DeleteBatch calls BatchWriteItem to create or update items in DynamoDB.
//
// DynamoDB BatchWriteItem api has a limit of 25 items per batch.
// DeleteBatch will automatically split the items into batches of 25 by default.
//
// You can override this default batch size using WithBatchSize(n) when you initialize the client.
func (c *Client) DeleteBatch(ctx context.Context, items ...Keyer) error {
	wr := make([]types.WriteRequest, len(items))
	for i, item := range items {
		keys, err := item.DDBKeys()
		if err != nil {
			return err
		}

		keyAttrs, err := attributevalue.MarshalMap(keys)
		if err != nil {
			return err
		}
		wr[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": keyAttrs["PK"],
					"SK": keyAttrs["SK"],
				},
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
