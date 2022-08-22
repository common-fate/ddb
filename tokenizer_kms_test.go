package ddb

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestKMSEncoder(t *testing.T) {
	_ = godotenv.Load("./.env")
	key := os.Getenv("AWS_KMS_KEY")
	if key == "" {
		t.Skip("AWS_KMS_KEY is not set")
	}
	kmsTokenizer, err := NewKMSTokenizer(context.Background(), key)
	if err != nil {
		t.Errorf(err.Error())
	}
	runEncoderTests(t, kmsTokenizer, encoderTestCases)
}
