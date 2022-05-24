package ddb

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// marshalItem turns an item into it's DynamoDB representation.
func marshalItem(item Keyer) (map[string]types.AttributeValue, error) {
	keys, err := item.DDBKeys()
	if err != nil {
		return nil, err
	}

	// marshal the object itself
	objAttrs, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(keys)
	// add the keys to the object
	for i := 0; i < v.NumField(); i++ {
		k := v.Type().Field(i).Name
		val := v.Field(i).String()

		// any fields which are empty strings are useless to write to DynamoDB.
		// when iterating through the object we ignore these.
		if val != "" {
			objAttrs[k] = &types.AttributeValueMemberS{Value: val}
		}
	}

	return objAttrs, nil
}
