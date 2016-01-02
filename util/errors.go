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

package util

import (
	"fmt"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errgo"
)

// Errors returns a string where each error has been separated by a line break.
// If an error implements errgo.Locationer the location will be added to the
// output to show you the file name and line number.
func Errors(errs ...error) string {
	if len(errs) == 0 || (len(errs) == 1 && errs[0] == nil) {
		return ""
	}
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	le := len(errs) - 1
	for i, e := range errs {
		if _, err := buf.WriteString(e.Error()); err != nil {
			return fmt.Sprintf("buf.WriteString (1) internal error (%s): %s\n%s", err, e, buf.String())
		}

		if lerr, ok := e.(errgo.Locationer); ok {
			if _, err := buf.WriteString("\n" + lerr.Location().String()); err != nil {
				return fmt.Sprintf("buf.WriteString (2) internal error (%s): %s\n%s", err, e, buf.String())
			}
		}

		if i < le {
			if _, err := buf.WriteString("\n"); err != nil {
				return fmt.Sprintf("buf.WriteString (3) internal error: %s\n%s", err, buf.String())
			}
		}
	}
	return buf.String()
}
