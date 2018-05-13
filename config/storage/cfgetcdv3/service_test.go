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

package cfgetcdv3_test

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage/cfgetcdv3"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dialTimeout = 5 * time.Second
	// requestTimeout = 10 * time.Second
	endpoints = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
)

func TestStorage_Get(t *testing.T) {
	var testData = []byte(`You should turn it to eleven.`)
	scpID := scope.Website.WithID(3)
	const path = "path/to/orion"
	p := config.MustNewPathWithScope(scope.Website.WithID(3), "path/to/orion")

	t.Run("Get found", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			GetKey:   []byte(cfgetcdv3.DefaultKeyPrefix + `websites/3/` + path),
			GetValue: testData,
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		require.NoError(t, err)

		haveData, found, err := s.Get(scpID, path)
		require.NoError(t, err)
		assert.True(t, found, "Value and path must be found")
		assert.Exactly(t, haveData, testData)
	})

	t.Run("Get not found", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			GetKey:   []byte(`websites/3/`),
			GetValue: testData,
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		require.NoError(t, err)

		haveData, found, err := s.Get(scpID, path)
		require.NoError(t, err)
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	testGetErrors := func(getErr error) func(*testing.T) {
		return func(t *testing.T) {

			mo := cfgetcdv3.FakeClient{
				GetError: getErr,
			}

			s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
			require.NoError(t, err)

			haveData, found, err := s.Get(scpID, path)
			require.NoError(t, err)
			assert.False(t, found, "Value and path must NOT be found")
			assert.Nil(t, haveData)
		}
	}

	t.Run("Get context.Canceled", testGetErrors(context.Canceled))
	t.Run("Get context.DeadlineExceeded", testGetErrors(context.DeadlineExceeded))
	t.Run("Get rpctypes.ErrEmptyKey", testGetErrors(rpctypes.ErrEmptyKey))
	t.Run("Get any other error", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			GetError: errors.ConnectionLost.Newf("Ups"),
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		require.NoError(t, err)

		haveData, found, err := s.Get(scpID, path)
		require.True(t, errors.ConnectionLost.Match(err), "Should have error kind connection lost")
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	t.Run("Set no error ", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		require.NoError(t, err)

		err = s.Set(scpID, path, testData)
		require.NoError(t, err)
	})
	t.Run("Set no error ", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			PutError: errors.ConnectionLost.Newf("Ups"),
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		require.NoError(t, err)

		err = s.Set(scpID, path, testData)
		require.True(t, errors.ConnectionLost.Match(err), "Should have error kind connection lost")
	})

}

func TestNewStorage_Integration(t *testing.T) {

	t.Skip("TODO integration tests")

	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, c)

}
