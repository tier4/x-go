package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
)

const (
	S3AWSAccessKeyID     = "AKIAIOSFODNN7DUMMY"                     // #nosec
	S3AWSSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYDUMMYKEY" // #nosec
)

type S3Factory struct{}

func (f *S3Factory) repository() string {
	return "versity/versitygw"
}

func (f *S3Factory) create(p *Pool, opt ContainerOption) (*state, error) {
	rOpt := &dockertest.RunOptions{
		Name:         opt.Name,
		Repository:   f.repository(),
		Tag:          opt.Tag,
		ExposedPorts: []string{"7070/tcp"},
		Env: []string{
			"VGW_BACKEND=posix",
			"VGW_BACKEND_ARG=/tmp",
			fmt.Sprintf("ROOT_ACCESS_KEY_ID=%s", S3AWSAccessKeyID),
			fmt.Sprintf("ROOT_SECRET_ACCESS_KEY=%s", S3AWSSecretAccessKey),
		},
	}
	resource, err := p.RunWithOptions(rOpt)
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %w", err)
	}
	return &state{
		ContainerName: opt.Name,
		Repository:    f.repository(),
		Tag:           opt.Tag,
		Env:           rOpt.Env,
		DSN:           fmt.Sprintf("http://localhost:%s", resource.GetPort("7070/tcp")),
		r:             resource,
	}, nil
}

func (f *S3Factory) ready(p *Pool, s *state) error {
	return p.Retry(func() error {
		cl := s3.New(s3.Options{
			Credentials:  credentials.NewStaticCredentialsProvider(S3AWSAccessKeyID, S3AWSSecretAccessKey, ""),
			BaseEndpoint: aws.String(s.DSN),
			Region:       "us-east-1",
		})
		_, err := cl.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		return err
	})
}
