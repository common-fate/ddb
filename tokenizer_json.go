package ddb

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type JSONTokenizer struct{}

func (e *JSONTokenizer) MarshalToken(ctx context.Context, item map[string]types.AttributeValue) (string, error) {
	if item == nil {
		return "", nil
	}

	b, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *JSONTokenizer) UnmarshalToken(ctx context.Context, s string) (map[string]types.AttributeValue, error) {
	if s == "" {
		return nil, nil
	}

	var tmp map[string]*types.AttributeValueMemberS
	err := json.Unmarshal([]byte(s), &tmp)
	if err != nil {
		return nil, err
	}

	// we can't use `tmp` as the return type, so copy the values over
	// into a new map of the right type.
	out := make(map[string]types.AttributeValue, len(tmp))

	for k, v := range tmp {
		out[k] = v
	}

	return out, nil
}
