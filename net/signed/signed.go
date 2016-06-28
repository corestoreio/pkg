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

package signed

import (
	"encoding/hex"
	"hash"
	"net/http"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/hashpool"
	"github.com/zenazn/goji/web/mutil"
)

// todo: refactor to use Service type and backendsigned package

// WithCompressor is a middleware applies the GZIP or deflate algorithm on
// the bytes writer. GZIP or deflate usage depends on the HTTP Accept
// Encoding header. Flush(), Hijack() and CloseNotify() interfaces will be
// preserved. No header set, no compression takes place. GZIP has priority
// before deflate.
func WithSignature(h func() hash.Hash) mw.Middleware {

	var hp = hashpool.New(h)
	var bp = bufferpool.New(h().Size())

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			buf := bp.Get()
			hsh := hp.Get()

			lw := mutil.WrapWriter(w)
			lw.Tee(hsh)

			h.ServeHTTP(lw, r)

			tmp := hsh.Sum(buf.Bytes())
			buf.Reset()
			_, _ = buf.Write(tmp)

			// does not work write into buffer and then hash it.
			// it will work with http.Trailer.Header
			WriteHTTPContentSignature(w, hex.EncodeToString, "", "", buf.Bytes())

			hp.Put(hsh)
			bp.Put(buf)
		})
	}
}
