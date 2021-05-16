package dockertestx

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"

	"github.com/tier4/x-go/idx"
)

// Pool represents a connection to the docker API and is used to create and remove docker images.
type Pool struct {
	*dockertest.Pool
	states stateList

	option PoolOption
}

type PoolOption struct {
	// KeepContainer or not
	// if true, write container state to StateStore when calling Pool.Save()
	KeepContainer bool

	// StateStore of docker container state for using container
	// this must be present when KeepContainer is true
	StateStore *os.File
}

func (o *PoolOption) validate() error {
	if !o.KeepContainer {
		return nil
	}
	if o.StateStore == nil {
		return errors.New("PoolOption: StateStore must be present when KeepContainer is true")
	}
	return nil
}

// New Pool instance
func New(opt PoolOption) (*Pool, error) {
	if err := opt.validate(); err != nil {
		return nil, errors.WithStack(err)
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.WithMessage(err, "Could not connect to docker")
	}

	states := make(stateList, 0, 3)

	if opt.StateStore != nil {
		ds := make(stateList, 0, 3)
		_ = json.NewDecoder(opt.StateStore).Decode(&ds)
		for i := range ds {
			ds[i].r, _ = pool.ContainerByName(ds[i].ContainerName)
			if ds[i].r != nil {
				// ignore not found container
				states = append(states, ds[i])
			}
		}
	}

	return &Pool{
		Pool:   pool,
		states: states,
		option: opt,
	}, nil
}

// ForcePurge regardless KeepContainer option
func (p *Pool) ForcePurge() error {
	var merr error
	for _, s := range p.states {
		if s.r == nil {
			continue
		}
		if err := p.Pool.Purge(s.r); err != nil {
			merr = multierror.Append(merr, err)
		}
	}
	return merr
}

// Purge if KeepContainer option is false
func (p *Pool) Purge() error {
	if p.option.KeepContainer {
		return nil
	}
	return p.ForcePurge()
}

// Save states to reuse containers next time
func (p *Pool) Save() error {
	if p.option.StateStore == nil {
		return nil
	}

	s := make([]state, 0)
	if p.option.KeepContainer {
		s = p.states
	}
	if err := p.option.StateStore.Truncate(0); err != nil {
		return errors.WithMessage(err, "Could not truncate container state")
	}
	if _, err := p.option.StateStore.Seek(0, 0); err != nil {
		return errors.WithMessage(err, "Could not change offset in state file")
	}
	if err := json.NewEncoder(p.option.StateStore).Encode(&s); err != nil {
		return errors.WithMessage(err, "Could not save container state")
	}
	return nil
}

type stateList []state

func (l stateList) find(name, repository, tag string) (*state, bool) {
	if len(name) > 0 {
		for _, s := range l {
			if s.ContainerName == name {
				return &s, true
			}
		}
		return nil, false
	}

	for _, s := range l {
		if s.Repository == repository && s.Tag == tag {
			return &s, true
		}
	}
	return nil, false
}

type state struct {
	ContainerName string   `json:"container_name"`
	Repository    string   `json:"repository"`
	Tag           string   `json:"tag"`
	Env           []string `json:"env"`
	DSN           string   `json:"dsn"`

	r *dockertest.Resource
}

type ContainerOption struct {
	// Name of docker container
	// it must be present when use same factories
	Name string

	// Tag of docker image
	Tag string
}

// ContainerFactory represents how to create docker container
type ContainerFactory interface {
	// repository of docker image name
	repository() string

	// create docker container in pool
	create(p *Pool, opt ContainerOption) (*state, error)

	// ready to access
	ready(p *Pool, s *state) error
}

func (p *Pool) NewResource(factory ContainerFactory, opt ContainerOption) (string, error) {
	if len(opt.Tag) == 0 {
		opt.Tag = "latest"
	}

	s, ok := p.states.find(opt.Name, factory.repository(), opt.Tag)
	if !ok {
		if len(opt.Name) == 0 {
			const namePrefix = "dockertestx"
			opt.Name = fmt.Sprintf("%s_%s", namePrefix, idx.ShortID())
		}

		var err error
		s, err = factory.create(p, opt)
		if err != nil {
			return "", errors.WithStack(err)
		}
		p.states = append(p.states, *s)
	}
	if err := factory.ready(p, s); err != nil {
		return "", errors.WithStack(err)
	}
	return s.DSN, nil
}
