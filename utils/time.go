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
	"fmt"
	"time"
)

// ParseTimeStrict parses a formatted string and returns the time value it
// represents. The output is identical to time.Parse except it returns an
// error for strings that don't format to the input value.
//
// An example where the output differs from time.Parse would be:
// parseTimeStrict("1/2/06", "11/31/15")
//
// - time.Parse returns "2015-12-01 00:00:00 +0000 UTC"
//
// - ParseTimeStrict returns an error
func ParseTimeStrict(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return t, err
	}
	if t.Format(layout) != value {
		return t, fmt.Errorf("invalid time: %q", value)
	}
	return t, nil
}
