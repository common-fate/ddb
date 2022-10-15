package ddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GetOpts struct {
	ConsistentRead bool
}

// GetConsistentRead customises strong read consistency.
// By default, Get() uses consistent reads.
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.ReadConsistency.html
func GetConsistentRead(enabled bool) func(*GetOpts) {
	return func(qo *GetOpts) {
		qo.ConsistentRead = enabled
	}
}

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
func (c *Client) Get(ctx context.Context, key GetKey, item Keyer, opts ...func(*GetOpts)) (*GetItemResult, error) {
	gopts := GetOpts{
		ConsistentRead: true,
	}
	for _, o := range opts {
		o(&gopts)
	}

	out, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: key.PK},
			"SK": &types.AttributeValueMemberS{Value: key.SK},
		},
		TableName:      &c.table,
		ConsistentRead: &gopts.ConsistentRead,
	})
	res := &GetItemResult{RawOutput: out}

	if err != nil {
		return res, err
	}
	err = attributevalue.UnmarshalMap(out.Item, item)
	return res, err
}
