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
	"context"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// WithLoadData reads the all keys and their values with the current or configured
// etcd key prefix and applies it to the config.service. This function option
// can be set when creating a new config.service or updating its internal DB.
func WithLoadData(c clientv3.KV, o Options) config.LoadDataOption {

	if o.KeyPrefix == "" {
		o.KeyPrefix = DefaultKeyPrefix
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
				return errors.Wrapf(err, "[config/storage/cfgetcdv3] With Path %q", p.String())
			}

			if err := s.Set(p, ev.Value); err != nil {
				return errors.Wrapf(err, "[config/storage/cfgetcdv3] With Path %q", p.String())
			}
			buf.Reset()
			p.Reset()
		}

		return nil
	}).WithUseStorageLevel(1)
}
