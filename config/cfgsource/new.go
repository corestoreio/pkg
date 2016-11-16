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

package cfgsource

import "github.com/corestoreio/csfw/util/errors"

// NewByString creates a new ValueLabelSlice (VLS) from key,value list.
// It panics when arguments are imbalanced. Example:
// 		mySlice := NewValueLabelSlice("http", "HTTP (unsecure)", "https", "HTTPS (TLS)")
// Error behaviour: NotValid.
func NewByString(vl ...string) (Slice, error) {
	if len(vl)%2 != 0 {
		return nil, errors.NewNotValidf("[source] Imbalanced Pairs: %v", vl)
	}
	vls := make(Slice, len(vl)/2)
	j := 0
	for i := 0; i < len(vl); i = i + 2 {
		vls[j] = Pair{
			String:  vl[i],
			NotNull: NotNullString,
			label:   vl[i+1],
		}
		j++
	}
	return vls, nil
}

// MustNewByString same as NewByString but panics on error.
func MustNewByString(vl ...string) Slice {
	vls, err := NewByString(vl...)
	if err != nil {
		panic(err)
	}
	return vls
}

// NewByStringValue all passed arguments are values.
func NewByStringValue(values ...string) Slice {
	vls := make(Slice, len(values))
	for i := 0; i < len(values); i++ {
		vls[i] = Pair{
			String:  values[i],
			NotNull: NotNullString,
		}
	}
	return vls
}

// Ints a slice only used as argument to NewByInt.
type Ints []struct {
	Value int
	Label string
}

// NewByInt creates a new slice with integer values
func NewByInt(vl Ints) Slice {
	vls := make(Slice, len(vl))
	for i := 0; i < len(vl); i++ {
		vls[i] = Pair{
			Int:     vl[i].Value,
			NotNull: NotNullInt,
			label:   vl[i].Label,
		}
	}
	return vls
}

// NewByIntValue all passed arguments are values.
func NewByIntValue(values ...int) Slice {
	vls := make(Slice, len(values))
	for i := 0; i < len(values); i++ {
		vls[i] = Pair{
			Int:     values[i],
			NotNull: NotNullInt,
		}
	}
	return vls
}

// F64s a slice only used as argument to NewByFloat64.
type F64s []struct {
	Value float64
	Label string
}

// NewByFloat64 creates a new slice with float64 values
func NewByFloat64(vl F64s) Slice {
	vls := make(Slice, len(vl))
	for i := 0; i < len(vl); i++ {
		vls[i] = Pair{
			Float64: vl[i].Value,
			NotNull: NotNullFloat64,
			label:   vl[i].Label,
		}
	}
	return vls
}

// Bools a slice only used as argument to NewByBool.
type Bools []struct {
	Value bool
	Label string
}

// NewByBool creates a new slice with bool values
func NewByBool(vl Bools) Slice {
	vls := make(Slice, len(vl))
	for i := 0; i < len(vl); i++ {
		vls[i] = Pair{
			Bool:    vl[i].Value,
			NotNull: NotNullBool,
			label:   vl[i].Label,
		}
	}
	return vls
}
