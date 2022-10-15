package ddb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestGetItemEntityType(t *testing.T) {
	tests := []struct {
		name    string
		give    map[string]types.AttributeValue
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			give: map[string]types.AttributeValue{
				"ddb:type": &types.AttributeValueMemberS{Value: "test"},
			},
			want: "test",
		},
		{
			name: "multiple fields",
			give: map[string]types.AttributeValue{
				"ddb:type": &types.AttributeValueMemberS{Value: "test"},
				"type":     &types.AttributeValueMemberS{Value: "other"},
				"PK":       &types.AttributeValueMemberS{Value: "other"},
				"SK":       &types.AttributeValueMemberS{Value: "other"},
			},
			want: "test",
		},
		{
			name: "not found",
			give: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: "other"},
				"SK": &types.AttributeValueMemberS{Value: "other"},
			},
			wantErr: true,
		},
		{
			name:    "empty object",
			give:    map[string]types.AttributeValue{},
			wantErr: true,
		},
		{
			name:    "nil object",
			give:    nil,
			wantErr: true,
		},
		{
			name: "wrong type",
			give: map[string]types.AttributeValue{
				"ddb:type": &types.AttributeValueMemberBOOL{Value: true},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetItemEntityType(tt.give)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemEntityType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetItemEntityType() = %v, want %v", got, tt.want)
			}
		})
	}
}
