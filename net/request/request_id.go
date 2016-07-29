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
)

var reqID = new(int64)

// requestIDService default prefix generator. Creates a prefix once the
// middleware is set up.
type RequestIDService struct {
	prefix string
}

// Init returns a unique prefix string for the current (micro) service. This
// id gets reset once you restart the service.
func (rp *RequestIDService) Init() {
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
	rp.prefix = fmt.Sprintf("%s/%s-", hostname, b64[0:10])
}

// NewID returns a new ID unique for the current compilation.
func (rp *RequestIDService) NewID(_ *http.Request) string {
	return rp.prefix + strconv.FormatInt(atomic.AddInt64(reqID, 1), 10)
}
