package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

// NewDynamoDB is to create AWS DynamoDB container and to return its connection and close function
func NewDynamoDB() (string, PurgeFunc, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, err
	}

	resource, err := pool.Run(
		"amazon/dynamodb-local",
		"latest",
		[]string{},
	)
	if err != nil {
		return "", nil, err
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"stub/aws_credentials.txt"}),
	)
	if err != nil {
		return "", nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%s", resource.GetPort("8000/tcp"))
	if err := pool.Retry(func() error {
		cl := dynamodb.New(dynamodb.Options{
			Credentials:      cfg.Credentials,
			EndpointResolver: dynamodb.EndpointResolverFromURL(endpoint),
		})
		_, err := cl.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(1),
		})
		return err
	}); err != nil {
		return "", nil, err
	}

	purgeFunc := func() error {
		if err := pool.Purge(resource); err != nil {
			return err
		}
		return nil
	}

	return endpoint, purgeFunc, nil
}
