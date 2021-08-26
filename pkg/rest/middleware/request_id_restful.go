package middleware

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/emicklei/go-restful"
)

// RequestIDRest filter
func RequestIDRest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	ctx := req.Request.Context()
	requestID := req.Request.Header.Get(RequestIDHeader)
	if requestID == "" {
		myid := atomic.AddUint64(&reqid, 1)
		requestID = fmt.Sprintf("%s-%06d", prefix, myid)
	}
	ctx = context.WithValue(ctx, RequestIDKey, requestID)

	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)

}
