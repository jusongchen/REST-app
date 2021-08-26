package swagger

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

type dummyServer struct{}

func (dummy dummyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dummy server, should not be called."))
	w.WriteHeader(http.StatusInternalServerError)
}

func TestHandlerSwaggerUIEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		urlPath        string
		wantStatusCode int
		wantErr        bool
		Parse2Swg      bool
		CheckContain   bool
	}{

		{name: "swaggerUI APIDocs.json endpoint",
			urlPath:        apidocsJSONPath,
			wantStatusCode: 201,
			wantErr:        false,
			Parse2Swg:      true,
		},
		{name: "swaggerUI APIDocURL endpoint",
			urlPath:        swaggerUIAPIDocURL,
			wantStatusCode: 201,
			wantErr:        false,
			CheckContain:   true,
		},

		{name: "swaggerUIHomeURL end point",
			urlPath:        swaggerUIHomeURL,
			wantStatusCode: 201,
			wantErr:        false,
			CheckContain:   true,
		},
	}

	ts := newSwaggerTestHTTPServer(t)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.urlPath)
			require.NoError(t, err)
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			if tt.Parse2Swg {
				swg := spec.Swagger{}
				err = swg.UnmarshalJSON(data)
				require.NoError(t, err)
			}

			if tt.CheckContain {

				require.Contains(t, string(data), "swagger-ui")
			}
		})
	}
}

func newSwaggerTestHTTPServer(t *testing.T) *httptest.Server {
	swaggerUIPath := "./testdata/swaggerUI"

	//we need a dummu server as we need to get the server URL for building Swagger Config
	dummy := dummyServer{}

	ts := httptest.NewUnstartedServer(dummy)

	var info = ServerInfo{
		Title:       "Demo app",
		Description: `test `,
		Contact:     "Jusong Chen",
		Email:       "wheresome@gmail.com",
		APIVersion:  "1.0.0",
	}

	u := UserResource{map[string]User{}}

	c, err := NewContainer(nil, ts.URL, swaggerUIPath, info, u.WebService())
	require.NoError(t, err)
	httpSrv := ts.Config
	httpSrv.Handler = c
	ts.Start()

	return ts
}
