package dockertestx_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestNewDynamoDB(t *testing.T) {
	t.Parallel()

	endpoint, purge, err := dockertestx.NewDynamoDB()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, purge())
	})

	// for test in GitHub Actions environment
	// if omit these lines, it goes fail due to no EC2 IMDS role found
	require.NoError(t, os.Setenv("AWS_ACCESS_KEY_ID", "dummy"))
	require.NoError(t, os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy"))
	t.Cleanup(func() {
		require.NoError(t, os.Unsetenv("AWS_ACCESS_KEY_ID"))
		require.NoError(t, os.Unsetenv("AWS_SECRET_ACCESS_KEY"))
	})

	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err)

	cl := dynamodb.New(dynamodb.Options{
		Credentials:      cfg.Credentials,
		EndpointResolver: dynamodb.EndpointResolverFromURL(endpoint),
	})
	_, err = cl.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String("Test"),
		BillingMode: types.BillingModePayPerRequest,
	})
	assert.NoError(t, err)
}
