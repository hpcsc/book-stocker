package store

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hpcsc/book-stocker/api/internal/config"
)

type dynamodbStore struct {
	client *dynamodb.Client
	cfg    *config.Configuration
}

func NewDynamoDbStore(cfg *config.Configuration) (Interface, error) {
	dbConfig, err := dynamoDbConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb config: %v", err)
	}

	return &dynamodbStore{
		client: dynamodb.NewFromConfig(dbConfig),
		cfg:    cfg,
	}, nil
}

func (s *dynamodbStore) Save(ctx context.Context, purchase StockRequest) error {
	marshalled, err := attributevalue.MarshalMap(purchase)
	if err != nil {
		return fmt.Errorf("failed to marshal purchase: %v", err)
	}

	if _, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(fmt.Sprintf("%s-stock-requests", s.cfg.Environment)),
		Item:      marshalled,
	}); err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	return nil
}

func dynamoDbConfig(cfg *config.Configuration) (aws.Config, error) {
	if cfg.AWSEndpoint == "" {
		return awsconfig.LoadDefaultConfig(context.TODO())
	}

	return awsconfig.LoadDefaultConfig(
		context.TODO(),
		endpointConfig(cfg),
	)
}

func endpointConfig(cfg *config.Configuration) awsconfig.LoadOptionsFunc {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.AWSEndpoint,
				SigningRegion: cfg.AWSRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	return awsconfig.WithEndpointResolverWithOptions(customResolver)
}
