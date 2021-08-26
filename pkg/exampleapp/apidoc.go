package exampleapp

import (
	swagger "github.com/jusongchen/REST-app/pkg/rest/swagger"
)

var info = swagger.ServerInfo{
	Title: "REST-app with Swagger UI",
	Description: ` A starter app example
`,
	Contact:    "Jusong Chen",
	Email:      "demoapp@gmail.com",
	APIVersion: "1.0.0",
}
