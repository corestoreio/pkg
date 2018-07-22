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
// Provide the field MinMax as a balanced slice where value n defines min and
// n+1 the max value. For JSON handling, see sub-package `json`.
type MinMaxInt64 struct {
	MinMax []int64 `json:"min_max,omitempty"`
	// PartialValidation if true only one of min/max pairs must be valid.
	PartialValidation bool `json:"partial_validation,omitempty"`
}

// NewMinMaxInt64 creates a new observer to check if a value is contained
// between min and max values. Argument MinMax must be balanced slice.
func NewMinMaxInt64(MinMax ...int64) (*MinMaxInt64, error) {
	return &MinMaxInt64{
		MinMax: MinMax,
	}, nil
}

func (v MinMaxInt64) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	lmm := len(v.MinMax)
	if lmm%2 == 1 || lmm < 1 {
		return nil, errors.NotAcceptable.Newf("[config/validation] MinMaxInt64 does not contain a balanced slice. Len: %d", lmm)
	}

	val, ok, err := byteconv.ParseInt(rawData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !ok {
		return rawData, nil
	}
	var validations int
	for i := 0; i < lmm; i = i + 2 {
		if left, right := v.MinMax[i], v.MinMax[i+1]; validation.InRangeInt64(val, left, right) {
			validations++
			if v.PartialValidation {
				return rawData, nil
			}
		}
	}

	if !v.PartialValidation && validations == lmm/2 {
		return rawData, nil
	}
	return nil, errors.OutOfRange.Newf("[config/validation] %q value out of range: %v", val, v.MinMax)
}
