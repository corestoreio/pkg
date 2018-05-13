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

package cfgetcdv3

import (
	"bytes"
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// DefaultKeyPrefix defines the global key prefix, which can be overwritten.
const DefaultKeyPrefix = "csv3/"

type Options struct {
	RequestTimeout time.Duration
	// KeyPrefix defines a global key prefix used for all keys
	KeyPrefix string
}

// service implemented interface config.Storager.
type service struct {
	options Options
	client  clientv3.KV
}

// NewService creates a new storage client with either a concret or a mocked
// object of the etcd v3.
func NewService(c clientv3.KV, o Options) (config.Storager, error) {

	s := &service{
		options: o,
		client:  c,
	}
	if s.options.KeyPrefix == "" {
		s.options.KeyPrefix = DefaultKeyPrefix
	}

	return s, nil
}

func (s *service) toKey(p *config.Path) (string, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	buf.WriteString(s.options.KeyPrefix)
	err := p.AppendFQ(buf)
	return buf.String(), err
}

// Set puts a key to the etcd service.
func (s *service) Set(p *config.Path, value []byte) error {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key, err := s.toKey(p)
	if err != nil {
		return errors.Wrapf(err, "[cfgetcdv3] toKey with key %q", key)
	}

	if _, err = s.client.Put(ctx, key, string(value)); err != nil {
		return errors.Wrapf(err, "[cfgetcdv3] Put failed with key %q", key)
	}
	return nil
}

// Get returns a value from the etcd service.
func (s *service) Get(p *config.Path) (v []byte, found bool, err error) {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key, err := s.toKey(p)
	if err != nil {
		return nil, false, errors.Wrapf(err, "[cfgetcdv3] toKey with key %q", key)
	}
	resp, err := s.client.Get(ctx, key)
	if err != nil {
		switch err {
		case context.Canceled:
			return nil, false, nil
		case context.DeadlineExceeded:
			return nil, false, nil
		case rpctypes.ErrEmptyKey:
			return nil, false, nil
		}
		return nil, false, errors.Wrapf(err, "[cfgetcdv3] Client Get with key %q", key)
	}

	keyBytes := []byte(key)
	for _, ev := range resp.Kvs {
		if bytes.Equal(ev.Key, keyBytes) { // maybe not necessary.
			return ev.Value, true, nil
		}
	}

	return nil, false, nil
}

// FakeClient implementation for testing purposes.
type FakeClient struct {
	PutFn    func(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	PutError error
	GetError error
	GetFn    func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	GetKey   []byte
	GetValue []byte
}

func (cm FakeClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if cm.PutFn != nil {
		return cm.PutFn(ctx, key, val, opts...)
	}
	return nil, cm.PutError
}

func (cm FakeClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if cm.GetFn != nil {
		return cm.GetFn(ctx, key, opts...)
	}
	if cm.GetError != nil {
		return nil, cm.GetError
	}
	return &clientv3.GetResponse{
		Kvs: []*mvccpb.KeyValue{
			{
				Key:   cm.GetKey,
				Value: cm.GetValue,
			},
		},
	}, nil
}

func (cm FakeClient) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return nil, errors.NotImplemented.Newf("[cfgetcdv3] Delete Not Implemented")
}

func (cm FakeClient) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, errors.NotImplemented.Newf("[cfgetcdv3] Compact Not Implemented")
}

func (cm FakeClient) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, errors.NotImplemented.Newf("[cfgetcdv3] Do Not Implemented")
}

func (cm FakeClient) Txn(ctx context.Context) clientv3.Txn { return nil }
