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

package responseproxy_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/net/responseproxy"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWrapTee(t *testing.T) {
	tw := responseproxy.WrapTee(httptest.NewRecorder())
	buf := new(bytes.Buffer)
	tw.Tee(buf)

	data := []byte(`“It can be very dangerous to see things from somebody else's point of view without the proper training.” - Douglas Adams`)
	n, err := tw.Write(data)
	assert.Exactly(t, len(data), n)
	assert.Exactly(t, len(data), tw.BytesWritten())
	assert.NoError(t, err)
	assert.Exactly(t, string(data), buf.String())
	assert.Contains(t, fmt.Sprintf("%#v", tw), `responseproxy.flushWriter`) // poor mans type testing ;-)
}
