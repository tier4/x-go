package dockertestx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ory/dockertest/v3"
)

type PrismFactory struct {
	// SpecURI is a file path or URL of the OpenAPI Specification.
	SpecURI string

	// HealthCheckPath is the path accessed to verify that the stub server has started.
	// The response status code is ignored.
	// The default is to use the base path.
	HealthCheckPath string
}

func (f *PrismFactory) repository() string {
	return "stoplight/prism"
}

func (f *PrismFactory) create(p *Pool, opt ContainerOption) (*state, error) {
	if opt.Tag == "latest" {
		opt.Tag = "5"
	}

	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
		Env:        []string{},
		Cmd:        []string{"mock", "--host", "0.0.0.0"},
	}

	if u, err := url.Parse(f.SpecURI); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		rOpt.Cmd = append(rOpt.Cmd, f.SpecURI)
	} else {
		fp, err := filepath.Abs(f.SpecURI)
		if err != nil {
			return nil, fmt.Errorf("could not resolve abstract path of the definition file: %w", err)
		}
		if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("could not stat the definition file: %w", err)
		}
		rOpt.Mounts = []string{fp + ":/tmp/oas.yml:ro"}
		rOpt.Cmd = append(rOpt.Cmd, "/tmp/oas.yml")
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
		// Use 127.0.0.1 explicitly to avoid IPv6 resolution of "localhost" in CI.
		DSN: fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("4010/tcp")),
		r:   resource,
	}, nil
}

func (f *PrismFactory) ready(p *Pool, s *state) error {
	u, err := url.JoinPath(s.DSN, f.HealthCheckPath)
	if err != nil {
		return fmt.Errorf("invalid heath check path: %w", err)
	}
	// Prism (a Node.js app) takes longer than the default 1-minute MaxWait to start in CI.
	p.Pool.MaxWait = 3 * time.Minute
	return p.Pool.Retry(func() error {
		// Fail immediately if the container has already exited.
		c, err := p.Pool.Client.InspectContainer(s.r.Container.ID)
		if err == nil && !c.State.Running {
			return backoff.Permanent(fmt.Errorf("prism container exited with code %d", c.State.ExitCode))
		}

		out, err := http.Get(u) // #nosec G107
		if err != nil {
			return err
		}
		defer out.Body.Close()
		_, _ = io.ReadAll(out.Body)
		return nil
	})
}
