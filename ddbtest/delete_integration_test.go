package ddbtest

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/stretchr/testify/assert"
)

type GetThing struct {
	Type   string
	ID     string
	Result *Thing `ddb:"result"`
}

func (g *GetThing) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk and SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: g.Type},
			":sk": &types.AttributeValueMemberS{Value: g.ID},
		},
	}
	return &qi, nil
}

func (g *GetThing) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) != 1 {
		return ddb.ErrNoItems
	}

	return attributevalue.UnmarshalMap(out.Items[0], &g.Result)
}

func TestDeleteIntegration(t *testing.T) {
	c := getTestClient(t)
	ctx := context.Background()

	// insert fixture data
	a := Thing{
		Type:  randomString(20),
		ID:    randomString(20),
		Color: "red",
	}
	PutFixtures(t, c, a)

	// verify the fixture is in the table
	q := &GetThing{Type: a.Type, ID: a.ID}
	err := c.Query(ctx, q)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Delete(ctx, a)
	if err != nil {
		t.Fatal(err)
	}

	// verify the fixture has been deleted
	q = &GetThing{Type: a.Type, ID: a.ID}
	err = c.Query(ctx, q)
	assert.Equal(t, ddb.ErrNoItems, err)
}
