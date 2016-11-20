// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loghttp

import (
	"net/http"
	"net/http/httputil"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
)

func addToHTTPRequest(key string, r *http.Request, dumpBody bool) func(log.AddStringFn) error {
	// copy the request and do not store it in the closure to avoid race
	// conditions (proof => see TestHTTPRequest_Race in sub package)
	r2 := new(http.Request)
	*r2 = *r
	return func(addString log.AddStringFn) error {
		b, err := httputil.DumpRequest(r2, dumpBody)
		if err != nil {
			return errors.Wrap(err, "[log] AddTo.HTTPRequest.DumpRequest")
		}
		addString(key, string(b))

		return nil
	}
}

// HTTPRequest transforms the request with the function httputil.DumpRequest(r,
// true) into a string. The body gets logged also. Not completely race condition
// free because it depends on the Body io.ReadCloser implementation.
//
// DumpRequest returns the given request in its HTTP/1.x wire representation. It
// should only be used by servers to debug client requests. The returned
// representation is an approximation only; some details of the initial request
// are lost while parsing it into an http.Request. In particular, the order and
// case of header field names are lost. The order of values in multi-valued
// headers is kept intact. HTTP/2 requests are dumped in HTTP/1.x form, not in
// their original binary representations.
func Request(key string, r *http.Request) log.Field {
	return log.StringFn(key, addToHTTPRequest(key, r, true))
}

// HTTPRequestHeader transforms the request with the function
// httputil.DumpRequest(r, false) into a string. The body gets not logged and
// hence it is race condition free.
func RequestHeader(key string, r *http.Request) log.Field {
	return log.StringFn(key, addToHTTPRequest(key, r, false))
}

// todo: add http.DumpRequestOut() with header+body and header only

// todo: add ResponseHeader()

// HTTPResponse transforms the response with the function
// httputil.DumpResponse(r, true) into a string. Same behaviour as
// HTTPRequest(). The body gets logged also. Not completely race condition free
// because it depends on the Body io.ReadCloser implementation.
func Response(key string, r *http.Response) log.Field {
	return log.StringFn(key, func(addString log.AddStringFn) error {
		b, err := httputil.DumpResponse(r, true)
		if err != nil {
			return errors.Wrap(err, "[log] AddTo.HTTPRequest.DumpResponse")
		}
		addString(key, string(b))

		return nil
	})
}
