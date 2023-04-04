//go:build component && api

package component

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

type StockRequest struct {
	ISBN     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

type StockResponse struct {
	Id         string `json:"id"`
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}

const testTableName = "local-stock-requests"

var (
	apiUrl      = os.Getenv("API_URL")
	awsEndpoint = os.Getenv("AWS_ENDPOINT")
	awsRegion   = os.Getenv("AWS_REGION")
)

func TestAPIComponent(t *testing.T) {
	t.Run("save stock request to database", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: "../../infra/examples/api_component_test",
		})

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		client := dynamodbClient(t)

		httpResponse := requestStock(t, apiUrl, StockRequest{
			ISBN:     "isbn-123",
			Quantity: 2,
		})

		require.Equal(t, http.StatusAccepted, httpResponse.StatusCode)
		var response StockResponse
		err := json.NewDecoder(httpResponse.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Successful)
		require.NotEmpty(t, response.Id)

		item, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(testTableName),
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberS{
					Value: response.Id,
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, item)
	})
}

func dynamodbClient(t *testing.T) *dynamodb.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	require.NoError(t, err)

	return dynamodb.NewFromConfig(cfg)
}

func requestStock(t *testing.T, baseUrl string, request StockRequest) *http.Response {
	url := fmt.Sprintf("%s/stock", baseUrl)

	requestBody, err := json.Marshal(request)
	require.NoError(t, err)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	require.NoError(t, err)

	return response
}
