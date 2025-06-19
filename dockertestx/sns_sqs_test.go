package dockertestx_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_SNSSQSFactory(t *testing.T) {
	t.Parallel()

	p, err := dockertestx.New(dockertestx.PoolOption{})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, p.Purge())
	})

	endpoint, err := p.NewResource(new(dockertestx.SNSSQSFactory), dockertestx.ContainerOption{})
	require.NoError(t, err)

	snsCli := sns.New(sns.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       "ap-northeast-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
			}, nil
		}),
	})

	sqsCli := sqs.New(sqs.Options{
		BaseEndpoint: aws.String(endpoint),
		Region:       "ap-northeast-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
			}, nil
		}),
	})

	ctx := context.Background()

	var (
		topicArn string
		queueUrl string
	)

	// Create a topic and a queue
	{
		createTopicRes, err := snsCli.CreateTopic(ctx, &sns.CreateTopicInput{
			Name: aws.String("test-topic.fifo"),
			Attributes: map[string]string{
				"FifoTopic": "true",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, createTopicRes.TopicArn)

		createQueueRes, err := sqsCli.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String("test-queue.fifo"),
			Attributes: map[string]string{
				"FifoQueue": "true",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, createQueueRes.QueueUrl)

		topicArn = *createTopicRes.TopicArn
		queueUrl = *createQueueRes.QueueUrl
	}

	// Subscribe the queue to the topic
	_, err = snsCli.Subscribe(ctx, &sns.SubscribeInput{
		Protocol: aws.String("sqs"),
		Attributes: map[string]string{
			"RawMessageDelivery": "true",
		},
		TopicArn: aws.String(topicArn),
		Endpoint: aws.String(queueUrl),
	})
	require.NoError(t, err)

	// Publish a message to the topic
	_, err = snsCli.Publish(ctx, &sns.PublishInput{
		Message:        aws.String("Hello, SNS!"),
		TopicArn:       aws.String(topicArn),
		MessageGroupId: aws.String("test-group"),
	})
	require.NoError(t, err)

	// Receive the message from the queue
	receiveRes, err := sqsCli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:                    aws.String(queueUrl),
		MaxNumberOfMessages:         1,
		MessageAttributeNames:       []string{"All"},
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll},
		ReceiveRequestAttemptId:     nil,
		VisibilityTimeout:           1,
		WaitTimeSeconds:             1,
	})
	require.NoError(t, err)
	assert.Len(t, receiveRes.Messages, 1)
	assert.Equal(t, "Hello, SNS!", aws.ToString(receiveRes.Messages[0].Body))

	// Clean up: delete the queue and topic
	_, err = sqsCli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: receiveRes.Messages[0].ReceiptHandle,
	})
	require.NoError(t, err)
}
