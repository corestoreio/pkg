// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package net

import (
	"path"
	"strings"
)

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. `head` will never contain a slash and
// `tail` will always be a rooted path without trailing slash.
// Source: http://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.IndexByte(p[1:], '/') + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
