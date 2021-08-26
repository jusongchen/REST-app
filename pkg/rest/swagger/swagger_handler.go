//Package swagger provides swagger UI handlers
package swagger

import (
	"fmt"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

const (
	swaggerUIHomeURL   = "/swagger-ui.html"
	swaggerUIAPIDocURL = "/apidocs/"
	apidocsJSONPath    = "/apidocs.json"
)

//ServerInfo server info for swagger
type ServerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
	Email       string `json:"email"`
	APIVersion  string `json:"APIVersion"`
}

//NewContainer returns a restful.Container with swagger handled
func NewContainer(webServicesURL, swaggerUIPath string, info ServerInfo, ws ...*restful.WebService) (*restful.Container, error) {

	if _, err := os.Stat(swaggerUIPath); os.IsNotExist(err) {

		return nil, fmt.Errorf("sawgger.NewContainer:swaggerUIPath does not exist:%s", swaggerUIPath)
	}
	c := restful.NewContainer()
	spec := buildConfig(c, webServicesURL, info, ws)
	// Swagger WebUI
	c.Add(restfulspec.NewOpenAPIService(*spec))
	c.Handle(swaggerUIHomeURL, handleSwaggerHomeUI())
	c.Handle(swaggerUIAPIDocURL, handleSwagger(swaggerUIPath))

	return c, nil
}

func handleSwagger(swaggerUIPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Infof("Request:from %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

		h := http.StripPrefix(swaggerUIAPIDocURL, http.FileServer(http.Dir(swaggerUIPath)))
		cors.Default().Handler(h).ServeHTTP(w, r)
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	}
}

func handleSwaggerHomeUI() http.HandlerFunc {
	return http.RedirectHandler(swaggerUIAPIDocURL, 302).ServeHTTP
}

func buildConfig(c *restful.Container, webServicesURL string, info ServerInfo, ws []*restful.WebService) *restfulspec.Config {

	for _, w := range ws {
		c.Add(w)
	}

	config := restfulspec.Config{
		WebServices:    c.RegisteredWebServices(),
		WebServicesURL: webServicesURL,
		APIPath:        apidocsJSONPath,
		PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
			swo.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       info.Title,
					Description: info.Description,
					Contact: &spec.ContactInfo{
						Name:  info.Contact,
						Email: info.Email,
					},
					Version: info.APIVersion,
				},
			}
		},
	}

	return &config
}
