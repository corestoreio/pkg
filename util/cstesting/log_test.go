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

package cstesting

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
)

var _ log.Logger = (*tLog)(nil)

type mockLog struct {
	*bytes.Buffer
}

func (tl mockLog) Log(args ...interface{}) {
	fmt.Fprint(tl, args...)
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	l := NewLogger(mockLog{buf})
	assert.True(t, l.IsInfo(), "IsInfo should be true")
	assert.True(t, l.IsDebug(), "IsDebug should be true")

	l = l.With(log.Int("newLogger", 123456))
	l.Info("Hello", log.String("key1", "val1"))
	l.Debug("Hallo", log.String("key2", "val2"))

	assert.Exactly(t, "[INFO] Hello newLogger: 123456 key1: \"val1\"\n[DEBUG] Hallo newLogger: 123456 key2: \"val2\"\n",
		buf.String())
}
