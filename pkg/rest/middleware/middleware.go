// Package middleware provides middlewares compatible to github.com/emicklei/go-restful
//
// all middleware must have this signiture:
//  	func RequestIDRest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain)
package middleware
