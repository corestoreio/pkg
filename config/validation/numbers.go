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

//go:generate easyjson $GOFILE

package validation

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/byteconv"
	"github.com/corestoreio/pkg/util/validation"
)

// MinMaxInt64 validates if a value is between or in range of min and max.
//easyjson:json
type MinMaxInt64 struct {
	Min int64 `json:"min,omitempty"`
	Max int64 `json:"max,omitempty"`
}

func (v MinMaxInt64) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	val, ok, err := byteconv.ParseInt(rawData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return rawData, nil
	}
	if !validation.InRangeInt64(val, v.Min, v.Max) {
		return nil, errors.OutOfRange.Newf("[config/validation] %q value out of range: %d < v:%d < %d", v.Min, val, v.Max)
	}
	return rawData, nil
}
