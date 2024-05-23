package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ory/dockertest/v3"
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
		return nil, fmt.Errorf("could not start resource: %w", err)
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
	return p.Pool.Retry(func() error {
		cl := dynamodb.New(dynamodb.Options{
			Credentials:  credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
			BaseEndpoint: aws.String(s.DSN),
		})
		_, err := cl.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(1),
		})
		return err
	})
}
