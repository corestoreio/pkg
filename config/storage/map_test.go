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

package storage_test

import (
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewMap_OneKey(t *testing.T) {
	t.Parallel()

	p := config.MustNewPathWithScope(scope.Store.WithID(55), "aa/bb/cc")
	sp := storage.NewMap()

	assert.NoError(t, sp.Set(p, []byte(`19.99`)))

	validateFoundGet(t, sp, scope.Store.WithID(55), "aa/bb/cc", "19.99")
	validateNotFoundGet(t, sp, scope.Store.WithID(55), "ff/gg/hh")
	validateNotFoundGet(t, sp, 0, "rr/ss/tt")

	fl, ok := sp.(flusher)
	if !ok {
		t.Fatalf("%#v must implement Flusher interface", sp)
	}

	assert.NoError(t, fl.Flush())

	validateNotFoundGet(t, sp, scope.Store.WithID(55), "aa/bb/cc")

}
