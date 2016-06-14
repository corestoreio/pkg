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
	"encoding/json"
	"encoding/xml"
	"io"
)

// Option provides convenience helper functions to apply various options while
// creating a new Processor type.
type Option func(*Processor) error

// XMLCodec is used to encode/decode XML
type XMLCodec struct{}

// NewEncoder returns a new xml encoder which writes to w
func (c XMLCodec) NewEncoder(w io.Writer) Encoder {
	return xml.NewEncoder(w)
}

// NewDecoder returns a new xml decoder which reads from r
func (c XMLCodec) NewDecoder(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

// JSONCodec is used to encode/decode JSON
type JSONCodec struct{}

// NewEncoder returns a new json encoder which writes to w
func (c JSONCodec) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

// NewDecoder returns a new json decoder which reads from r
func (c JSONCodec) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

// GobCodec is used to encode/decode using the Gob format. You must use
// gob.Register to add new types to a pooled gob encoder.
type GobCodec struct{}

// NewEncoder returns a new gob encoder which writes to w
func (c GobCodec) NewEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

// NewDecoder returns a new gob decoder which reads from r
func (c GobCodec) NewDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}

// WithEncoder sets a custom encoder and decoder.
func WithEncoder(codec Codecer) Option {
	return func(p *Processor) error {
		p.Codec = codec
		return nil
	}
}

// WithPooledEncoder creates new encoder/decoder with a sync.Pool to reuse the
// objects. Providing argument primeObjects causes the encoder/decode to prime
// the data which means that no type information will be stored in the cache.
// If you use gob you must use gob.Register() for your types.
func WithPooledEncoder(codec Codecer, primeObjects ...interface{}) Option {
	return func(p *Processor) error {
		p.Codec = NewPooledCodec(codec, primeObjects...)
		return nil
	}
}

// WithCache sets a custom cache type. Examples in the subpackages.
func WithCache(c Cacher) Option {
	return func(p *Processor) error {
		p.Cache = c
		return nil
	}
}
