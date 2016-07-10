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

package scopedservice

import (
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestOptionsError(t *testing.T) {
	opts := OptionsError(errors.NewAlreadyClosedf("Something has already been closed."))
	s, err := New(opts...)
	assert.Nil(t, s)
	assert.True(t, errors.IsAlreadyClosed(err), "%+v", err)
}

func TestOptionFactories(t *testing.T) {

	var off OptionFactoryFunc = func(config.ScopedGetter) []Option {
		return []Option{
			withValue(scope.Store, 1, "a value for the store 1 scope"),
			withValue(scope.Website, 2, "a value for the website 2 scope"),
		}
	}

	of := NewOptionFactories()
	of.Register("key", off)
	assert.Exactly(t, []string{"key"}, of.Names())

	off2, err := of.Lookup("key")
	assert.NoError(t, err)
	assert.Exactly(t, fmt.Sprintf("%#v", off), fmt.Sprintf("%#v", off2)) // yes weird but it does the job

	off3, err := of.Lookup("not found")
	assert.Nil(t, off3)
	assert.True(t, errors.IsNotFound(err), "%+v", err)
}
