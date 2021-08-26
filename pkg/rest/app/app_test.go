package app

import (
	"net/http"
	"os"
	"testing"

	"github.com/jusongchen/REST-app/pkg/rest/swagger"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
)

func TestApp_Run(t *testing.T) {

	var info = swagger.ServerInfo{
		Title:       "Demo app",
		Description: ` Demo app `,
		Contact:     "Jusong Chen",
		Email:       "demo@example.com",
		APIVersion:  "1.0.0",
	}

	tests := []struct {
		name    string
		envs    map[string]string
		wantErr bool
	}{
		{
			name: "all_config_with_default_OK",
			envs: map[string]string{
				"LOG_FORMAT":      "text",
				"SWAGGER_UI_PATH": "./testdata/swaggerUI",
				"HOST":            "127.0.0.1",
				"PORT":            "0",
			},
			wantErr: false,
		},

		{
			name: "Missing_env_with_default_OK",
			envs: map[string]string{
				"SWAGGER_UI_PATH": "./testdata/swaggerUI",
				"PORT":            "0",
			},
			wantErr: false,
		},

		{
			name: "BAD_swagger_UI_path",
			envs: map[string]string{
				"SWAGGER_UI_PATH": "./testdata/swaggerUI_NOT_EXITST",
				"PORT":            "0",
			},
			wantErr: true,
		},

		{
			name: "BAD_HOST_IP",
			envs: map[string]string{
				"SWAGGER_UI_PATH": "./testdata/swaggerUI",
				"HOST":            "badIP27.0.0.1",
				"PORT":            "0",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}

			spec := Config{}
			err := envconfig.Process("", &spec)
			require.NoError(t, err)

			u := UserResource{map[string]User{}}

			a, err := New(spec, info, u.WebService())

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			a.Start()
			defer a.Close()
			baseURL := a.Svr.URL

			paths := []string{
				ReadyzPath,
				HealthzPath,
				HomePath,
				userResourceRootPath,
			}
			for _, p := range paths {

				resp, err := http.Get(baseURL + p)
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equal(t, 200, resp.StatusCode)
			}
		})
	}
}
