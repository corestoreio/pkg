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

	"github.com/alecthomas/assert"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/config/storage/cfgetcdv3"
	"github.com/corestoreio/pkg/store/scope"
)

func TestWithLoadData_Success(t *testing.T) {

	fc := cfgetcdv3.FakeClient{
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
	inMem := storage.NewMap()
	cfgSrv, err := config.NewService(
		inMem, config.Options{},
		cfgetcdv3.WithLoadData(fc, cfgetcdv3.Options{}),
	)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	p := config.MustNewPathWithScope(scope.Website.WithID(2), "payment/datatr/sha1")

	assert.Exactly(t, `"fc9d6fd2d8db223be4a7484a8619f26b"`, cfgSrv.Get(p).String())
	assert.Exactly(t, `"46aaccbebf47d8f8fce8c02d621aa573"`, cfgSrv.Get(p.BindStore(1)).String())
	assert.Exactly(t, `"e30d8df9810bc36105c96ad3ae76ffd3"`, cfgSrv.Get(p.BindDefault()).String())

}
