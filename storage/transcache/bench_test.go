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

package transcache_test

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/corestoreio/csfw/storage/transcache/tcbigcache"
	"github.com/corestoreio/csfw/storage/transcache/tcboltdb"
	"github.com/corestoreio/csfw/storage/transcache/tcredis"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/ugorji/go/codec"
)

// removed "gopkg.in/vmihailenco/msgpack.v2" because not worth it

func benchmark_country_enc(iterationsSetGet int, opts ...transcache.Option) func(b *testing.B) {
	return func(b *testing.B) {
		p, err := transcache.NewProcessor(opts...)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if err := p.Cache.Close(); err != nil {
				b.Fatal(err)
			}
		}()
		cntry := getTestCountry(b) // type already gob.Registered ...
		const wantCountryISO = "US"
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int64
			for pb.Next() {
				key := strconv.AppendInt(nil, i, 10) // 1 alloc
				i++

				if err := p.Set(key, cntry); err != nil {
					b.Fatal(errors.PrintLoc(err))
				}
				// Double execution might detect storing of type information in streaming encoders
				for j := 0; j < iterationsSetGet; j++ {
					var newCntry = new(Country)
					if err := p.Get(key, newCntry); err != nil {
						b.Fatal(errors.PrintLoc(err))
					}
					if newCntry.Country.IsoCode != wantCountryISO {
						b.Fatalf("Country ISO Code must be %q, Have %q", wantCountryISO, newCntry.Country.IsoCode)
					}
				}
			}
		})
	}
}

func benchmark_stores_enc(iterationsSetGet int, opts ...transcache.Option) func(b *testing.B) {
	return func(b *testing.B) {
		p, err := transcache.NewProcessor(opts...)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if err := p.Cache.Close(); err != nil {
				b.Fatal(err)
			}
		}()
		ts := getTestStores() // type already gob.Registered ...
		const wantStoreCode = "nz"
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int64
			for pb.Next() {
				key := strconv.AppendInt(nil, i, 10) // 1 alloc
				i++

				if err := p.Set(key, ts); err != nil {
					b.Fatal(err)
				}

				// Double execution might detect storing of type information in streaming encoders
				for j := 0; j < iterationsSetGet; j++ {
					var newTS TableStoreSlice
					if err := p.Get(key, &newTS); err != nil {
						b.Fatal(err)
					}
					if have := newTS[5].Code.String; have != wantStoreCode {
						b.Fatalf("Store Code in slice position 5 must be %q, Have %q", wantStoreCode, have)
					}
				}
			}
		})
	}
}

func Benchmark_BigCache_Country(b *testing.B) {
	b.Run("Gob_1x", benchmark_country_enc(1, tcbigcache.With(), transcache.WithPooledEncoder(transcache.GobCodec{}, Country{})))
	b.Run("Gob_2x", benchmark_country_enc(2, tcbigcache.With(), transcache.WithPooledEncoder(transcache.GobCodec{}, Country{})))
	b.Run("JSON_1x", benchmark_country_enc(1, tcbigcache.With(), transcache.WithPooledEncoder(transcache.JSONCodec{})))
	b.Run("JSON_2x", benchmark_country_enc(2, tcbigcache.With(), transcache.WithPooledEncoder(transcache.JSONCodec{})))
	b.Run("MsgPack_1x", benchmark_country_enc(1, tcbigcache.With(), transcache.WithEncoder(newMsgPackCodec())))
	b.Run("MsgPack_2x", benchmark_country_enc(2, tcbigcache.With(), transcache.WithEncoder(newMsgPackCodec())))
}

func Benchmark_BigCache_Stores(b *testing.B) {
	b.Run("Gob_1x", benchmark_stores_enc(1, tcbigcache.With(), transcache.WithPooledEncoder(transcache.GobCodec{}, TableStoreSlice{})))
	b.Run("Gob_2x", benchmark_stores_enc(2, tcbigcache.With(), transcache.WithPooledEncoder(transcache.GobCodec{}, TableStoreSlice{})))
	b.Run("JSON_1x", benchmark_stores_enc(1, tcbigcache.With(), transcache.WithPooledEncoder(transcache.JSONCodec{})))
	b.Run("JSON_2x", benchmark_stores_enc(2, tcbigcache.With(), transcache.WithPooledEncoder(transcache.JSONCodec{})))
	b.Run("MsgPack_1x", benchmark_stores_enc(1, tcbigcache.With(), transcache.WithEncoder(newMsgPackCodec())))
	b.Run("MsgPack_2x", benchmark_stores_enc(2, tcbigcache.With(), transcache.WithEncoder(newMsgPackCodec())))
}

func getTempFile(t interface {
	Fatal(...interface{})
}) string {
	f, err := ioutil.TempFile("", "tcboltdb_")
	if err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func Benchmark_BoltDB_Gob(b *testing.B) {
	f := getTempFile(b)
	defer os.Remove(f)
	b.Run("Country_1x", benchmark_country_enc(1, transcache.WithPooledEncoder(transcache.GobCodec{}, Country{}), tcboltdb.WithFile(f, 0600)))
	b.Run("Stores_1x", benchmark_stores_enc(1, tcboltdb.WithFile(f, 0600), transcache.WithPooledEncoder(transcache.GobCodec{}, TableStoreSlice{})))
}

func Benchmark_Redis_Gob(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	b.Run("Country_1x", benchmark_country_enc(1, tcredis.WithURL(redConURL, nil), transcache.WithPooledEncoder(transcache.GobCodec{}, Country{})))
	b.Run("Country_2x", benchmark_country_enc(2, tcredis.WithURL(redConURL, nil), transcache.WithPooledEncoder(transcache.GobCodec{}, Country{})))
	b.Run("Stores_1x", benchmark_stores_enc(1, tcredis.WithURL(redConURL, nil), transcache.WithPooledEncoder(transcache.GobCodec{}, TableStoreSlice{})))
	b.Run("Stores_2x", benchmark_stores_enc(2, tcredis.WithURL(redConURL, nil), transcache.WithPooledEncoder(transcache.GobCodec{}, TableStoreSlice{})))
}

func Benchmark_Redis_MsgPack(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	b.Run("Country_1x", benchmark_country_enc(1, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newMsgPackCodec())))
	b.Run("Country_2x", benchmark_country_enc(2, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newMsgPackCodec())))
	b.Run("Stores_1x", benchmark_stores_enc(1, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newMsgPackCodec())))
	b.Run("Stores_2x", benchmark_stores_enc(2, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newMsgPackCodec())))
}

var ugmsgPackHandle codec.MsgpackHandle

// msgPackCodec cannot be pooled because then it uses too much allocs and slows down.
type msgPackCodec struct{}

func newMsgPackCodec() msgPackCodec {
	return msgPackCodec{}
}

// NewEncoder returns a new json encoder which writes to w
func (c msgPackCodec) NewEncoder(w io.Writer) transcache.Encoder {
	return codec.NewEncoder(w, &ugmsgPackHandle)
}

// NewDecoder returns a new json decoder which reads from r
func (c msgPackCodec) NewDecoder(r io.Reader) transcache.Decoder {
	return codec.NewDecoder(r, &ugmsgPackHandle)
}
