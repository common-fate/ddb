# ddb

[![Go Reference](https://pkg.go.dev/badge/github.com/common-fate/ddb.svg)](https://pkg.go.dev/github.com/common-fate/ddb)

Common Fate helpers for working with DynamoDB.

## Integration testing

You can provision an example table for testing as follows.

```bash
go run cmd/create/main.go
```

To run the integration tests, you need to set the `TESTING_DYNAMODB_TABLE` to be the name of the test table you created.

```bash
export TESTING_DYNAMODB_TABLE=ddb-testing
```

To cleanup the table:

```
go run cmd/destroy/main.go
```
