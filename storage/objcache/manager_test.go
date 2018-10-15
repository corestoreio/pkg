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

package objcache

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func withError() Option {
	return Option{
		fn: func(p *Manager) error {
			return errors.NotImplemented.Newf("What?")
		},
	}
}

func TestNewProcessor_NewError(t *testing.T) {
	p, err := NewManager(withError())
	assert.Nil(t, p)
	assert.True(t, errors.NotImplemented.Match(err), "Error: %s", err)
}
