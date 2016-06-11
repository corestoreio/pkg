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
	"reflect"

	"github.com/corestoreio/csfw/util/errors"
)

// Options allows to set custom cache storage and encoder and decoder
type Option func(*Processor) error

// WithGobPriming allows to prime the gob encoder and decoder with the expected
// types for the Set() and Get() functions. It panics if the Encoder and Decoder
// are nil (not set) or not a type of the encoding/gob package. This function is
// not the same as gob.Register(). You shall not use this function more than
// once during a Get/Set operation for one type. Gob can only handle efficiently
// for multiple Get calls just one type. If you do not use this function you but
// still apply WithGob(), without an object to prime, you will get the error
// "extra data in buffer" from gob. This error means that the type information
// has also been stored in the cache and gob tries to set the internal cache a
// 2nd time with the type information. See the example for more details.
func WithGobPriming(primeObject interface{}) Option {
	return func(p *Processor) error {

		to := reflect.TypeOf(primeObject)
		encVal := reflect.Zero(to)
		decVal := reflect.New(to)

		for i := 0; i < encodeShards; i++ {
			p.enc[i].Lock()
			p.dec[i].Lock()
			defer p.enc[i].Unlock()
			defer p.dec[i].Unlock()

			enc := p.enc[i].Encoder.(*gob.Encoder)
			dec := p.dec[i].Decoder.(*gob.Decoder)

			if err := enc.EncodeValue(encVal); err != nil {
				return errors.NewFatal(err, "[transcache] WithGobPriming failed to encode a prime object")
			}
			if _, err := io.Copy(p.dec[i].buf, p.enc[i].buf); err != nil {
				return errors.NewFatal(err, "[transcache] WithGobPriming failed to copy from encode buffer to decode buffer")
			}
			p.enc[i].buf.Reset()

			if err := dec.DecodeValue(decVal); err != nil {
				return errors.NewFatal(err, "[transcache] WithGobPriming failed to decode a previously encoded object")
			}
			// decode buffer is now empty

			// encoder and decoder are no primed with the types for later use in Set() and
			// Get()
		}
		return nil
	}
}

// WithGob uses encoding/gob and allows you to prime one type. This function
// should only be used once because it creates the encoder and decoder. For more
// details please read documentation at WithGobPriming().
func WithGob(primeObject ...interface{}) Option {
	return func(p *Processor) error {
		for i := 0; i < encodeShards; i++ {
			p.enc[i].Encoder = gob.NewEncoder(p.enc[i].buf)
			p.dec[i].Decoder = gob.NewDecoder(p.dec[i].buf)
		}
		if len(primeObject) == 1 {
			return WithGobPriming(primeObject[0])(p)
		}
		return nil
	}
}

// WithEncoder sets a custom encoder and decoder like message-pack or protobuf,
// captnproto, JSON, XML ...
func WithEncoder(enc func(io.Writer) Encoder, dec func(io.Reader) Decoder) Option {
	return func(p *Processor) error {
		for i := 0; i < encodeShards; i++ {
			p.enc[i].Encoder = enc(p.enc[i].buf)
			p.dec[i].Decoder = dec(p.dec[i].buf)
		}
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
