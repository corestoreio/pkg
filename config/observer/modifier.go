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

package observer

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/gzippool"
	"github.com/corestoreio/pkg/util/hashpool"
)

// alter vs modification vs change | Don't know - non native speaker and
// StackExchange isn't helpful.

// ModificateFn defines the function signature for altering the data.
type ModificateFn func(*config.Path, []byte) ([]byte, error)

type modReg struct {
	sync.RWMutex
	pool map[string]ModificateFn
}

var modifierRegistry = &modReg{
	pool: map[string]ModificateFn{
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

// RegisterModifier adds a new modification function to the global registry
// and might overwrite previously set entries. Access to the global registry can
// be achieved via function NewModifier.
func RegisterModifier(typeName string, h ModificateFn) {
	modifierRegistry.Lock()
	defer modifierRegistry.Unlock()
	modifierRegistry.pool[typeName] = h
}

// ModifierArg defines the modifiers to use to alter a string received from the
// config.Service.
//easyjson:json
type ModifierArg struct {
	// Funcs defines a list of function names. Currently supported: upper,
	// lower, trim, title, base64_encode, base64_decode, sha256 (must one time
	// be registered in hashpool package), gzip, gunzip. Additional all other
	// custom modifier functions registered via RegisterModifier are
	// supported.
	Funcs []string `json:"funcs,omitempty"`
}

// NewModifier creates a new type specific modifier.
func NewModifier(data ModifierArg) (config.Observer, error) {
	ia := &modifiers{
		opType: append([]string{}, data.Funcs...), // copy data
		opFns:  make([]ModificateFn, 0, len(data.Funcs)),
	}

	modifierRegistry.RLock()
	defer modifierRegistry.RUnlock()

	for _, mod := range data.Funcs {
		h, ok := modifierRegistry.pool[mod]
		if !ok || h == nil {
			return nil, errors.NotSupported.Newf("[config/observer] Modifier %q not yet supported.", mod)
		}
		ia.opFns = append(ia.opFns, h)
	}

	return ia, nil
}

// MustNewModifier same as NewModifier but panics on error.
func MustNewModifier(data ModifierArg) config.Observer {
	o, err := NewModifier(data)
	if err != nil {
		panic(err)
	}
	return o
}

// modifiers must be used to prevent race conditions during initialization.
// That is the reason we have a separate struct for JSON handling. Having two
// structs allows to refrain from using Locks.
type modifiers struct {
	opType []string
	opFns  []ModificateFn
}

// Observe validates the given rawData value. This functions runs in a hot path.
func (v *modifiers) Observe(p config.Path, rawData []byte, found bool) (rawData2 []byte, err error) {
	rawData2 = rawData
	p2 := &p
	for idx, valFn := range v.opFns {
		if rawData2, err = valFn(p2, rawData2); err != nil {
			return nil, errors.Interrupted.New(err, "[config/observer] Function %q interrupted", v.opType[idx])
		}
	}
	return rawData2, nil
}

// as long as we don't see a use case for those modifiers in other packages,
// they stay private. might be refactored later.

func trim(_ *config.Path, data []byte) ([]byte, error) {
	return bytes.TrimSpace(data), nil
}

func toUpper(_ *config.Path, data []byte) ([]byte, error) {
	return bytes.ToUpper(data), nil
}

func toLower(_ *config.Path, data []byte) ([]byte, error) {
	return bytes.ToLower(data), nil
}

func toTitle(_ *config.Path, data []byte) ([]byte, error) {
	return bytes.Title(data), nil
}

func base64Encode(_ *config.Path, src []byte) (dst []byte, _ error) {
	dst = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return
}

func base64Decode(_ *config.Path, src []byte) (dst []byte, _ error) {
	dst = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)
	return
}

func hexEncode(_ *config.Path, src []byte) (dst []byte, _ error) {
	dst = make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst, nil
}

func hexDecode(_ *config.Path, src []byte) (dst []byte, _ error) {
	dst = make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	return dst, nil
}

// hash256 prefix the fully qualified path to src and then hashes it. Higher
// security.
func hash256(p *config.Path, src []byte) ([]byte, error) {
	tnk, err := hashpool.FromRegistry("sha256")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := p.AppendFQ(buf); err != nil {
		return nil, errors.Wrapf(err, "[config/observer] SHA256 with path %q", p.String())
	}
	buf.Write(src)
	var dst [sha256.Size]byte
	return tnk.Sum(buf.Bytes(), dst[:0]), nil
}

func dataGzip(_ *config.Path, src []byte) (dst []byte, _ error) {
	var buf bytes.Buffer
	buf.Grow(len(src) * 9 / 10) // *0.9
	zw := gzippool.GetWriter(&buf)
	defer gzippool.PutWriter(zw)
	zw.Write(src)
	if err := zw.Close(); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func dataGunzip(_ *config.Path, src []byte) (dst []byte, _ error) {
	r := bufferpool.GetReader(src)
	zr := gzippool.GetReader(r)
	defer func() {
		bufferpool.PutReader(r)
		gzippool.PutReader(zr)
	}()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(zr); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := zr.Close(); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}
