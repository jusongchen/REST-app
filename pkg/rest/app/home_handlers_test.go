package app

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerHome(t *testing.T) {

	require.HTTPSuccess(t, home(Instance{}), "GET", HomePath, nil)

}

func TestReadyz(t *testing.T) {
	isReady := &atomic.Value{}
	isReady.Store(false)
	require.HTTPError(t, readyz(isReady), "GET", ReadyzPath, nil)

	isReady.Store(true)
	require.HTTPSuccess(t, readyz(isReady), "GET", ReadyzPath, nil)
}

func TestHealthz(t *testing.T) {

	require.HTTPBodyContains(t, healthz(), "GET", HealthzPath, nil, "alive")
}
