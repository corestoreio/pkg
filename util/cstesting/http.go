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

package cstesting

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HttpTrip used for mocking the Transport field in http.Client.
type HttpTrip struct {
	Resp *http.Response
	Err  error
	Req  *http.Request
}

// NewHttpTrip creates a new trip.
func NewHttpTrip(code int, body string, err error) *HttpTrip {
	return &HttpTrip{
		Resp: &http.Response{
			StatusCode: code,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		},
		Err: err,
	}
}

// RoundTrip implements http.RoundTripper and adds the Request to the
// field Req for later inspection.
func (tp *HttpTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	tp.Req = r
	return tp.Resp, tp.Err
}
