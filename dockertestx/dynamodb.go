package dockertestx

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
)

type DynamoDBFactory struct{}

func (f *DynamoDBFactory) repository() string {
	return "amazon/dynamodb-local"
}

func (f *DynamoDBFactory) create(p *Pool, opt ContainerOption) (*state, error) {
	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
		Env:        []string{},
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
		DSN:           fmt.Sprintf("http://localhost:%s", resource.GetPort("8000/tcp")),
		r:             resource,
	}, nil
}

func (f *DynamoDBFactory) ready(p *Pool, s *state) error {
	{
		clean, err := temporaryEnv("AWS_ACCESS_KEY_ID", "dummy")
		if err != nil {
			return errors.WithStack(err)
		}
		defer clean()
	}
	{
		clean, err := temporaryEnv("AWS_SECRET_ACCESS_KEY", "dummy")
		if err != nil {
			return errors.WithStack(err)
		}
		defer clean()
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return errors.WithStack(err)
	}

	return p.Pool.Retry(func() error {
		cl := dynamodb.New(dynamodb.Options{
			Credentials:      cfg.Credentials,
			EndpointResolver: dynamodb.EndpointResolverFromURL(s.DSN),
		})
		_, err := cl.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(1),
		})
		return err
	})
}

func temporaryEnv(key, value string) (func(), error) {
	v, ok := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		return nil, err
	}

	if ok {
		return func() {
			_ = os.Setenv(key, v)
		}, nil
	}
	return func() {
		_ = os.Unsetenv(key)
	}, nil
}
