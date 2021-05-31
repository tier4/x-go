package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
)

const (
	S3AWSAccessKeyID     = "AKIAIOSFODNN7DUMMY"                     // #nosec
	S3AWSSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYDUMMYKEY" // #nosec
)

type S3Factory struct{}

func (f *S3Factory) repository() string {
	return "minio/minio"
}

func (f *S3Factory) create(p *Pool, opt ContainerOption) (*state, error) {
	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
		Env: []string{
			fmt.Sprintf("MINIO_ACCESS_KEY=%s", S3AWSAccessKeyID),
			fmt.Sprintf("MINIO_SECRET_KEY=%s", S3AWSSecretAccessKey),
		},
		Cmd: []string{"server", "/data"},
	}
	resource, err := p.Pool.RunWithOptions(rOpt)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not start resource")
	}
	return &state{
		ContainerName: opt.Name,
		Repository:    f.repository(),
		Tag:           opt.Tag,
		Env:           rOpt.Env,
		DSN:           fmt.Sprintf("http://localhost:%s", resource.GetPort("9000/tcp")),
		r:             resource,
	}, nil
}

func (f *S3Factory) ready(p *Pool, s *state) error {
	return p.Pool.Retry(func() error {
		cl := s3.New(s3.Options{
			Credentials:      credentials.NewStaticCredentialsProvider(S3AWSAccessKeyID, S3AWSSecretAccessKey, ""),
			EndpointResolver: s3.EndpointResolverFromURL(s.DSN),
		})
		_, err := cl.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		return err
	})
}
