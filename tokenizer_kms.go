package ddb

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type KMSTokenizer struct {
	keyID  string
	client *kms.Client
}

func (e *KMSTokenizer) MarshalToken(ctx context.Context, item map[string]types.AttributeValue) (string, error) {
	if item == nil {
		return "", nil
	}

	b, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	input := &kms.EncryptInput{
		KeyId:     aws.String(e.keyID),
		Plaintext: b,
	}

	result, err := e.client.Encrypt(ctx, input)
	if err != nil {
		return "", err
	}
	b64EncodedToken := base64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return b64EncodedToken, nil
}

func (e *KMSTokenizer) UnmarshalToken(ctx context.Context, s string) (map[string]types.AttributeValue, error) {
	if s == "" {
		return nil, nil
	}
	b64DecodedToken, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	input := &kms.DecryptInput{
		KeyId:          aws.String(e.keyID),
		CiphertextBlob: b64DecodedToken,
	}

	result, err := e.client.Decrypt(ctx, input)
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	var tmp map[string]*types.AttributeValueMemberS
	err = json.Unmarshal([]byte(result.Plaintext), &tmp)
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

func NewKMSTokenizer(ctx context.Context, key string) (*KMSTokenizer, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := kms.NewFromConfig(cfg)

	return &KMSTokenizer{key, client}, nil
}
