package ddb

import (
	"math/rand"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randomString is used to generate random IDs for integration testing
func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Thing struct {
	// a random type to prevent queries from overlapping with one another
	Type  string
	ID    string
	Color string
}

func (a Thing) DDBKeys() (Keys, error) {
	return Keys{PK: a.Type, SK: a.ID}, nil
}

type ListThingStructTag struct {
	Type   string
	Result []Thing `ddb:"result"`
}

func (l *ListThingStructTag) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}

type ListThingsCustomUnmarshal struct {
	Type   string
	Result []Thing
}

func (l *ListThingsCustomUnmarshal) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}

func (l *ListThingsCustomUnmarshal) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	return attributevalue.UnmarshalListOfMaps(out.Items, &l.Result)
}

func randomThings(t string) []Thing {
	var apples []Thing
	for i := 0; i < 5; i++ {
		a := Thing{
			Type:  t,
			ID:    randomString(20),
			Color: "red",
		}
		apples = append(apples, a)
	}
	return apples
}

func TestQueryIntegration(t *testing.T) {
	c := getTestClient(t)

	// prevent queries from overlapping with any existing fixture data
	// in the database by using a random string as a type
	typ := randomString(20)

	// insert fixture data
	apples := randomThings(typ)
	PutFixtures(t, c, apples)

	testcases := []QueryTestCase{
		{
			Name:  "struct tag",
			Query: &ListThingStructTag{Type: typ},
			Want:  &ListThingStructTag{Type: typ, Result: apples},
		},
		{
			Name:  "custom unmarshal",
			Query: &ListThingsCustomUnmarshal{Type: typ},
			Want:  &ListThingsCustomUnmarshal{Type: typ, Result: apples},
		},
	}
	RunQueryTests(t, c, testcases)
}
