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

var _ config.Storager = (*storage.Multi)(nil)

func TestMakeMulti(t *testing.T) {
	scpID := scope.Store.Pack(44)
	const path = "aa/bb/cc"
	cmpValue := func(t *testing.T, s config.Storager, wantData []byte) {
		v, found, err := s.Value(scpID, path)
		require.NoError(t, err)
		assert.True(t, found)
		assert.Exactly(t, wantData, v)
	}

	inMem1 := config.NewInMemoryStore()
	inMem2 := config.NewInMemoryStore()

	m := storage.MakeMulti(inMem1, inMem2)
	testVal := []byte(`I'm your bro-grammer'`)

	t.Run("write,read to,from all", func(t *testing.T) {
		require.NoError(t, m.Set(scpID, path, testVal))

		cmpValue(t, inMem1, testVal)
		cmpValue(t, inMem2, testVal)
		cmpValue(t, m, testVal)

	})

	t.Run("write timeout", func(t *testing.T) {
		testVal := []byte(`A bro-grammer has a hammer`)
		m.Backends = append(m.Backends, sleepWriter{d: time.Millisecond * 100})
		m.ContextTimeout = time.Millisecond * 20

		err := m.Set(scpID, path, testVal)
		require.Error(t, err)
		assert.Exactly(t, context.DeadlineExceeded.Error(), err.Error())

		cmpValue(t, inMem1, testVal)
		cmpValue(t, inMem2, testVal)
		cmpValue(t, m, testVal)

	})

	t.Run("write error", func(t *testing.T) {
		testVal := []byte(`You are a bro-grammer'`)
		m.Backends = append(m.Backends, sleepWriter{setErr: errors.AlreadyInUse.Newf("resource in use")})

		err := m.Set(scpID, path, testVal)
		require.Error(t, err)
		assert.Exactly(t, "resource in use", err.Error())

		cmpValue(t, inMem1, testVal)
		cmpValue(t, inMem2, testVal)
		cmpValue(t, m, testVal)
	})

	t.Run("found nothing", func(t *testing.T) {
		v, found, err := m.Value(scope.Website.Pack(44), path)
		require.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, v)
	})

}

type sleepWriter struct {
	d      time.Duration
	setErr error
}

func (sw sleepWriter) Set(scp scope.TypeID, path string, value []byte) error {
	if sw.d > 0 {
		time.Sleep(sw.d)
	}
	return sw.setErr
}

func (sw sleepWriter) Value(scp scope.TypeID, path string) (v []byte, found bool, err error) {
	return
}
