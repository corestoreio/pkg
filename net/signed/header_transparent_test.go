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

package signed_test

import (
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/signed"
	"github.com/stretchr/testify/assert"
)

type cacherMock struct {
	SetFn func(hash []byte, ttl time.Duration)
	HasFn func(hash []byte) bool
}

func (cm cacherMock) Set(hash []byte, ttl time.Duration) { cm.SetFn(hash, ttl) }
func (cm cacherMock) Has(hash []byte) bool               { return cm.HasFn(hash) }

func TestMakeTransparent(t *testing.T) {

	cm := cacherMock{
		SetFn: func(hash []byte, ttl time.Duration) {},
		HasFn: func(hash []byte) bool {
			return false
		},
	}

	tp := signed.MakeTransparent(cm, time.Millisecond*500)
	assert.Empty(t, tp.HeaderKey())
}
