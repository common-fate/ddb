package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	name = flag.String("name", "ddb-testing", "the name of the test database to create")
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := dynamodb.NewFromConfig(cfg)
	res, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: name,
	})
	if err != nil {
		return err
	}

	fmt.Printf("destroyed table, ARN = %s\n", *res.TableDescription.TableArn)
	return nil
}
