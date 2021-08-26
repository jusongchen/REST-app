// Package app provides a common way to launch a RESTFUL server with built-in swagger UI support.
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jusongchen/REST-app/pkg/rest/swagger"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

const (
	swaggerUIHomeURL   = "/swagger-ui.html"
	swaggerUIAPIDocURL = "/apidocs/"
	apidocsJSONPath    = "/apidocs.json"
	//HomePath for app version
	HomePath = "/home"
	//HealthzPath for k8s health probe
	HealthzPath = "/healthz"
	//ReadyzPath for k8s readiness probe
	ReadyzPath = "/readyz"
	// MetricsPath for exposing Prometheus metrics
	MetricsPath = "/metrics"
	// UIPath is the default path for application UI access
	UIPath = "/ui/"
)

//Config is used to keep common App config
type Config struct {
	SwaggerDir string `json:"swagger_dir,omitempty" required:"true" envconfig:"SWAGGER_UI_PATH" env:"SWAGGER_UI_PATH,required"`
	Port       uint   `json:"port,omitempty" default:"0" required:"true" envconfig:"PORT" env:"PORT,default=0"`
	Host       string `json:"host,omitempty" default:"0.0.0.0" required:"true" envconfig:"HOST" env:"HOST,default=0.0.0.0"`
	About      string `json:"about,omitempty" `
}

var _ fmt.Stringer = Config{}

func (s Config) String() string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

//Instance struct represents an REST Instance with Swagger UI enabled
type Instance struct {
	StartupTime time.Time `json:"startup_time,omitempty"`
	Config
	Container *restful.Container `json:"-"`
	Svr       *Server            `json:"-"`
	isReady   *atomic.Value
}

//New init a new application instance
func New(conf Config, info swagger.ServerInfo, ws ...*restful.WebService) (*Instance, error) {

	var err error

	a := Instance{Config: conf}

	a.StartupTime = time.Now()
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get CWD:%v", err)
	}

	address := net.ParseIP(a.Host)
	if address == nil {
		return nil, fmt.Errorf("cannot parse host IP:%s", a.Host)
	}
	if _, err := os.Stat(a.SwaggerDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("swaggerUI dir does not exist:%s.  working dir:%s", a.SwaggerDir, cwd)
	}

	addr := net.JoinHostPort(address.String(), strconv.FormatUint(uint64(a.Port), 10))

	svr := newUnstartedServer(addr, healthz())
	a.Svr = svr

	c, err := swagger.NewContainer(svr.URL, a.SwaggerDir, info, ws...)
	if err != nil {
		return nil, err
	}
	svr.Config.Handler = c
	a.Container = c

	c.Handle(HealthzPath, healthz())
	c.Handle(HomePath, home(a))

	a.isReady = &atomic.Value{}
	c.Handle(ReadyzPath, readyz(a.isReady))

	return &a, nil

}

//Start starts a server and return immediately
func (a *Instance) Start() {

	a.Svr.Start()
	a.isReady.Store(true)
	log.Infof("server %s is ready to serve", a.Svr.URL)

}

//Close ends a server execution
func (a *Instance) Close() {
	a.Svr.Close()
	log.Infof("server %s shut down. exit.", a.Svr.URL)
}

//Run starts a server and keep running until either it gets a SIGINTR or ctx is Done.
func (a *Instance) Run(ctx context.Context) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	a.Start()
	//TODO detect signals.
	select {
	case <-ctx.Done():
		log.Infof("app get done signal. shutting down http server ...")
	case <-interrupt:
		log.Infof("Got SIGINT or SIGTERM, shutting down http server ...")
	}
	log.Infof("server %s is shutting down ...", a.Svr.URL)
	a.Close()
}
