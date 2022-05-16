# ddb

[![Go Reference](https://pkg.go.dev/badge/github.com/common-fate/ddb.svg)](https://pkg.go.dev/github.com/common-fate/ddb)

Common Fate helpers for working with DynamoDB.

## Integration testing

You can provision an example table for testing as follows.

```bash
aws dynamodb create-table --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE --billing-mode=PAY_PER_REQUEST  --table-name ddb-testing
```

To run the integration tests, you need to set the `TESTING_DYNAMODB_TABLE` to be the name of the test table you created.

```bash
export TESTING_DYNAMODB_TABLE=ddb-testing
```

To cleanup the table:

```
aws dynamodb delete-table --table-name ddb-testing
```
