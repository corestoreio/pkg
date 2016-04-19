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

package blacklist

import (
	"time"

	"github.com/coocood/freecache"
)

// FreeCache high performance cache for concurrent/parallel use cases
// like in net/http needed.
type FreeCache struct {
	*freecache.Cache
	emptyVal []byte
}

// NewFreeCache creates a new cache instance with a minimum size to be
// set to 512KB.
// If the size is set relatively large, you should call `debug.SetGCPercent()`,
// set it to a much smaller value to limit the memory consumption and GC pause time.
func NewFreeCache(size int) *FreeCache {
	return &FreeCache{
		Cache:    freecache.NewCache(size),
		emptyVal: []byte(`1`),
	}
}

// Set adds a token to the blacklist and may perform a
// purge operation. If expires <=0 the cached item will not expire. Set should
// be called when you log out a user. Set must make sure to copy away the
// token bytes or hash them.
func (fc *FreeCache) Set(token []byte, expires time.Duration) error {
	return fc.Cache.Set(token, fc.emptyVal, int(expires.Seconds()))
}

// Has checks if a token has been stored in the blacklist and may
// delete the token if expiration time is up.
func (fc *FreeCache) Has(token []byte) bool {
	val, err := fc.Cache.Get(token)
	if err == freecache.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}
	return val != nil
}
