package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	tableName  = aws.String("thoth")
	bucketName = aws.String("thoth-games")
)

type conn struct {
	db     *dynamodb.Client
	bucket *s3.Client
}

func NewConnection() *conn {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	bucket := s3.NewFromConfig(cfg)

	return &conn{
		db:     svc,
		bucket: bucket,
	}
}
