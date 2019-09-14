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

// +build csall etcdv3

package storage

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

func init() {
	flag.BoolVar(&runIntegration, "integration", false, "Enables etcdv3 integration tests")
}

var (
	runIntegration bool
	dialTimeout    = 1 * time.Second
	// requestTimeout = 2 * time.Second
	endpoints = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
)

func TestStorage_Get(t *testing.T) {
	testData := []byte(`You should turn it to eleven.`)
	// scpID := scope.Website.WithID(3)
	const path = "path/to/orion"
	p := config.MustNewPathWithScope(scope.Website.WithID(3), "path/to/orion")

	t.Run("Get found", func(t *testing.T) {
		mo := Etcdv3FakeClient{
			GetKey:   []byte(Etcdv3DefaultKeyPrefix + `websites/3/` + path),
			GetValue: testData,
		}

		s, err := NewEtcdv3Client(mo, Etcdv3Options{})
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.NoError(t, err)
		assert.True(t, found, "Value and path must be found")
		assert.Exactly(t, haveData, testData)
	})

	t.Run("Get not found", func(t *testing.T) {
		mo := Etcdv3FakeClient{
			GetKey:   []byte(`websites/3/`),
			GetValue: testData,
		}

		s, err := NewEtcdv3Client(mo, Etcdv3Options{})
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.NoError(t, err)
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	testGetErrors := func(getErr error) func(*testing.T) {
		return func(t *testing.T) {
			mo := Etcdv3FakeClient{
				GetError: getErr,
			}

			s, err := NewEtcdv3Client(mo, Etcdv3Options{})
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
		mo := Etcdv3FakeClient{
			GetError: errors.ConnectionLost.Newf("Ups"),
		}

		s, err := NewEtcdv3Client(mo, Etcdv3Options{})
		assert.NoError(t, err)

		haveData, found, err := s.Get(p)
		assert.True(t, errors.ConnectionLost.Match(err), "Should have error kind connection lost")
		assert.False(t, found, "Value and path must NOT be found")
		assert.Nil(t, haveData)
	})

	t.Run("Set no error ", func(t *testing.T) {
		mo := Etcdv3FakeClient{}

		s, err := NewEtcdv3Client(mo, Etcdv3Options{})
		assert.NoError(t, err)

		err = s.Set(p, testData)
		assert.NoError(t, err)
	})
	t.Run("Set no error ", func(t *testing.T) {
		mo := Etcdv3FakeClient{
			PutError: errors.ConnectionLost.Newf("Ups"),
		}

		s, err := NewEtcdv3Client(mo, Etcdv3Options{})
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

	srv, err := NewEtcdv3Client(c, Etcdv3Options{})
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

func TestWithLoadData_Success(t *testing.T) {
	fc := Etcdv3FakeClient{
		GetFn: func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
			return &clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{
						Key:   []byte(`websites/2/payment/datatr/sha1`),
						Value: []byte(`fc9d6fd2d8db223be4a7484a8619f26b`),
					},
					{
						Key:   []byte(`stores/1/payment/datatr/sha1`),
						Value: []byte(`46aaccbebf47d8f8fce8c02d621aa573`),
					},
					{
						Key:   []byte(`default/0/payment/datatr/sha1`),
						Value: []byte(`e30d8df9810bc36105c96ad3ae76ffd3`),
					},
				},
			}, nil
		},
	}
	inMem := NewMap()
	cfgSrv, err := config.NewService(
		inMem, config.Options{},
		WithLoadFromEtcdv3(fc, Etcdv3Options{}),
	)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	p := config.MustNewPathWithScope(scope.Website.WithID(2), "payment/datatr/sha1")

	assert.Exactly(t, `"fc9d6fd2d8db223be4a7484a8619f26b"`, cfgSrv.Get(p).String())
	assert.Exactly(t, `"46aaccbebf47d8f8fce8c02d621aa573"`, cfgSrv.Get(p.BindStore(1)).String())
	assert.Exactly(t, `"e30d8df9810bc36105c96ad3ae76ffd3"`, cfgSrv.Get(p.BindDefault()).String())
}
