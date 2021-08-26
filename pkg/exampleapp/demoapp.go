package exampleapp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jusongchen/REST-app/pkg/logging"
	restapp "github.com/jusongchen/REST-app/pkg/rest/app"
	"github.com/sethvargo/go-envconfig"
)

//specification has app config
type specification struct {
	RestConfig restapp.Config `json:"rest_config,omitempty"`

	LogFormat string `env:"LOG_FORMAT,default=text" json:"log_format,omitempty" `
	LogLevel  string `env:"LOG_LEVEL,default=INFO" json:"log_level,omitempty" `

	//if set to true, start service but do not mount any RESTFUL resources
	BootstrapMode bool `env:"BOOTSTRAP_MODE,default=false" json:"bootstrap_mode"`

	//RootCAPath location of RootCA cert
	RootCAPath     string `env:"ROOT_CA_PATH,default=/etc/identity/ca/cacerts.pem" json:"root_ca_path,omitempty"`
	ClientCertPath string `env:"CLIENT_CERT_PATH,default=/etc/identity/client/certificates/client.pem" json:"client_cert_path,omitempty"`
	ClientKeyPath  string `env:"CLIENT_KEY_PATH,default=/etc/identity/client/keys/client-key.pem" json:"client_key_path,omitempty"`
}

var _ fmt.Stringer = specification{}

//New init a new application
func New() (*restapp.Instance, error) {
	// log.Infof("Environment:%s", os.Environ())
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	return NewWith(ctx, envconfig.OsLookuper())
}

//NewWith init a new app using passed in env config lookuper
func NewWith(ctx context.Context, l envconfig.Lookuper) (*restapp.Instance, error) {
	var spec specification
	var err error

	if err := envconfig.ProcessWith(ctx, &spec, l); err != nil {
		return nil, err
	}

	logger := logging.DefaultLogger().Named("Demoapp")
	logger.Infof(`App Specification: %s`, spec)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd:%v", err)
	}

	about := struct {
		Release     string
		BuildTime   string
		Commit      string
		Cwd         string
		Spec        specification
		EnvironVars string
	}{
		Release,
		BuildTime,
		Commit,
		cwd,
		spec,
		strings.Join(os.Environ(), " "),
	}
	data, err := json.MarshalIndent(about, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error json marshal %+v: %v", about, err)

	}
	spec.RestConfig.About = string(data)

	u := UserResource{users: map[string]User{}}

	a, err := restapp.New(spec.RestConfig, info, u.WebService())
	if err != nil {
		logger.Errorf("app init:%v", err)
		return nil, err
	}

	return a, nil
}

func (a specification) String() string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}
