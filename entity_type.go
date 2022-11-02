package ddb

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// If EntityType is implemented, ddb will write a special
// 'ddb:type' when marshalling the object.
// This can be used to implement custom unmarshalling for
// queries which return multiple object types.
//
// For example, you may wish to save an invoice along with
// its line items as separate rows in DynamoDB.
// The `EntityType` of the invoice can be "invoice",
// and then `EntityType` of the line item can be "lineItem".
// When querying the database and unmarshalling these objects
// back into Go structs, you can check the type of them
// by looking at the value of 'ddb:type'.
type EntityTyper interface {
	EntityType() string
}

var ErrNoEntityType = errors.New("item does not have a 'ddb:type' field")

// GetItemEntityType gets the entity type of a raw DynamoDB item.
// The item must have a 'ddb:type' field on it.
// If it doesn't, an error is returned.
//
// To use this, implement ddb.EntityTyper on your objects.
func GetItemEntityType(item map[string]types.AttributeValue) (string, error) {
	var tmp struct {
		Type string `dynamodbav:"ddb:type"`
	}

	err := attributevalue.UnmarshalMap(item, &tmp)
	if err != nil {
		return "", err
	}

	if tmp.Type == "" {
		return "", ErrNoEntityType
	}
	return tmp.Type, nil
}
