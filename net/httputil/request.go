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

package httputil

import (
	"net"
	"net/http"
	"strings"
	"unicode"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/csnet"
)

// GetRealIP extracts the remote address from a request and takes
// care of different headers in which an IP address can be stored.
// Checks if the IP in one of the header fields lies in csnet.PrivateIPRanges.
// Return value can be nil.
func GetRealIP(r *http.Request) net.IP {
	// Courtesy https://husobee.github.io/golang/ip-address/2015/12/17/remote-ip-go.html
	for _, h := range [2]string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			// header can contain spaces too, strip those out.
			addr := filterIP(addresses[i])
			if addr == "" {
				continue
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}
			realIP := net.ParseIP(host)
			if !realIP.IsGlobalUnicast() || csnet.PrivateIPRanges.In(realIP) {
				// bad address, go to next
				continue
			}

			if realIP != nil {
				return realIP
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return net.ParseIP(filterIP(host))
}

func filterIP(ip string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	for _, r := range ip {
		switch {
		case unicode.IsDigit(r), unicode.IsLetter(r), unicode.IsPunct(r):
			_, _ = buf.WriteRune(r)
		}
	}
	return buf.String()
}
