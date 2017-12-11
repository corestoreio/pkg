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

package dml

import (
	"database/sql/driver"
	"time"
)

type nowSentinel struct{}

var now = time.Now

// Now is a value that serializes to the current time. Should only be used to
// testing purposes. This exported variable can be removed at any time.
var Now = nowSentinel{}

// Value implements the driver.Valuer interface.
func (n nowSentinel) Value() (driver.Value, error) {
	fnow := n.UTC().Format(timeFormat)
	return fnow, nil
}

// String returns the time string in format time.RFC3339Nano.
func (n nowSentinel) String() string {
	return n.UTC().Format(timeFormat)
}

// UTC returns the UTC time
func (n nowSentinel) UTC() time.Time {
	return now().UTC()
}
