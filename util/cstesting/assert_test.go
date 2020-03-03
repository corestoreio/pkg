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

package cstesting_test

import (
	"fmt"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/cstesting"
)

type mockErrorf struct {
	data string
}

func (m *mockErrorf) Errorf(format string, args ...interface{}) {
	m.data = fmt.Sprintf(format, args...)
}

func TestContainsCount(t *testing.T) {
	t.Parallel()
	me := &mockErrorf{}
	cstesting.ContainsCount(me, "Hello Gopher", "Rust", 1)
	assert.Exactly(t, "\"Hello Gopher\" should contain \"Rust\" times 1 Have: 0 Want: 1", me.data)
	me.data = ""
	cstesting.ContainsCount(me, "Hello Gopher", "Gopher", 1)
	assert.Empty(t, me.data)
}
