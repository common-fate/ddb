package ddb

import (
	"context"
	"encoding/json"
	"fmt"

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

	result, err := EncryptText(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error encrypting data:")
		fmt.Println(err)
		return
	}

	return string(b), nil
}

func (e *KMSTokenizer) UnmarshalToken(s string) (map[string]types.AttributeValue, error) {
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

// EncryptText encrypts some text using an AWS Key Management Service (AWS KMS) customer master key (CMK).
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, an EncryptOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to Encrypt.
func EncryptText(c context.Context, api KMSEncryptAPI, input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return api.Encrypt(c, input)
}

func NewKMSTokenizer(ctx context.Context) (*KMSTokenizer, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := kms.NewFromConfig(cfg)

	return &KMSTokenizer{"", client}, nil
}
