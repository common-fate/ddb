package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
)

var (
	name = flag.String("name", "ddb-testing", "the name of the test database to create")
	wait = flag.Bool("wait", false, "wait until the table is ready")
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
	res, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		BillingMode: types.BillingModePayPerRequest,
		TableName:   name,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI2PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI2SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI1PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI1SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
			},
			{
				IndexName: aws.String("GSI2"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI2PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI2SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
			},
		},
	})
	riu := &types.ResourceInUseException{}
	if errors.As(err, &riu) {
		fmt.Printf("table %s already exists\n", *name)
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("created table, ARN = %s\n", *res.TableDescription.TableArn)

	if *wait {
		var ready bool
		for !ready {
			desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: name,
			})
			if err != nil {
				return err
			}

			if desc.Table.TableStatus == types.TableStatusActive {
				ready = true
			} else {
				sleepFor := time.Second * 3
				fmt.Printf("waiting for table to become active (status=%s)\n", desc.Table.TableStatus)
				time.Sleep(sleepFor)
			}
		}
	}

	envVar := "TESTING_DYNAMODB_TABLE"

	// write the table name to the .env file for local development.
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		fmt.Printf(".env file not found, so skipping writing %s flag\n", envVar)
		return nil
	}

	myEnv, err := godotenv.Read()
	if err != nil {
		return err
	}

	myEnv[envVar] = *name
	err = godotenv.Write(myEnv, ".env")
	if err != nil {
		return err
	}

	fmt.Printf("wrote %s to .env\n", envVar)

	return nil
}
