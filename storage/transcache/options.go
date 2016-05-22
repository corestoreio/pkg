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

package transcache

import (
	"encoding/gob"
	"io"
)

// Options allows to set custom cache storage and encoder and decoder
type Option func(*Processor) error

// withGob defines the default encoder
func withGob() Option {
	return func(p *Processor) error {
		for i := 0; i < encodeShards; i++ {
			p.enc[i].Encoder = gob.NewEncoder(p.enc[i].buf)
			p.dec[i].Decoder = gob.NewDecoder(p.dec[i].buf)
		}
		return nil
	}
}

// WithEncoder sets a custom encoder and decoder like message pack or protobuf.
func WithEncoder(enc func(io.Writer) Encoder, dec func(io.Reader) Decoder) Option {
	return func(p *Processor) error {
		for i := 0; i < encodeShards; i++ {
			p.enc[i].Encoder = enc(p.enc[i].buf)
			p.dec[i].Decoder = dec(p.dec[i].buf)
		}
		return nil
	}
}

// WithCache sets a custom cache type for example Redis or MySQL.
func WithCache(c Cacher) Option {
	return func(p *Processor) error {
		p.Cache = c
		return nil
	}
}
