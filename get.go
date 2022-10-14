package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetKey is the key of the item to get.
// The GetItem API always uses the primary key.
type GetKey struct {
	PK string
	SK string
}

type GetItemResult struct {
	// RawOutput is the DynamoDB API response. Usually you won't need this,
	// as results are parsed onto the item argument.
	RawOutput *dynamodb.GetItemOutput
}

// Get calls GetItem to get an item in DynamoDB.
// Get defaults to using consistent reads.
func (c *Client) Get(ctx context.Context, key GetKey, item interface{}) (*GetItemResult, error) {
	out, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: key.PK},
			"SK": &types.AttributeValueMemberS{Value: key.SK},
		},
		TableName:      &c.table,
		ConsistentRead: aws.Bool(true),
	})
	res := &GetItemResult{RawOutput: out}

	if err != nil {
		return res, err
	}
	err = attributevalue.UnmarshalMap(out.Item, item)
	return res, err
}
