package dockertestx_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewSQS(t *testing.T) {
	t.Parallel()

	p, err := dockertestx.New(dockertestx.PoolOption{})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, p.Purge())
	})

	endpoint, err := p.NewResource(new(dockertestx.SQSFactory), dockertestx.ContainerOption{})
	require.NoError(t, err)

	cl := sqs.New(sqs.Options{
		BaseEndpoint: aws.String(endpoint),
	})
	_, err = cl.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: aws.String("test"),
	})
	assert.NoError(t, err)
}
