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
	"net/http"

	"github.com/corestoreio/csfw/net"
	"github.com/corestoreio/csfw/util/bufferpool"
)

// WriteHTTPContentSignature writes the content signature header using an encoder, which can be hex or base64.
// 	Content-Signature: keyId="rsa-key-1",algorithm="rsa-sha256",signature="Hex|Base64(RSA-SHA256(signing string))"
// 	Content-Signature: keyId="hmac-key-1",algorithm="hmac-sha1",signature="Hex|Base64(HMAC-SHA1(signing string))"
func WriteHTTPContentSignature(w http.ResponseWriter, encoder func(src []byte) string, keyID, algorithm string, signature []byte) {
	buf := bufferpool.Get()
	buf.WriteString(`keyId="`)
	buf.WriteString(keyID)
	buf.WriteString(`",algorithm="`)
	buf.WriteString(algorithm)
	buf.WriteString(`",signature="`)
	buf.WriteString(encoder(signature))
	buf.WriteRune('"')
	w.Header().Set(net.ContentSignature, buf.String())
	bufferpool.Put(buf)
}
