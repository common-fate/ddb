package ddb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Marshal is a lightweight wrapper over attributevalue.MarshalMap
// which enforces the `json` tag key is used rather than `dynamodbav`.
func Marshal(in interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMapWithOptions(in, encoderWithJSONTagKey)
}

// Unmarshal is a lightweight wrapper over attributevalue.UnmarshalMap
// which enforces the `json` tag key is used rather than `dynamodbav`.
func Unmarshal(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMapWithOptions(m, &out, decoderWithJSONTagKey)
}

func encoderWithJSONTagKey(eo *attributevalue.EncoderOptions) {
	eo.TagKey = "json"
}

func decoderWithJSONTagKey(do *attributevalue.DecoderOptions) {
	do.TagKey = "json"
}
