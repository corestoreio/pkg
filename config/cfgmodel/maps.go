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

package cfgmodel

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/errors"
)

// MapIntResolver maps an integer value into a string value. The passed
// scoped configuration allows detecting the current scope.
type MapIntResolver interface {
	IntToStr(config.Scoped, int) (string, error)
}

// MapIntStr backend model for mapping stored integer values into string values. For
// example the configuration stores a country ID but we want the country's name.
type MapIntStr struct {
	Int
	MapIntResolver
}

// NewMapIntStr creates a new integer to string mapper type. It will panic while
// calling later Get() when the MapIntResolver has not been set.
func NewMapIntStr(path string, opts ...Option) MapIntStr {
	return MapIntStr{
		Int: NewInt(path, opts...),
	}
}

// Get returns an encrypted value decrypted. Panics if Encryptor interface is
// nil.
func (p MapIntStr) Get(sg config.Scoped) (string, error) {
	i, err := p.Int.Get(sg)
	if err != nil {
		return "", errors.Wrap(err, "[cfgmodel] MapIntStr.Byte.Get")
	}
	s, err := p.IntToStr(sg, i)
	return s, errors.Wrap(err, "[cfgmodel] MapIntStr.Get.Decrypt")
}
