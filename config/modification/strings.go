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

package modification

import (
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// alter vs modification vs change | Don't know - non native speaker and
// StackExchange isn't helpful.

// Operator defines the function signature for altering the data.
type Operator func(*config.Path, []byte) ([]byte, error)

type ops struct {
	sync.RWMutex
	pool map[string]Operator
}

var operatorRegistry = &ops{
	pool: map[string]Operator{
		"upper":         toUpper,
		"lower":         toLower,
		"trim":          trim,
		"title":         toTitle,
		"base64_encode": base64Encode,
		"base64_decode": base64Decode,
		"hex_encode":    hexEncode,
		"hex_decode":    hexDecode,
		"sha256":        hash256, // as an example, but must register
		"gzip":          dataGzip,
		"gunzip":        dataGunzip,
	},
}

// RegisterOperator adds a new operator to the global registry and might
// overwrite previously set entries.
func RegisterOperator(typeName string, h Operator) {
	operatorRegistry.Lock()
	defer operatorRegistry.Unlock()
	operatorRegistry.pool[typeName] = h
}

// Strings defines the modificators to use to alter a string received from the
// config.Service.
//easyjson:json
type Strings struct {
	// Modificators currently supported: upper, lower, trim, title,
	// base64_encode, base64_decode, sha256 (must one time be registered in
	// hashpool package), gzip, gunzip.
	Modificators []string `json:"modificators,omitempty"`
}

// NewStrings creates a new type specific modificator.
func NewStrings(data Strings) (config.Observer, error) {
	ia := &observeStrings{
		opType: append([]string{}, data.Modificators...), // copy data
		opFns:  make([]Operator, 0, len(data.Modificators)),
	}

	operatorRegistry.RLock()
	defer operatorRegistry.RUnlock()

	for _, mod := range data.Modificators {
		h, ok := operatorRegistry.pool[mod]
		if !ok || h == nil {
			return nil, errors.NotSupported.Newf("[config/validation] Modificator %q not yet supported.", mod)
		}
		ia.opFns = append(ia.opFns, h)
	}

	return ia, nil
}

// MustNewStrings same as NewStrings but panics on error.
func MustNewStrings(data Strings) config.Observer {
	o, err := NewStrings(data)
	if err != nil {
		panic(err)
	}
	return o
}

// observeStrings must be used to prevent race conditions during initialization.
// That is the reason we have a separate struct for JSON handling. Having two
// structs allows to refrain from using Locks.
type observeStrings struct {
	opType []string
	opFns  []Operator
}

// Observe validates the given rawData value. This functions runs in a hot path.
func (v *observeStrings) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	rawData2 = rawData
	p2 := &p
	for idx, valFn := range v.opFns {
		if rawData2, err = valFn(p2, rawData2); err != nil {
			return nil, errors.Interrupted.New(err, "[config/modification] Function %q interrupted", v.opType[idx])
		}
	}
	return rawData2, nil
}
