// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 6.8.8: DO NOT EDIT

package esapi

import (
	"context"
	"net/http"
	"strings"
)

func newXPackWatcherStopFunc(t Transport) XPackWatcherStop {
	return func(o ...func(*XPackWatcherStopRequest)) (*Response, error) {
		var r = XPackWatcherStopRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// XPackWatcherStop - http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-stop.html
//
type XPackWatcherStop func(o ...func(*XPackWatcherStopRequest)) (*Response, error)

// XPackWatcherStopRequest configures the X Pack Watcher Stop API request.
//
type XPackWatcherStopRequest struct {
	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r XPackWatcherStopRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "POST"

	path.Grow(len("/_xpack/watcher/_stop"))
	path.WriteString("/_xpack/watcher/_stop")

	params = make(map[string]string)

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f XPackWatcherStop) WithContext(v context.Context) func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f XPackWatcherStop) WithPretty() func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f XPackWatcherStop) WithHuman() func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f XPackWatcherStop) WithErrorTrace() func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f XPackWatcherStop) WithFilterPath(v ...string) func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f XPackWatcherStop) WithHeader(h map[string]string) func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f XPackWatcherStop) WithOpaqueID(s string) func(*XPackWatcherStopRequest) {
	return func(r *XPackWatcherStopRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
