package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/ory/dockertest/v3"
)

type SQSFactory struct{}

func (f *SQSFactory) repository() string {
	return "softwaremill/elasticmq-native"
}

func (f *SQSFactory) create(p *Pool, opt ContainerOption) (*state, error) {
	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
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
		DSN:           fmt.Sprintf("http://localhost:%s", resource.GetPort("9324/tcp")),
		r:             resource,
	}, nil
}

func (f *SQSFactory) ready(p *Pool, s *state) error {
	return p.Pool.Retry(func() error {
		cl := sqs.New(sqs.Options{
			BaseEndpoint: aws.String(s.DSN),
		})
		_, err := cl.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
		return err
	})
}
