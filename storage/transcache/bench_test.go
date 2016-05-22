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
	"encoding/json"
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

func benchmark_country_enc(b *testing.B, opts ...transcache.Option) {
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
			if err := p.Set(key, cntry); err != nil {
				b.Fatal(errors.PrintLoc(err))
			}
			var newCntry = new(Country)
			if err := p.Get(key, newCntry); err != nil {
				b.Fatal(errors.PrintLoc(err))
			}
			if newCntry.Country.IsoCode != wantCountryISO {
				b.Fatalf("Country ISO Code must be %q, Have %q", wantCountryISO, newCntry.Country.IsoCode)
			}
			i++
		}
	})
}

func benchmark_stores_enc(b *testing.B, opts ...transcache.Option) {
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

func Benchmark_BigCache_Country_Gob(b *testing.B) {
	benchmark_country_enc(b, tcbigcache.With())
}

func Benchmark_BigCache_Country_JSON(b *testing.B) {
	benchmark_country_enc(b, tcbigcache.With(), transcache.WithEncoder(newJSONEncoder, newJSONDecoder))
}

func Benchmark_BigCache_Country_UgorjiMsgPack(b *testing.B) {
	benchmark_country_enc(b, tcbigcache.With(), transcache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
}

func Benchmark_BigCache_Stores_Gob(b *testing.B) {
	benchmark_stores_enc(b, tcbigcache.With())
}

func Benchmark_BigCache_Stores_JSON(b *testing.B) {
	benchmark_stores_enc(b, tcbigcache.With(), transcache.WithEncoder(newJSONEncoder, newJSONDecoder))
}

func Benchmark_BigCache_Stores_UgorjiMsgPack(b *testing.B) {
	benchmark_stores_enc(b, tcbigcache.With(), transcache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
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

func Benchmark_BoltDB_Country_Gob(b *testing.B) {
	f := getTempFile(b)
	defer os.Remove(f)
	benchmark_country_enc(b, tcboltdb.WithFile(f, 0600))
}

func Benchmark_BoltDB_Stores_Gob(b *testing.B) {
	f := getTempFile(b)
	defer os.Remove(f)
	benchmark_stores_enc(b, tcboltdb.WithFile(f, 0600))
}

func Benchmark_Redis_Country_Gob(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	benchmark_country_enc(b, tcredis.WithURL(redConURL, nil))
}

func Benchmark_Redis_Stores_Gob(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	benchmark_stores_enc(b, tcredis.WithURL(redConURL, nil))
}

func Benchmark_Redis_Country_UgorjiMsgPack(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	benchmark_country_enc(b, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
}

func Benchmark_Redis_Stores_UgorjiMsgPack(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	benchmark_stores_enc(b, tcredis.WithURL(redConURL, nil), transcache.WithEncoder(newUgorjiMsgPackEncoder, newUgorjiMsgPackDecoder))
}

func newJSONEncoder(w io.Writer) transcache.Encoder { return json.NewEncoder(w) }
func newJSONDecoder(r io.Reader) transcache.Decoder { return json.NewDecoder(r) }

var ugmsgPackHandle codec.MsgpackHandle

func newUgorjiMsgPackDecoder(r io.Reader) transcache.Decoder {
	return codec.NewDecoder(r, &ugmsgPackHandle)
}

func newUgorjiMsgPackEncoder(w io.Writer) transcache.Encoder {
	return codec.NewEncoder(w, &ugmsgPackHandle)
}
