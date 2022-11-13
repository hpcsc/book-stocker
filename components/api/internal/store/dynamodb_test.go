//go:build integration

package store

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/hpcsc/book-stocker/api/internal/config"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

var (
	awsEndpoint = os.Getenv("AWS_ENDPOINT")
	awsRegion   = os.Getenv("AWS_REGION")
)

const testTableName = "local-stock-requests"

func TestDynamodbStore_Save(t *testing.T) {
	t.Run("create when stock request not exists", func(t *testing.T) {
		cfg := config.Configuration{
			AWSEndpoint: awsEndpoint,
			AWSRegion:   awsRegion,
			Environment: "local",
		}

		client := testDynamoDbClient(t, &cfg)
		ensureTestTableExists(t, client)

		requestStore, err := NewDynamoDbStore(&cfg)
		require.NoError(t, err)

		id := uuid.New().String()
		err = requestStore.Save(context.TODO(), StockRequest{
			Id:       id,
			ISBN:     "some-isbn",
			Quantity: 2,
		})
		require.NoError(t, err)

		existingItem := stockRequestWithId(t, client, id)
		require.Equal(t, "some-isbn", existingItem.ISBN)
		require.Equal(t, 2, existingItem.Quantity)
	})

	t.Run("override when stock request with the same id exists", func(t *testing.T) {
		cfg := config.Configuration{
			AWSEndpoint: awsEndpoint,
			AWSRegion:   awsRegion,
			Environment: "local",
		}
		requestStore, err := NewDynamoDbStore(&cfg)
		require.NoError(t, err)

		client := testDynamoDbClient(t, &cfg)
		ensureTestTableExists(t, client)

		id := uuid.New().String()
		putTestStockRequest(t, client, StockRequest{
			Id:       id,
			ISBN:     "old-isbn",
			Quantity: 1,
		})

		err = requestStore.Save(context.TODO(), StockRequest{
			Id:       id,
			ISBN:     "some-isbn",
			Quantity: 2,
		})
		require.NoError(t, err)

		existingItem := stockRequestWithId(t, client, id)
		require.Equal(t, "some-isbn", existingItem.ISBN)
		require.Equal(t, 2, existingItem.Quantity)
	})
}

func ensureTestTableExists(t *testing.T, client *dynamodb.Client) {
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(testTableName),
	})
	if err != nil && !strings.Contains(err.Error(), "resource not found") {
		require.Fail(t, err.Error())
	}

	_, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	require.NoError(t, err)
}

func putTestStockRequest(t *testing.T, client *dynamodb.Client, stockRequest StockRequest) {
	item, err := attributevalue.MarshalMap(stockRequest)
	require.NoError(t, err)

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(testTableName),
		Item:      item,
	})
	require.NoError(t, err)
}

func testDynamoDbClient(t *testing.T, cfg *config.Configuration) *dynamodb.Client {
	dbConfig, err := dynamoDbConfig(cfg)
	require.NoError(t, err)
	client := dynamodb.NewFromConfig(dbConfig)
	return client
}

func stockRequestWithId(t *testing.T, client *dynamodb.Client, key string) StockRequest {
	item, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{
				Value: key,
			},
		},
	})
	require.NoError(t, err)

	var stock StockRequest
	require.NoError(t, attributevalue.UnmarshalMap(item.Item, &stock))
	return stock
}
