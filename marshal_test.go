package ddb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type exampleEntityType struct{}

func (e exampleEntityType) DDBKeys() (Keys, error) {
	k := Keys{
		PK: "PK",
		SK: "SK",
	}
	return k, nil
}

func (e exampleEntityType) EntityType() string {
	return "example"
}

func Test_marshalItem(t *testing.T) {
	tests := []struct {
		name    string
		give    Keyer
		want    map[string]types.AttributeValue
		wantErr bool
	}{
		{
			name: "entity type",
			give: exampleEntityType{},
			want: map[string]types.AttributeValue{
				"PK":       &types.AttributeValueMemberS{Value: "PK"},
				"SK":       &types.AttributeValueMemberS{Value: "SK"},
				"ddb:type": &types.AttributeValueMemberS{Value: "example"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalItem(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
