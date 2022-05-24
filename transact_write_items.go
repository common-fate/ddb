package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TransactWriteItem is a wrapper over the DynamoDB TransactWriteItem type.
// Currently, only the Put option is exposed. The API supports other operations
// which can be added to this struct.
//
// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/transaction-apis.html
type TransactWriteItem struct {
	Put Keyer
}

func (c *Client) TransactWriteItems(ctx context.Context, tx []TransactWriteItem) error {
	twi := dynamodb.TransactWriteItemsInput{
		TransactItems: make([]types.TransactWriteItem, len(tx)),
	}

	for i := range tx {
		item, err := marshalItem(tx[i].Put)
		if err != nil {
			return err
		}
		twi.TransactItems[i] = types.TransactWriteItem{
			Put: &types.Put{
				Item:      item,
				TableName: &c.table,
			},
		}
	}

	_, err := c.client.TransactWriteItems(ctx, &twi)
	return err
}
