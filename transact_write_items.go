package ddb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TransactWriteItem is a wrapper over the DynamoDB TransactWriteItem type.
// Currently, only the Put option is exposed. The API supports other operations
// which can be added to this struct.
//
// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html
type TransactWriteItem struct {
	Put    Keyer
	Delete Keyer
}

func (c *Client) TransactWriteItems(ctx context.Context, tx []TransactWriteItem) error {
	twi := dynamodb.TransactWriteItemsInput{
		TransactItems: make([]types.TransactWriteItem, len(tx)),
	}

	for i := range tx {
		// a transaction must be either a put or a delete, not both.
		entry := tx[i]
		if entry.Put == nil && entry.Delete == nil {
			return errors.New("no operation defined for transaction")
		}

		if entry.Put != nil && entry.Delete != nil {
			return errors.New("both Put and Delete operations were defined for a transaction")
		}

		if entry.Put != nil {
			item, err := marshalItem(entry.Put)
			if err != nil {
				return err
			}
			twi.TransactItems[i] = types.TransactWriteItem{
				Put: &types.Put{
					Item:      item,
					TableName: &c.table,
				},
			}
		} else if entry.Delete != nil {
			keys, err := entry.Delete.DDBKeys()
			if err != nil {
				return err
			}

			keyAttrs, err := attributevalue.MarshalMap(keys)
			if err != nil {
				return err
			}
			twi.TransactItems[i] = types.TransactWriteItem{
				Delete: &types.Delete{
					Key: map[string]types.AttributeValue{
						"PK": keyAttrs["PK"],
						"SK": keyAttrs["SK"],
					},
					TableName: &c.table,
				},
			}
		}
	}

	_, err := c.client.TransactWriteItems(ctx, &twi)
	return err
}
