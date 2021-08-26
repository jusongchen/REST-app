package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/emicklei/go-restful"
)

const (
	//LengthClipResponseBody when logging, clip response body content if body size is bigger than this threshold
	LengthClipResponseBody = 1000
)

// Logging Filter
func Logging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {

	r := req.Request
	reqID := GetReqID(r.Context())

	now := time.Now()
	requestDump, err := httputil.DumpRequest(req.Request, true)
	// requestDump, err := httputil.DumpRequest(req.Request, log.IsLevelEnabled(log.DebugLevel))
	if err != nil {
		log.Errorf("fail to dump request:%v", err)
		requestDump = []byte{}
	}

	body := ""
	chunks := strings.Split(string(requestDump), "\r\n\r\n")
	head := chunks[0]
	if len(chunks) > 1 {
		body = chunks[1]
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	log.WithFields(
		log.Fields{
			"request_id":  reqID,
			"remote_addr": req.Request.RemoteAddr,
			"scheme":      scheme,
			"host":        r.Host,
			"uri":         r.RequestURI,
			"proto":       r.Proto,
			"head":        head,
			"body":        body,
		}).Info("Request")

	c := NewResponseCapture(resp.ResponseWriter)
	resp.ResponseWriter = c

	chain.ProcessFilter(req, resp)

	h := strings.Builder{}
	c.Header().Write(&h)

	b := string(c.Bytes())
	count := len(b)
	if count > LengthClipResponseBody {
		b = b[:LengthClipResponseBody]
		b += fmt.Sprintf("...\n[rest clipped, total %d bytes]", count)
	}

	duration := time.Now().Sub(now)
	log.WithFields(
		log.Fields{
			"request_id":  reqID,
			"status_code": c.StatusCode(),
			"duration":    duration,
			"headers":     h.String(),
			"body":        b,
		}).Info("Response")

}

//ResponseCapture capture http response
type ResponseCapture struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
	body        *bytes.Buffer
}

// NewResponseCapture init a ResponseCapure
func NewResponseCapture(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{
		ResponseWriter: w,
		wroteHeader:    false,
		body:           new(bytes.Buffer),
	}
}

//Header reads respnse Header
func (c ResponseCapture) Header() http.Header {
	return c.ResponseWriter.Header()
}

//Write writes response
func (c ResponseCapture) Write(data []byte) (int, error) {
	if !c.wroteHeader {
		c.WriteHeader(http.StatusOK)
	}
	c.body.Write(data)
	return c.ResponseWriter.Write(data)
}

//WriteHeader write http headers
func (c *ResponseCapture) WriteHeader(statusCode int) {
	c.status = statusCode
	c.wroteHeader = true
	c.ResponseWriter.WriteHeader(statusCode)
}

//Bytes returns response body
func (c ResponseCapture) Bytes() []byte {
	return c.body.Bytes()
}

//StatusCode return status code
func (c ResponseCapture) StatusCode() int {
	return c.status
}
