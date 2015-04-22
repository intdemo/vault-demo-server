package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Request is a raw request configuration structure used to initiate
// API requests to the Vault server.
type Request struct {
	Method   string
	URL      *url.URL
	Params   url.Values
	Obj      interface{}
	Body     io.Reader
	BodySize int64
}

// SetJSONBody is used to set a request body that is a JSON-encoded value.
func (r *Request) SetJSONBody(val interface{}) error {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	r.Obj = val
	r.Body = buf
	r.BodySize = int64(buf.Len())
	return nil
}

// ResetJSONBody is used to reset the body for a redirect
func (r *Request) ResetJSONBody() error {
	if r.Body == nil {
		return nil
	}
	return r.SetJSONBody(r.Obj)
}

// ToHTTP turns this request into a valid *http.Request for use with the
// net/http package.
func (r *Request) ToHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.URL.RawQuery = r.Params.Encode()

	// Create the HTTP request
	req, err := http.NewRequest(r.Method, r.URL.RequestURI(), r.Body)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = r.URL.Scheme
	req.URL.Host = r.URL.Host
	req.Host = r.URL.Host

	return req, nil
}
