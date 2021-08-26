package middleware

import (
	"testing"

	"github.com/emicklei/go-restful"
)

func TestRequestIDRest(t *testing.T) {
	type args struct {
		req   *restful.Request
		resp  *restful.Response
		chain *restful.FilterChain
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RequestIDRest(tt.args.req, tt.args.resp, tt.args.chain)
		})
	}
}
