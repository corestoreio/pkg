// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package element

import (
	"strings"
)

// Separator used in the database table core_config_data and in config.Service
// to separate the path parts.
const Separator byte = '/'

const sSeparator = "/"

// JoinRoutes joins multiple strings with the Separator constant.
func JoinRoutes(paths ...string) string {
	var buf strings.Builder
	for i, p := range paths {
		if i > 0 {
			buf.WriteByte(Separator)
		}
		buf.WriteString(p)
	}
	return buf.String()
}

// Sectioner at the moment only for testing
type Sectioner interface {
	// Defaults generates the default configuration from all fields. Key is the
	// path and value the value.
	Defaults() (DefaultMap, error)
}

// DefaultMap contains the default aka global configuration of a package. string
// is the fully qualified configuration path of scope default.
type DefaultMap map[string]interface{}
