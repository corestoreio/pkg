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

package http_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/log"
	loghttp "github.com/corestoreio/csfw/log/http"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/buffer"
)

const testKey = "MyTestKey"

func TestField_Request(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`

	req := httptest.NewRequest("GET", "https://corestore.io", strings.NewReader(data))
	req.Header.Set("X-CoreStore-ID", "349:44")

	f := loghttp.Request(testKey, req)

	buf := &bytes.Buffer{}
	wt := log.WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"GET https://corestore.io HTTP/1.1\\r\\nX-Corestore-Id: 349:44\\r\\n\\r\\n35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams\"", buf.String())
}

func TestField_RequestHeader(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`

	req := httptest.NewRequest("GET", "https://corestore.io", strings.NewReader(data))
	req.Header.Set("X-CoreStore-ID", "349:44")

	f := loghttp.RequestHeader(testKey, req)

	buf := &bytes.Buffer{}
	wt := log.WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"GET https://corestore.io HTTP/1.1\\r\\nX-Corestore-Id: 349:44\\r\\n\\r\\n\"", buf.String())
}

type errCloser struct {
	error
	io.Reader
}

func (ec errCloser) Close() error { return ec.error }

func closerWithErr(err error, r io.Reader) io.ReadCloser {
	return errCloser{error: err, Reader: r}
}

func TestField_Request_Error(t *testing.T) {

	testR := httptest.NewRequest("GET", "/", closerWithErr(errors.New("XErr"), buffer.NewReader(nil)))

	f := loghttp.Request(testKey, testR)

	buf := &bytes.Buffer{}
	wt := log.WriteTypes{W: buf}

	assert.EqualError(t, f.AddTo(wt), `[log] AddTo.StringFn: [log] AddTo.HTTPRequest.DumpRequest: XErr`)
}

func TestField_Response(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`

	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"X-CoreStore-ID": []string{"987654321"},
		},
		Body:          ioutil.NopCloser(strings.NewReader(data)),
		ContentLength: int64(len(data)),
	}

	f := loghttp.Response(testKey, res)

	buf := &bytes.Buffer{}
	wt := log.WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"HTTP/1.0 200 OK\\r\\nContent-Length: 85\\r\\nX-CoreStore-ID: 987654321\\r\\n\\r\\n35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams\"", buf.String())
}

func TestField_Response_Error(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`

	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"X-CoreStore-ID": []string{"987654321"},
		},
		Body:          closerWithErr(errors.New("XErr"), buffer.NewReader(nil)),
		ContentLength: int64(len(data)),
	}

	f := loghttp.Response(testKey, res)

	buf := &bytes.Buffer{}
	wt := log.WriteTypes{W: buf}
	assert.EqualError(t, f.AddTo(wt), `[log] AddTo.StringFn: [log] AddTo.HTTPRequest.DumpResponse: XErr`)
}
