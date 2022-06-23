package ddbtest

import (
	"context"
	"crypto/rand"
	"math/big"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// randomString is used to generate random IDs for integration testing
func randomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

type Thing struct {
	// a random type to prevent queries from overlapping with one another
	Type  string
	ID    string
	Color string
}

// DDBKeys returns the item's keys for DynamoDB.
// We set both PK/SK as well as a GSI so that we can test the behaviour of
// access patterns that use a GSI.
func (a Thing) DDBKeys() (ddb.Keys, error) {
	return ddb.Keys{PK: a.Type, SK: a.ID, GSI1PK: a.Type, GSI1SK: a.ID}, nil
}

type ListThingStructTag struct {
	Type   string
	Result []Thing `ddb:"result"`
}

func (l *ListThingStructTag) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}

type ListThingsCustomUnmarshal struct {
	Type   string
	Result []Thing
}

func (l *ListThingsCustomUnmarshal) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}

func (l *ListThingsCustomUnmarshal) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	return attributevalue.UnmarshalListOfMaps(out.Items, &l.Result)
}

// Things is a helper method used for asserting the result value in testing.
func (l *ListThingsCustomUnmarshal) Things() []Thing {
	return l.Result
}

type ListThingGSI struct {
	Type   string
	Result []Thing `ddb:"result"`
}

// Things is a helper method used for asserting the result value in testing.
func (l *ListThingGSI) Things() []Thing {
	return l.Result
}

func (l *ListThingGSI) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		KeyConditionExpression: aws.String("GSI1PK = :pk"),
		IndexName:              aws.String("GSI1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: l.Type},
		},
	}
	return &qi, nil
}

func randomThings(t string, count int) []Thing {
	var apples []Thing
	for i := 0; i < count; i++ {
		a := Thing{
			Type:  t,
			ID:    randomString(20),
			Color: "red",
		}
		apples = append(apples, a)
	}
	sort.SliceStable(apples, func(i, j int) bool {
		a := apples[i]
		b := apples[j]
		return a.ID < b.ID
	})
	return apples
}

func TestQueryIntegration(t *testing.T) {
	_ = godotenv.Load("../.env")
	c := getTestClient(t)

	// prevent queries from overlapping with any existing fixture data
	// in the database by using a random string as a type
	typ := randomString(20)

	// insert fixture data
	apples := randomThings(typ, 5)
	PutFixtures(t, c, apples)

	testcases := []QueryTestCase{
		{
			Name:  "struct tag",
			Query: &ListThingStructTag{Type: typ},
			Want:  &ListThingStructTag{Type: typ, Result: apples},
		},
		{
			Name:  "custom unmarshal",
			Query: &ListThingsCustomUnmarshal{Type: typ},
			Want:  &ListThingsCustomUnmarshal{Type: typ, Result: apples},
		},
	}
	RunQueryTests(t, c, testcases)
}

// TestPaginationIntegration tests that pagination works as expected against a live
// DynamoDB database.
//
// The test runs by running two queries with the ddb.Page() option set.
// We assert that the first query returns the first page (matching 'WantPage1' in the test case)
// and that the second query returns the second page (match 'WantPage2' in the test case).
func TestPaginationIntegration(t *testing.T) {
	_ = godotenv.Load("../.env")
	c := getTestClient(t, ddb.WithPageTokenizer(&ddb.JSONTokenizer{}))

	// prevent queries from overlapping with any existing fixture data
	// in the database by using a random string as a type
	typ := randomString(20)

	// insert fixture data
	apples := randomThings(typ, 5)
	PutFixtures(t, c, apples)

	type testcase struct {
		Name      string
		Query     ddb.QueryBuilder
		PageSize  int32
		WantPage1 ddb.QueryBuilder
		WantPage2 ddb.QueryBuilder
	}

	testcases := []testcase{
		{
			// check that pagination works on a GSI
			Name:      "pagination gsi",
			Query:     &ListThingGSI{Type: typ},
			PageSize:  1,
			WantPage1: &ListThingGSI{Type: typ, Result: apples[0:1]},
			WantPage2: &ListThingGSI{Type: typ, Result: apples[1:2]},
		},
		{
			Name:      "pagination",
			Query:     &ListThingsCustomUnmarshal{Type: typ},
			PageSize:  1,
			WantPage1: &ListThingsCustomUnmarshal{Type: typ, Result: apples[0:1]},
			WantPage2: &ListThingsCustomUnmarshal{Type: typ, Result: apples[1:2]},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			// run the query for the first page
			res, err := c.Query(context.Background(), tc.Query, ddb.Limit(1))
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.WantPage1, tc.Query)

			// run a second query for the second page
			_, err = c.Query(context.Background(), tc.Query, ddb.Limit(1), ddb.Page(res.NextPage))
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.WantPage2, tc.Query)
		})
	}
}

// ensure that the pagination token is overwritten to an empty string
// if we run out of pages.
func TestPaginationTokenEmptyIntegration(t *testing.T) {
	_ = godotenv.Load("../.env")
	c := getTestClient(t, ddb.WithPageTokenizer(&ddb.JSONTokenizer{}))

	// prevent queries from overlapping with any existing fixture data
	// in the database by using a random string as a type
	typ := randomString(20)

	// insert 2 items as fixture data
	apples := randomThings(typ, 2)
	PutFixtures(t, c, apples)

	q := &ListThingGSI{Type: typ}

	// run the query for the first page
	res, err := c.Query(context.Background(), q, ddb.Limit(1))
	if err != nil {
		t.Fatal(err)
	}
	// token shouldn't be empty, as we use it for the next page
	assert.NotEmpty(t, res.NextPage)

	// run a second query for the second page
	res, err = c.Query(context.Background(), q, ddb.Limit(1), ddb.Page(res.NextPage))
	if err != nil {
		t.Fatal(err)
	}

	// token should now be empty, as we've run out of items
	assert.Empty(t, res.NextPage)
}
