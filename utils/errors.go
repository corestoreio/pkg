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

package utils

import (
	"bytes"
	"fmt"
)

// Errors returns a string where each error has been separated by a line break.
func Errors(errs ...error) string {
	if len(errs) == 0 || (len(errs) == 1 && errs[0] == nil) {
		return ""
	}
	var buf bytes.Buffer
	le := len(errs) - 1
	for i, e := range errs {
		if _, err := buf.WriteString(e.Error()); err != nil {
			return fmt.Sprintf("buf.WriteString internal error (%s): %s", err, e)
		}
		if i < le {
			if _, err := buf.WriteString("\n"); err != nil {
				return fmt.Sprintf("buf.WriteString internal error: %s when writing line break", err)
			}

		}
	}
	return buf.String()
}
