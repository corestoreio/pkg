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

// +build etcdv3

// build tag above: for some reason it cannot be added build tag "csall" because
// when running in parent directory `config` the command:`$ go test -tags csall
// ./...` fails with error: `panic: http: multiple registrations for
// /debug/requests`. Somehow x/net/trace gets loaded twice, one time as vendored
// and the other time as x/net/trace ...

package storage

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/bufferpool"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// Etcdv3DefaultKeyPrefix defines the global key prefix, which can be overwritten.
const Etcdv3DefaultKeyPrefix = "csv3/"

type Etcdv3Options struct {
	RequestTimeout time.Duration
	// KeyPrefix defines a global key prefix used for all keys
	KeyPrefix string
}

// service implemented interface config.Storager.
type etcdv3Client struct {
	options Etcdv3Options
	client  clientv3.KV
}

// NewEtcdv3Client creates a new storage client with either a concret or a mocked
// object of the etcd v3.
func NewEtcdv3Client(c clientv3.KV, o Etcdv3Options) (config.Storager, error) {

	s := &etcdv3Client{
		options: o,
		client:  c,
	}
	if s.options.KeyPrefix == "" {
		s.options.KeyPrefix = Etcdv3DefaultKeyPrefix
	}

	return s, nil
}

func (s *etcdv3Client) toKey(p *config.Path) (string, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	buf.WriteString(s.options.KeyPrefix)
	err := p.AppendFQ(buf)
	return buf.String(), err
}

// Set puts a key to the etcd service.
func (s *etcdv3Client) Set(p *config.Path, value []byte) error {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key, err := s.toKey(p)
	if err != nil {
		return errors.Wrapf(err, "[storage/etcdv3] toKey with key %q", key)
	}

	if _, err = s.client.Put(ctx, key, string(value)); err != nil {
		return errors.Wrapf(err, "[storage/etcdv3] Put failed with key %q", key)
	}
	return nil
}

// Get returns a value from the etcd service.
func (s *etcdv3Client) Get(p *config.Path) (v []byte, found bool, err error) {
	ctx := context.Background()
	if s.options.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.RequestTimeout)
		defer cancel()
	}

	key, err := s.toKey(p)
	if err != nil {
		return nil, false, errors.Wrapf(err, "[storage/etcdv3] toKey with key %q", key)
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
		return nil, false, errors.Wrapf(err, "[storage/etcdv3] Client Get with key %q", key)
	}

	keyBytes := []byte(key)
	for _, ev := range resp.Kvs {
		if bytes.Equal(ev.Key, keyBytes) { // maybe not necessary.
			return ev.Value, true, nil
		}
	}

	return nil, false, nil
}

// Etcdv3FakeClient implementation for testing purposes.
type Etcdv3FakeClient struct {
	PutFn    func(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	PutError error
	GetError error
	GetFn    func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	GetKey   []byte
	GetValue []byte
}

func (cm Etcdv3FakeClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if cm.PutFn != nil {
		return cm.PutFn(ctx, key, val, opts...)
	}
	return nil, cm.PutError
}

func (cm Etcdv3FakeClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
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

func (cm Etcdv3FakeClient) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return nil, errors.NotImplemented.Newf("[storage/etcdv3] Delete Not Implemented")
}

func (cm Etcdv3FakeClient) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, errors.NotImplemented.Newf("[storage/etcdv3] Compact Not Implemented")
}

func (cm Etcdv3FakeClient) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, errors.NotImplemented.Newf("[storage/etcdv3] Do Not Implemented")
}

func (cm Etcdv3FakeClient) Txn(ctx context.Context) clientv3.Txn { return nil }

// WithLoadFromEtcdv3 reads the all keys and their values with the current or configured
// etcd key prefix and applies it to the config.service. This function option
// can be set when creating a new config.service or updating its internal DB.
func WithLoadFromEtcdv3(c clientv3.KV, o Etcdv3Options) config.LoadDataOption {

	if o.KeyPrefix == "" {
		o.KeyPrefix = Etcdv3DefaultKeyPrefix
	}

	return config.MakeLoadDataOption(func(s *config.Service) error {

		ctx := context.Background()
		if o.RequestTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), o.RequestTimeout)
			defer cancel()
		}

		resp, err := c.Get(ctx, o.KeyPrefix, clientv3.WithPrefix())
		if err != nil {
			return errors.WithStack(err)
		}
		p := new(config.Path)
		var buf strings.Builder
		for _, ev := range resp.Kvs {
			buf.Write(ev.Key)

			if err := p.Parse(buf.String()); err != nil {
				return errors.Wrapf(err, "[storage/etcdv3] With Path %q", p.String())
			}

			if err := s.Set(p, ev.Value); err != nil {
				return errors.Wrapf(err, "[storage/etcdv3] With Path %q", p.String())
			}
			buf.Reset()
			p.Reset()
		}

		return nil
	}).WithUseStorageLevel(1)
}
