package ddb

import (
	"testing"
)

func TestJSONEncoder(t *testing.T) {
	runEncoderTests(t, &JSONTokenizer{}, encoderTestCases)
}
