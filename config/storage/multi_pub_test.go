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
	"context"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeMulti(t *testing.T) {

	p := config.MustNewPathWithScope(scope.Store.WithID(44), "aa/bb/cc")

	cmpGet := func(t *testing.T, s config.Storager, wantData []byte) {
		v, found, err := s.Get(p)
		require.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, wantData, v)
	}

	testVal := []byte(`I'm your bro-grammer'`)

	t.Run("write,read to,from all", func(t *testing.T) {
		inMem1 := storage.NewMap()
		inMem2 := storage.NewMap()
		m := storage.MakeMulti(storage.MultiOptions{}, inMem1, inMem2)

		require.NoError(t, m.Set(p, testVal))

		cmpGet(t, inMem1, testVal)
		cmpGet(t, inMem2, testVal)
		cmpGet(t, m, testVal)

	})

	t.Run("write timeout", func(t *testing.T) {
		inMem1 := storage.NewMap()
		inMem2 := storage.NewMap()

		m := storage.MakeMulti(storage.MultiOptions{
			ContextTimeout: time.Millisecond * 20,
		}, inMem1, inMem2, sleepWriter{d: time.Millisecond * 100})

		testVal := []byte(`A bro-grammer has a hammer`)

		err := m.Set(p, testVal)
		require.Error(t, err)
		assert.Exactly(t, context.DeadlineExceeded.Error(), err.Error())

		cmpGet(t, inMem1, testVal)
		cmpGet(t, inMem2, testVal)
		cmpGet(t, m, testVal)

	})

	t.Run("write error (racy)", func(t *testing.T) {
		inMem1 := storage.NewMap()
		inMem2 := storage.NewMap()

		m2 := storage.MakeMulti(storage.MultiOptions{}, sleepWriter{setErr: errors.AlreadyInUse.Newf("resource in use")})
		m := storage.MakeMulti(storage.MultiOptions{}, inMem1, inMem2, m2)

		testVal := []byte(`You are a bro-grammer'`)

		err := m.Set(p, testVal)
		require.Error(t, err)
		assert.Exactly(t, "resource in use", err.Error())

		cmpGet(t, inMem1, testVal)
		cmpGet(t, inMem2, testVal)
		cmpGet(t, m, testVal)
	})

	t.Run("found nothing", func(t *testing.T) {
		m := storage.MakeMulti(storage.MultiOptions{})

		validateNotFoundGet(t, m, scope.Website.WithID(44), "aa/bb/cc")
	})
}

type sleepWriter struct {
	d      time.Duration
	setErr error
}

func (sw sleepWriter) Set(_ *config.Path, _ []byte) error {
	if sw.d > 0 {
		time.Sleep(sw.d)
	}
	return sw.setErr
}

func (sw sleepWriter) Get(_ *config.Path) (v []byte, found bool, err error) {
	return
}
