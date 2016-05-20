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

package typecache

import (
	"encoding/gob"
	"io"
	"time"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/csfw/util/errors"
)

// WithBigCache allows to set custom configuration options to the bigcache
// instance. Bigcache has been selected as the default cache if you do
// not apply any cache option.
// Default option: shards 256, LifeWindow one hour, Verbose false
func WithBigCache(c ...bigcache.Config) Options {
	def := bigcache.Config{
		// optimize this ...
		Shards:             256,
		LifeWindow:         time.Hour,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   0,
	}
	if len(c) == 1 {
		def = c[0]
	}
	return func(p *Processor) {
		if p.optionError != nil {
			return
		}
		p.Cache, p.optionError = bigcache.NewBigCache(def)
		p.optionError = errors.NewFatal(p.optionError, "[typecache] bigcache.NewBigCache")
	}
}

// withGob defines the default encoder
func withGob() Options {
	return func(p *Processor) {
		p.enc = gob.NewEncoder(p.encBuf)
		p.dec = gob.NewDecoder(p.decBuf)
	}
}

// WithEncoder sets a custom encoder and decoder like message pack or protobuf.
func WithEncoder(enc func(io.Writer) Encoder, dec func(io.Reader) Decoder) Options {
	return func(p *Processor) {
		p.enc = enc(p.encBuf)
		p.dec = dec(p.decBuf)
	}
}

// WithCache sets a custom cache type for example Redis or MySQL.
func WithCache(c Cacher) Options {
	return func(p *Processor) {
		p.Cache = c
	}
}
