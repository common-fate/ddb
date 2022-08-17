package ddb

import (
	"context"
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

func (e *KMSTokenizer) MarshalToken(item map[string]types.AttributeValue) (string, error) {
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

	result, err := e.client.Encrypt(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return string(result.CiphertextBlob), nil
}

func (e *KMSTokenizer) UnmarshalToken(s string) (map[string]types.AttributeValue, error) {
	if s == "" {
		return nil, nil
	}

	input := &kms.DecryptInput{
		KeyId:          aws.String(e.keyID),
		CiphertextBlob: []byte(s),
	}

	result, err := e.client.Decrypt(context.TODO(), input)
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
	var opts []func(*config.LoadOptions) error

	opts = append(opts, config.WithRegion("ap-southeast-2"))

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	client := kms.NewFromConfig(cfg)

	return &KMSTokenizer{key, client}, nil
}
