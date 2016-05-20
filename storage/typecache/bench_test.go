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

package typecache_test

import (
	"encoding/json"
	"io"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/ugorji/go/codec"
	vmihailencoMsgPack "gopkg.in/vmihailenco/msgpack.v2"
)

func benchmark_country_enc(b *testing.B, opts ...typecache.Options) {
	p, err := typecache.NewProcessor(opts...)
	if err != nil {
		b.Fatal(err)
	}
	cntry := getTestCountry(b) // type already gob.Registered ...
	const wantCountryISO = "US"
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := strconv.Itoa(i) // 1 alloc
			if err := p.Set(key, cntry); err != nil {
				b.Fatal(err)
			}
			var newCntry = new(Country)
			if err := p.Get(key, newCntry); err != nil {
				b.Fatal(err)
			}
			if newCntry.Country.IsoCode != wantCountryISO {
				b.Fatalf("Country ISO Code must be %q, Have %q", wantCountryISO, newCntry.Country.IsoCode)
			}
			i++
		}
	})
}

func benchmark_stores_enc(b *testing.B, opts ...typecache.Options) {
	p, err := typecache.NewProcessor(opts...)
	if err != nil {
		b.Fatal(err)
	}
	ts := getTestStores() // type already gob.Registered ...
	const wantStoreCode = "nz"
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := strconv.Itoa(i) // 1 alloc
			if err := p.Set(key, ts); err != nil {
				b.Fatal(err)
			}
			var newTS TableStoreSlice
			if err := p.Get(key, &newTS); err != nil {
				b.Fatal(err)
			}
			if have := newTS[5].Code.String; have != wantStoreCode {
				b.Fatalf("Store Code in slice position 5 must be %q, Have %q", wantStoreCode, have)
			}
			i++
		}
	})
}

func Benchmark_SetGet_Country_Gob(b *testing.B) {
	benchmark_country_enc(b)
}

func Benchmark_SetGet_Country_VmihailencoMsgPack(b *testing.B) {
	benchmark_country_enc(b, typecache.WithEncoder(newVmihailencoMsgPackEnc, newVmihailencoMsgPackDec))
}

func Benchmark_SetGet_Country_JSON(b *testing.B) {
	benchmark_country_enc(b, typecache.WithEncoder(newJSONEncoder, newJSONDecoder))
}

func Benchmark_SetGet_Country_UgorjiMsgPack(b *testing.B) {
	benchmark_country_enc(b, typecache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
}

func Benchmark_SetGet_Stores_Gob(b *testing.B) {
	benchmark_stores_enc(b)
}

func Benchmark_SetGet_Stores_VmihailencoMsgPack(b *testing.B) {
	benchmark_stores_enc(b, typecache.WithEncoder(newVmihailencoMsgPackEnc, newVmihailencoMsgPackDec))
}

func Benchmark_SetGet_Stores_JSON(b *testing.B) {
	benchmark_stores_enc(b, typecache.WithEncoder(newJSONEncoder, newJSONDecoder))
}

func Benchmark_SetGet_Stores_UgorjiMsgPack(b *testing.B) {
	benchmark_stores_enc(b, typecache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
}

func newJSONEncoder(w io.Writer) typecache.Encoder { return json.NewEncoder(w) }
func newJSONDecoder(r io.Reader) typecache.Decoder { return json.NewDecoder(r) }

type VmihailencoMsgPackEnc struct {
	enc *vmihailencoMsgPack.Encoder
	dec *vmihailencoMsgPack.Decoder
}

func newVmihailencoMsgPackEnc(w io.Writer) typecache.Encoder {
	return VmihailencoMsgPackEnc{
		enc: vmihailencoMsgPack.NewEncoder(w),
	}
}

func (m VmihailencoMsgPackEnc) Encode(src interface{}) error {
	return m.enc.Encode(src)
}

func (m VmihailencoMsgPackEnc) Decode(dst interface{}) error {
	return m.dec.Decode(dst)
}

func newVmihailencoMsgPackDec(r io.Reader) typecache.Decoder {
	return VmihailencoMsgPackEnc{
		dec: vmihailencoMsgPack.NewDecoder(r),
	}
}

var ugmsgPackHandle codec.MsgpackHandle

func newUgorjiMsgPackDecoder(r io.Reader) typecache.Decoder {
	return codec.NewDecoder(r, &ugmsgPackHandle)
}

func newUgorjiMsgPackEncoder(w io.Writer) typecache.Encoder {
	return codec.NewEncoder(w, &ugmsgPackHandle)
}
