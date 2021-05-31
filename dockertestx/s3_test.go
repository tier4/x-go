package dockertestx_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewS3(t *testing.T) {
	t.Parallel()

	p, err := dockertestx.New(dockertestx.PoolOption{})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, p.Purge())
	})

	endpoint, err := p.NewResource(new(dockertestx.S3Factory), dockertestx.ContainerOption{})
	require.NoError(t, err)

	cl := s3.New(s3.Options{
		Credentials:      credentials.NewStaticCredentialsProvider(dockertestx.S3AWSAccessKeyID, dockertestx.S3AWSSecretAccessKey, ""),
		EndpointResolver: s3.EndpointResolverFromURL(endpoint),
		UsePathStyle:     true,
	})
	_, err = cl.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String("test"),
	})
	assert.NoError(t, err)
}
