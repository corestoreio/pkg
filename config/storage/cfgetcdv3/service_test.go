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
	"flag"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage/cfgetcdv3"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
)

func init() {
	flag.BoolVar(&runIntegration, "integration", false, "Enables dml integration tests")
}

var (
	runIntegration bool
	dialTimeout    = 1 * time.Second
	// requestTimeout = 2 * time.Second
	endpoints = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
)

func TestStorage_Get(t *testing.T) {
	var testData = []byte(`You should turn it to eleven.`)
	//scpID := scope.Website.WithID(3)
	const path = "path/to/orion"
	p := config.MustNewPathWithScope(scope.Website.WithID(3), "path/to/orion")

	t.Run("Get found", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			GetKey:   []byte(cfgetcdv3.DefaultKeyPrefix + `websites/3/` + path),
			GetValue: testData,
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.NoError(t, err)
		assert.True(t, found, "Value and path must be found")
		assert.Exactly(t, haveData, testData)
	})

	t.Run("Get not found", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			GetKey:   []byte(`websites/3/`),
			GetValue: testData,
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.NoError(t, err)
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	testGetErrors := func(getErr error) func(*testing.T) {
		return func(t *testing.T) {

			mo := cfgetcdv3.FakeClient{
				GetError: getErr,
			}

			s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
			assert.NoError(t, err)

			haveData, found, err := s.Get(p)
			assert.NoError(t, err)
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
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.True(t, errors.ConnectionLost.Match(err), "Should have error kind connection lost")
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	t.Run("Set no error ", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		assert.NoError(t, err)

		err = s.Set(p, testData)
		assert.NoError(t, err)
	})
	t.Run("Set no error ", func(t *testing.T) {

		mo := cfgetcdv3.FakeClient{
			PutError: errors.ConnectionLost.Newf("Ups"),
		}

		s, err := cfgetcdv3.NewService(mo, cfgetcdv3.Options{})
		assert.NoError(t, err)

		err = s.Set(p, testData)
		assert.True(t, errors.ConnectionLost.Match(err), "Should have error kind connection lost")
	})

}

func TestNewStorage_Integration(t *testing.T) {
	if !runIntegration {
		t.Skip("Skipped. To enable use -integration=1")
	}
	/*
	   $ etcdctl get --prefix csv3
	   csv3/default/0/tax/calculation/rate
	   19.0
	   csv3/stores/2/tax/calculation/rate
	   19.2
	*/
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		t.Skipf("etcd daemon seems not to be running: %s", err)
	}
	assert.NotNil(t, c)

	srv, err := cfgetcdv3.NewService(c, cfgetcdv3.Options{})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	p := config.MustNewPath("tax/calculation/rate")
	p2 := p.BindStore(2)

	assert.NoError(t, srv.Set(p, []byte(`19.0`)))
	assert.NoError(t, srv.Set(p2, []byte(`19.2`)))

	data, _, err := srv.Get(p)
	assert.NoError(t, err)
	assert.Exactly(t, []byte(`19.0`), data)

	data, _, err = srv.Get(p2)
	assert.NoError(t, err)
	assert.Exactly(t, []byte(`19.2`), data)

}
