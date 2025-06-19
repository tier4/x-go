package dockertestx

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/ory/dockertest/v3"
)

type SNSSQSFactory struct{}

func (f *SNSSQSFactory) repository() string {
	return "admiralpiett/goaws"
}

func (f *SNSSQSFactory) create(p *Pool, opt ContainerOption) (*state, error) {
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
		DSN:           fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("4100/tcp")),
		r:             resource,
	}, nil
}

func (f *SNSSQSFactory) ready(p *Pool, s *state) error {
	return p.Pool.Retry(func() error {
		cl := sns.New(sns.Options{
			BaseEndpoint: aws.String(s.DSN),
		})
		_, err := cl.ListTopics(context.TODO(), &sns.ListTopicsInput{})
		return err
	})
}
