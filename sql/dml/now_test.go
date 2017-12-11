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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	now = func() time.Time {
		return time.Date(2006, 1, 2, 15, 4, 5, 02, time.FixedZone("UTC-4", -4*60*60))
	}
}

func TestNowSentinel_String(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "2006-01-02 19:04:05", Now.String())
}
