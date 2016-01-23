// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"sync"

	"github.com/corestoreio/csfw/util"
)

// Storager is the underlying data storage for holding the keys and its values.
// Implementations can be spf13/viper or MySQL backed. Default Storager
// is a simple mutex protected map[string]interface{}.
// ProTip: If you use MySQL as Storager don't execute function
// ApplyCoreConfigData()
type Storager interface {
	Set(key string, value interface{}) error
	Get(key string) interface{}
	AllKeys() []string
}

var _ Storager = (*simpleStorage)(nil)

type simpleStorage struct {
	sync.Mutex
	data map[string]interface{}
}

func newSimpleStorage() *simpleStorage {
	return &simpleStorage{
		data: make(map[string]interface{}),
	}
}

func (sp *simpleStorage) Set(key string, value interface{}) error {
	sp.Lock()
	sp.data[key] = value
	sp.Unlock()
	return nil
}

func (sp *simpleStorage) Get(key string) interface{} {
	sp.Lock()
	defer sp.Unlock()
	if data, ok := sp.data[key]; ok {
		return data
	}
	return nil
}
func (sp *simpleStorage) AllKeys() []string {
	sp.Lock()
	defer sp.Unlock()

	var ret = make(util.StringSlice, len(sp.data))
	i := 0
	for k := range sp.data {
		ret[i] = k
		i++
	}
	return ret.Sort()
}
