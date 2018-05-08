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
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
)

type Options struct {
	RequestTimeout time.Duration
}

// Storage implemented interface config.Storager.
type Storage struct {
	options Options
	client  clientv3.KV
}

// NewStorage creates a new storage client with either a concret or a mocked
// object of the etcd v3.
func NewStorage(c clientv3.KV, o Options) (*Storage, error) {

	s := &Storage{
		options: o,
		client:  c,
	}

	return s, nil
}

func toKey(scp scope.TypeID, route string) string {
	s, id := scp.Unpack()

	var buf strings.Builder
	buf.WriteString(s.StrType())
	buf.WriteByte('/')
	buf.WriteString(strconv.FormatInt(id, 10))
	buf.WriteByte('/')
	buf.WriteString(route)
	return buf.String()
}

// Set puts a key to the etcd service.
func (s *Storage) Set(scp scope.TypeID, route string, value []byte) error {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key := toKey(scp, route)
	_, err := s.client.Put(ctx, key, string(value))
	if err != nil {
		return errors.Wrapf(err, "[cfgetcdv3] Put failed with key %q", key)
	}
	return nil
}

// Get returns a value from the etcd service.
func (s *Storage) Get(scp scope.TypeID, route string) (v []byte, found bool, err error) {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key := toKey(scp, route)
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

// ClientMock implementation for testing purposes.
type ClientMock struct {
	PutFn    func(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	PutError error
	GetError error
	GetFn    func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	GetKey   []byte
	GetValue []byte
}

func (cm ClientMock) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if cm.PutFn != nil {
		return cm.PutFn(ctx, key, val, opts...)
	}
	return nil, cm.PutError
}

func (cm ClientMock) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
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

func (cm ClientMock) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return nil, errors.NotImplemented.Newf("[cfgetcdv3] Delete Not Implemented")
}

func (cm ClientMock) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, errors.NotImplemented.Newf("[cfgetcdv3] Compact Not Implemented")
}

func (cm ClientMock) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, errors.NotImplemented.Newf("[cfgetcdv3] Do Not Implemented")
}

func (cm ClientMock) Txn(ctx context.Context) clientv3.Txn { return nil }
