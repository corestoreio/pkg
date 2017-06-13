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

package dbr

import (
	"database/sql/driver"
	"time"
)

type nowSentinel struct{}

var now = time.Now

// Now is a value that serializes to the current time
var Now = nowSentinel{}

// Value implements a valuer for compatibility
func (n nowSentinel) Value() (driver.Value, error) {
	fnow := n.UTC().Format(timeFormat)
	return fnow, nil
}

// String returns the time string in format "2006-01-02 15:04:05"
func (n nowSentinel) String() string {
	return n.UTC().Format(timeFormat)
}

// UTC returns the UTC time
func (n nowSentinel) UTC() time.Time {
	return now().UTC()
}
