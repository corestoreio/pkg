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
	"github.com/coreos/etcd/clientv3"
	"github.com/corestoreio/pkg/config"
)

// WithEtcdV3 reads the ...
func WithEtcd(ec *clientv3.Client, o Options) config.Option {
	return func(s *config.Service) error {
		// TODO read all applicable keys from etcd and pass them to the service
		return nil
	}
}
