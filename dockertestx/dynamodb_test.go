package dockertestx_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewDynamoDB(t *testing.T) {
	t.Parallel()

	p, err := dockertestx.New(dockertestx.PoolOption{})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, p.Purge())
	})

	endpoint, err := p.NewResource(new(dockertestx.DynamoDBFactory), dockertestx.ContainerOption{})
	require.NoError(t, err)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"testdata/aws_credentials.txt"}),
	)
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
