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

package httputil

import (
	"net"
	"net/http"
)

// GetRemoteAddr extracts the remote address from a request and takes
// care of different headers in which an IP address can be stored.
// Return value can be nil.
func GetRemoteAddr(req *http.Request) net.IP {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}
	host, _, err := net.SplitHostPort(addr)

	if err != nil {
		host = addr
	}
	return net.ParseIP(host)
}
