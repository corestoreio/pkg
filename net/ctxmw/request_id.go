// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxmw

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"os"
	"strconv"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputils"
	"golang.org/x/net/context"
)

// reqID is a global Counter used to create new request ids. This ID is not unique
// across multiple micro services.
var reqID int64

// RequestPrefixGenerator defines the functions needed to generate a request
// prefix id.
type RequestPrefixGenerator interface {
	// Generate returns a unique string for the current (micro) service.
	Generate() string
}

// DefaultRequestPrefix default prefix generator. Creates a prefix once the middleware
// is set up.
type DefaultRequestPrefix struct{}

// Generate returns a unique string for the current (micro) service.
func (rp DefaultRequestPrefix) Generate() string {
	// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L40-L52
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	return fmt.Sprintf("%s/%s-", hostname, b64[0:10])
}

// WithRequestID is a middleware that injects a request ID into the response header
// of each request. Retrieve it using:
// 		w.Header().Get(httputils.RequestIDHeader)
// If the incoming request has a RequestIDHeader header then that value is used
// otherwise a random value is generated. You can specify your own generator by
// providing the RequestPrefixGenerator once or pass no argument to use the default request
// prefix generator.
func WithRequestID(rpg ...RequestPrefixGenerator) ctxhttp.Middleware {
	var pf RequestPrefixGenerator
	pf = DefaultRequestPrefix{}
	if len(rpg) == 1 && rpg[0] != nil {
		pf = rpg[0]
	}

	reqPrefix := pf.Generate() // init once for all requests

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			id := r.Header.Get(httputils.RequestIDHeader)
			if id == "" {
				id = reqPrefix + strconv.FormatInt(atomic.AddInt64(&reqID, 1), 10)
			}
			w.Header().Set(httputils.RequestIDHeader, id)
			return hf(ctx, w, r)
		}
	}
}
