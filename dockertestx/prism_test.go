package dockertestx_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewPrism(t *testing.T) {
	t.Parallel()

	p, err := dockertestx.New(dockertestx.PoolOption{})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, p.Purge())
	})

	endpoint, err := p.NewResource(&dockertestx.PrismFactory{
		SpecURI:         "testdata/oas.yml",
		HealthCheckPath: "/health",
	}, dockertestx.ContainerOption{})
	require.NoError(t, err)

	u, err := url.JoinPath(endpoint, "books")
	require.NoError(t, err)
	resp, err := http.Get(u) // #nosec G107
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
