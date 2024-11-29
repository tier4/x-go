package dockertestx_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

type container struct {
	name        string
	factory     dockertestx.ContainerFactory
	tag         string
	previousDSN string
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("reuse container when same factory and no name given", func(t *testing.T) {
		t.Parallel()

		statePath := "testdata/cannot_use_same_factory_without_name.tmp.json"
		f, err := os.Create(statePath)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = f.Close()
			_ = os.Remove(statePath)
		})

		p, err := dockertestx.New(dockertestx.PoolOption{
			KeepContainer: true,
			StateStore:    f,
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = p.ForcePurge()
		})

		dsn1, err := p.NewResource(new(dockertestx.PostgresFactory), dockertestx.ContainerOption{
			Tag: "alpine",
		})
		require.NoError(t, err)
		dsn2, err := p.NewResource(new(dockertestx.PostgresFactory), dockertestx.ContainerOption{
			Tag: "alpine",
		})
		require.NoError(t, err)
		assert.Equal(t, dsn1, dsn2)
	})

	t.Run("reuse multiple containers", func(t *testing.T) {
		t.Parallel()

		statePath := "testdata/reuse_multiple_containers.tmp.json"
		t.Cleanup(func() {
			_ = os.Remove(statePath)
		})

		cs := []container{
			{
				name:    fmt.Sprintf("dockertestx_postgres_01_%s", dockertestx.ShortID()),
				factory: new(dockertestx.PostgresFactory),
				tag:     "alpine",
			},
			{
				name:    fmt.Sprintf("dockertestx_postgres_02_%s", dockertestx.ShortID()),
				factory: new(dockertestx.PostgresFactory),
				tag:     "alpine",
			},
			{
				name:    fmt.Sprintf("dockertestx_dynamodb_01_%s", dockertestx.ShortID()),
				factory: new(dockertestx.DynamoDBFactory),
				tag:     "latest",
			},
		}

		for _, c := range []struct {
			name              string
			fileFlag          int
			keepContainer     bool
			expectStateLength int
		}{
			{
				name:              "new state",
				fileFlag:          os.O_RDWR | os.O_CREATE | os.O_TRUNC,
				keepContainer:     true,
				expectStateLength: len(cs),
			},
			{
				name:              "reuse/leave state",
				fileFlag:          os.O_RDWR,
				keepContainer:     true,
				expectStateLength: len(cs),
			},
			{
				name:              "reuse/clean state",
				fileFlag:          os.O_RDWR,
				keepContainer:     false,
				expectStateLength: 0,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				f, err := os.OpenFile(statePath, c.fileFlag, 0666)
				require.NoError(t, err)
				t.Cleanup(func() {
					_ = f.Close()
				})

				p, err := dockertestx.New(dockertestx.PoolOption{
					KeepContainer: c.keepContainer,
					StateStore:    f,
				})
				require.NoError(t, err)
				t.Cleanup(func() {
					_ = p.Purge()
				})

				for i := range cs {
					t.Run(fmt.Sprintf("i=%d name=%s tag=%s", i, cs[i].name, cs[i].tag), func(t *testing.T) {
						dsn, err := p.NewResource(cs[i].factory, dockertestx.ContainerOption{
							Name: cs[i].name,
							Tag:  cs[i].tag,
						})
						require.NoError(t, err)

						assert.NotEmpty(t, dsn)
						if len(cs[i].previousDSN) > 0 {
							assert.Equal(t, cs[i].previousDSN, dsn)
						}
						cs[i].previousDSN = dsn
					})
				}

				require.NoError(t, p.Save())
				require.NoError(t, f.Close())
				require.NoError(t, p.Purge())

				{
					finalState, err := os.OpenFile(statePath, os.O_RDONLY, 0644)
					require.NoError(t, err)
					t.Cleanup(func() {
						require.NoError(t, finalState.Close())
					})
					var s []interface{}
					require.NoError(t, json.NewDecoder(finalState).Decode(&s))
					assert.Len(t, s, c.expectStateLength)
				}
			})
		}
	})
}
