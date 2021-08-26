package middleware

import (
	"github.com/emicklei/go-restful"
)

//EnableCORS adds a response header that enables CORS
func EnableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	resp.AddHeader("Access-Control-Allow-Origin", "*")
	chain.ProcessFilter(req, resp)
}
