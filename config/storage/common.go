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

package storage

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
)

type cacheKey struct {
	scope.TypeID
	string // route
}

func makeCacheKey(s scope.TypeID, r string) cacheKey {
	return cacheKey{s, r}
}

// WithLoadStrings loads a balanced fully qualified path and its stringified
// value pair into the config.Service. It does not panic when the fqPathValue
// slice argument isn't balanced, but returns an error. This functional option
// allows to load immutable values.
func WithLoadStrings(fqPathValue ...string) config.LoadDataOption {
	return config.MakeLoadDataOption(func(s *config.Service) (err error) {
		if lpv := len(fqPathValue); lpv%2 != 0 {
			return errors.NotAcceptable.Newf("[config/storage] WithLoadStrings: the argument slice fqPathValue is not balanced: length %d", lpv)
		}

		p := new(config.Path)
		for i := 0; i < len(fqPathValue) && err == nil; i = i + 2 {
			if err = p.ParseFQ(fqPathValue[i]); err != nil {
				return errors.WithStack(err)
			}
			err = s.Set(p, []byte(fqPathValue[i+1]))
		}
		return
	}).WithUseStorageLevel(1)
}
