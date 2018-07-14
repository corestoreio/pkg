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

package request

// crypto/rand => http://blog.sgmansfield.com/2016/06/managing-syscall-overhead-with-crypto-rand/

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/corestoreio/log"
	loghttp "github.com/corestoreio/log/http"
)

// HeaderIDKeyName defines the name of the header used to transmit the request
// ID.
const HeaderIDKeyName = "X-Request-Id"

// ID represents a middleware for request Id generation and adding the ID to the
// header.
type ID struct {
	// HeaderIDKeyName identifies the key name in the request header. Can be
	// empty and falls back to constant HeaderIDKeyName.
	HeaderIDKeyName string
	// Count defines the optional start value. If nil starts at zero. To access
	// securely "Count" you must use the atomic package.
	Count *uint64
	// NewIDFunc generates a new ID. Can be nil.
	NewIDFunc func(*http.Request) string
	Log       log.Logger
}

func (iw *ID) newID() func(*http.Request) string {

	// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L40-L52
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		if _, err := rand.Read(buf[:]); err != nil {
			panic(err) // todo remove panic without giving up error reporting
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	prefix := fmt.Sprintf("%s/%s-", hostname, b64[0:10])

	return func(*http.Request) string {
		return prefix + strconv.FormatUint(atomic.AddUint64(iw.Count, 1), 10)
	}
}

// With is a middleware that injects a request ID into the response header of
// each request. Retrieve it using:
// 		w.Header().Get(HeaderIDKeyName)
// If the incoming request has a HeaderIDKeyName header then that value is used
// otherwise a random value is generated. You can specify your own generator by
// providing the NewIDFunc in an option. No options uses the default request
// prefix generator.
// The returned function is compatible to type mw.Middleware.
func (iw *ID) With() func(h http.Handler) http.Handler {
	if iw.HeaderIDKeyName == "" {
		iw.HeaderIDKeyName = HeaderIDKeyName
	}
	if iw.NewIDFunc == nil {
		iw.NewIDFunc = iw.newID()
	}
	if iw.Count == nil {
		iw.Count = new(uint64)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(iw.HeaderIDKeyName)
			if id == "" {
				id = iw.NewIDFunc(r)
			}
			if iw.Log != nil && iw.Log.IsDebug() {
				iw.Log.Debug("request.ID.With", log.String("id", id), loghttp.Request("request", r))
			}
			w.Header().Set(iw.HeaderIDKeyName, id)
			h.ServeHTTP(w, r)
		})
	}
}
