package dockertestx

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

// NewDynamoDB is to create AWS DynamoDB container and to return its connection and close function
func NewDynamoDB() (string, PurgeFunc, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, err
	}

	resource, err := pool.Run(
		"amazon/dynamodb-local",
		"latest",
		[]string{},
	)
	if err != nil {
		return "", nil, err
	}

	purgeFunc := func() error {
		if err := pool.Purge(resource); err != nil {
			return err
		}
		return nil
	}

	endpoint := fmt.Sprintf("http://localhost:%s", resource.GetPort("8000/tcp"))
	return endpoint, purgeFunc, nil
}
