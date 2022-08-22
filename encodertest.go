package ddb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type encoderTestCase struct {
	name string
	give map[string]types.AttributeValue
}

var encoderTestCases = []encoderTestCase{
	{
		name: "ok",
		give: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "1"},
			"SK": &types.AttributeValueMemberS{Value: "2"},
		},
	},
	{
		name: "empty",
		give: nil,
	},
}

// runEncoderTests the test suite against a PageEncoder.
func runEncoderTests(t *testing.T, e Tokenizer, testcases []encoderTestCase) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := e.MarshalToken(context.Background(), tc.give)
			if err != nil {
				t.Fatal(err)
			}
			got, err := e.UnmarshalToken(context.Background(), s)
			if err != nil {
				t.Fatal(err)
			}
			// unmarshalling should give us the original object
			assert.Equal(t, tc.give, got)
		})
	}
}
