package exampleapp

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/require"
)

func TestApp_Run_LocalDevNoDBNoAuthN(t *testing.T) {

	t.Run("app_run_local_dev", func(t *testing.T) {

		ctx1, cancelFn1 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn1()

		l := envconfig.MapLookuper(LocalDevEnv(t))
		a, err := NewWith(ctx1, l)

		require.NoError(t, err)
		require.NotNil(t, a)

		a.Start()
		defer a.Close()
		baseURL := a.Svr.URL

		paths := []string{
			"/home",
			"/healthz",
			"/apidocs",
		}
		for _, p := range paths {

			resp, err := http.Get(baseURL + p)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, 200, resp.StatusCode)
		}
	})
}

//LocalDevEnv returns local dev run minimal env var setting
func LocalDevEnv(tb testing.TB) map[string]string {

	env := map[string]string{

		"ROOT_CA_PATH":     "testdata/etc/identity/ca/cacerts.pem",
		"CLIENT_CERT_PATH": "testdata/etc/identity/client/certificates/client.pem",
		"CLIENT_KEY_PATH":  "testdata/etc/identity/client/keys/client-key.pem",
		"LOG_FORMAT":       "json",
		"SWAGGER_UI_PATH":  "./testdata/swaggerUI",
		"PORT":             "0",
		"HOST":             "127.0.0.1",
	}

	return env
}
