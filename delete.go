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
