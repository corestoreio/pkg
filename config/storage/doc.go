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

// Package storage provides the available storage engines for level 1 and level
// 2 caches.
//
// Use Go build tags to enable special storage clients or file format loading
// functions. Supported tags are: bigcache (store in big cache), db (store in
// MySQL/MariaDB), etcdv3 (store in etcd cluster/server), load from json and
// yaml.
package storage
