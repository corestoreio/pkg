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

// +build bigcache gob redis csall

package objcache_test

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/strs"
	"golang.org/x/sync/errgroup"
)

var _ io.Closer = (*objcache.Service)(nil)

func TestNewProcessor_EncoderError(t *testing.T) {
	t.Parallel()
	p, err := objcache.NewService(nil, objcache.NewCacheSimpleInmemory, newSrvOpt(gobCodec{}))
	assert.NoError(t, err)

	ch := struct {
		ErrChan chan error
	}{
		ErrChan: make(chan error),
	}
	err = p.Put(context.TODO(), "key1", ch, 0)
	assert.EqualError(t, err, "[objcache] With key \"key1\" and dst type struct { ErrChan chan error }: gob: type struct { ErrChan chan error } has no exported fields", "Error: %s", err)
}

func debugSliceBytes(t testing.TB, data ...[]byte) {
	for i, d := range data {
		t.Logf("IDX:%02d: IsNil %t | %q", i, d == nil, string(d))
	}
}

type myString struct {
	data string
	err  error
}

func (ms *myString) Unmarshal(data []byte) error {
	ms.data = string(data)
	return ms.err
}
func (ms *myString) Marshal() ([]byte, error) {
	return []byte(ms.data), ms.err
}

func TestService_Encoding(t *testing.T) {
	p, err := objcache.NewService(nil, objcache.NewCacheSimpleInmemory, newSrvOpt(gobCodec{}))
	assert.NoError(t, err)
	defer assert.NoError(t, p.Close())

	t.Run("marshal error", func(t *testing.T) {
		dErr := &myString{err: errors.BadEncoding.Newf("Bad encoding")}
		err := p.Put(context.TODO(), "dErr", dErr, 0)
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})
	t.Run("unmarshal error", func(t *testing.T) {
		err := p.Put(context.TODO(), "dErr2", 1, 0)
		assert.NoError(t, err)
		dErr := &myString{err: errors.BadEncoding.Newf("Bad encoding")}
		err = p.Get(context.TODO(), "dErr2", dErr)
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("marshal success", func(t *testing.T) {
		d1 := &myString{data: "HelloWorld"}
		d2 := &myString{data: "HalloWelt"}

		err = p.PutMulti(context.TODO(), []string{"d1x", "d2x"}, []interface{}{d1, d2}, nil)
		assert.NoError(t, err)

		d1.data = ""
		d2.data = ""
		err = p.GetMulti(context.TODO(), []string{"d1x", "d2x"}, []interface{}{d1, d2})
		assert.NoError(t, err)

		assert.Exactly(t, "HelloWorld", d1.data)
		assert.Exactly(t, "HalloWelt", d2.data)
	})
}

const iterations = 30

func testCountry(p *objcache.Service, key string) func() error {
	var val = mustGetTestCountry()
	return func() error {
		if err := p.Put(context.TODO(), key, val, 0); err != nil {
			return errors.WithStack(err)
		}

		for i := 0; i < iterations; i++ {
			var newVal = new(Country)
			if err := p.Get(context.TODO(), key, newVal); err != nil {
				return errors.WithStack(err)
			}
			if !reflect.DeepEqual(val, newVal) {
				return errors.Mismatch.Newf("%#v\n!=\n%#v", val, newVal)
			}
		}

		if err := p.Put(context.TODO(), key, Country{}, 0); err != nil {
			return errors.WithStack(err)
		}

		for i := 0; i < iterations; i++ {
			if err := p.Put(context.TODO(), key, val, 0); err != nil {
				return errors.WithStack(err)
			}
		}
		var newVal = new(Country)
		if err := p.Get(context.TODO(), key, newVal); err != nil {
			return errors.WithStack(err)
		}
		if !reflect.DeepEqual(val, newVal) {
			return errors.Mismatch.Newf("%#v\n!=\n%#v", val, newVal)
		}
		return nil
	}
}

func testStoreSlice(p *objcache.Service, key string) func() error {
	return func() error {

		var val = getTestStores()
		if err := p.Put(context.TODO(), key, val, 0); err != nil {
			return errors.WithStack(err)
		}
		for i := 0; i < iterations; i++ {
			var newVal TableStoreSlice
			if err := p.Get(context.TODO(), key, &newVal); err != nil {
				return errors.WithStack(err)
			}
			if !reflect.DeepEqual(val, newVal) {
				return errors.Mismatch.Newf("%#v\n!=\n%#v", val, newVal)
			}
		}
		if err := p.Put(context.TODO(), key, TableStoreSlice{}, 0); err != nil {
			return errors.WithStack(err)
		}

		for i := 0; i < iterations; i++ {
			if err := p.Put(context.TODO(), key, val, 0); err != nil {
				return errors.WithStack(err)
			}
		}
		var newVal TableStoreSlice
		if err := p.Get(context.TODO(), key, &newVal); err != nil {
			return errors.WithStack(err)
		}

		if !reflect.DeepEqual(val, newVal) {
			return errors.Mismatch.Newf("%#v\n!=\n%#v", val, newVal)
		}

		return nil
	}
}

func init() {
	gob.Register(&Country{})
	gob.Register(&TableStoreSlice{})
}

type Country struct {
	// IP contains the request IP address even if we run behind a proxy
	IP   net.IP `json:"ip,omitempty"`
	City struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameID  uint              `json:"geoname_id,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"city,omitempty"`
	Continent struct {
		Code      string            `json:"code,omitempty"`
		GeoNameID uint              `json:"geoname_id,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
	} `json:"continent,omitempty"`
	Country struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameID  uint              `json:"geoname_id,omitempty"`
		IsoCode    string            `json:"iso_code,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"country,omitempty"`
	Location struct {
		AccuracyRadius    int     `json:"accuracy_radius,omitempty"`
		AverageIncome     int     `json:"average_income,omitempty"`
		Latitude          float64 `json:"latitude,omitempty"`
		Longitude         float64 `json:"longitude,omitempty"`
		MetroCode         int     `json:"metro_code,omitempty"`
		PopulationDensity int     `json:"population_density,omitempty"`
		TimeZone          string  `json:"time_zone,omitempty"`
	} `json:"location,omitempty"`
	Postal struct {
		Code       string `json:"code,omitempty"`
		Confidence int    `json:"confidence,omitempty"`
	} `json:"postal,omitempty"`
	RegisteredCountry struct {
		GeoNameID uint              `json:"geoname_id,omitempty"`
		IsoCode   string            `json:"iso_code,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
	} `json:"registered_country,omitempty"`
	RepresentedCountry struct {
		GeoNameID uint              `json:"geoname_id,omitempty"`
		IsoCode   string            `json:"iso_code,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
		Type      string            `json:"type,omitempty"`
	} `json:"represented_country,omitempty"`
	Subdivision []struct {
		Confidence int               `json:"confidence,omitempty"`
		GeoNameId  uint              `json:"geoname_id,omitempty"`
		IsoCode    string            `json:"iso_code,omitempty"`
		Names      map[string]string `json:"names,omitempty"`
	} `json:"subdivisions,omitempty"`
	Traits struct {
		AutonomousSystemNumber       int    `json:"autonomous_system_number,omitempty"`
		AutonomousSystemOrganization string `json:"autonomous_system_organization,omitempty"`
		Domain                       string `json:"domain,omitempty"`
		IsAnonymousProxy             bool   `json:"is_anonymous_proxy,omitempty"`
		IsSatelliteProvider          bool   `json:"is_satellite_provider,omitempty"`
		Isp                          string `json:"isp,omitempty"`
		IpAddress                    string `json:"ip_address,omitempty"`
		Organization                 string `json:"organization,omitempty"`
		UserType                     string `json:"user_type,omitempty"`
	} `json:"traits,omitempty"`
	MaxMind struct {
		QueriesRemaining int `json:"queries_remaining,omitempty"`
	} `json:"maxmind,omitempty"`
}

func mustGetTestCountry() *Country {
	td, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		panic(err)
	}
	c := new(Country)
	if err := json.Unmarshal(td, c); err != nil {
		panic(err)
	}
	return c
}

// TableStoreSlice represents a collection type for DB table store
// Generated via tableToStruct.
type TableStoreSlice []*TableStore

// TableStore represents a type for DB table store
// Generated via tableToStruct.
type TableStore struct {
	StoreID   int64  `db:"store_id" json:",omitempty"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment
	Code      string `db:"code" json:",omitempty"`       // code varchar(32) NULL UNI
	WebsiteID int64  `db:"website_id" json:",omitempty"` // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	GroupID   int64  `db:"group_id" json:",omitempty"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	Name      string `db:"name" json:",omitempty"`       // name varchar(255) NOT NULL
	SortOrder int64  `db:"sort_order" json:",omitempty"` // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'
	IsActive  bool   `db:"is_active" json:",omitempty"`  // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'
}

func getTestStores() TableStoreSlice {
	return TableStoreSlice{
		&TableStore{StoreID: 0, Code: "admin", WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		&TableStore{StoreID: 5, Code: "au", WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 1, Code: "de", WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 4, Code: "uk", WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 2, Code: "at", WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
		&TableStore{StoreID: 6, Code: "nz", WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		&TableStore{IsActive: false, StoreID: 3, Code: "ch", WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
	}
}

func newServiceComplexParallelTest(t *testing.T, level2 objcache.NewStorageFn, so *objcache.ServiceOptions) {
	if so == nil {
		so = newSrvOpt(gobCodec{}, Country{}, TableStoreSlice{})
	}
	p, err := objcache.NewService(objcache.NewBlackHoleClient(nil), level2, so)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, p.Close())
	}()

	// to detect race conditions run with -race
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(testCountry(p, "country_one"))
	eg.Go(testStoreSlice(p, "stores_one"))
	eg.Go(testCountry(p, "country_two"))
	eg.Go(testStoreSlice(p, "stores_two"))
	eg.Go(testStoreSlice(p, "stores_three"))
	eg.Go(testCountry(p, "country_three"))

	assert.NoError(t, eg.Wait())
}

func newTestServiceDelete(t *testing.T, level2 objcache.NewStorageFn) {

	p, err := objcache.NewService(objcache.NewBlackHoleClient(nil), level2, newSrvOpt(JSONCodec{}))
	assert.NoError(t, err)
	defer func() { assert.NoError(t, p.Close()) }()

	t.Run("single key", func(t *testing.T) {
		err = p.Put(context.TODO(), "bc_delete", 1970, 0)
		assert.NoError(t, err)

		var bcInt int
		err = p.Get(context.TODO(), "bc_delete", &bcInt)
		assert.NoError(t, err)
		assert.Exactly(t, 1970, bcInt)

		err = p.Delete(context.TODO(), "bc_delete")
		assert.NoError(t, err)

		bcInt = 0
		err = p.Get(context.TODO(), "bc_delete", &bcInt)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, 0, bcInt)
	})

	t.Run("multiple keys", func(t *testing.T) {

		bcInt1 := 1971
		bcInt2 := 1972
		keys := []string{"bc_delete1", "bc_delete2"}
		vals := []interface{}{&bcInt1, &bcInt2}
		err = p.PutMulti(context.TODO(), keys, vals, nil)
		assert.NoError(t, err)

		bcInt1 = 0
		bcInt2 = 0
		err = p.GetMulti(context.TODO(), keys, vals)
		assert.NoError(t, err, "\n%+v", err)
		assert.Exactly(t, 1971, bcInt1)
		assert.Exactly(t, 1972, bcInt2)

		err = p.Delete(context.TODO(), "bc_delete1", "bc_delete2")
		assert.NoError(t, err)

		bcInt1 = 0
		bcInt2 = 0
		err = p.GetMulti(context.TODO(), keys, vals)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, 0, bcInt1)
		assert.Exactly(t, 0, bcInt2)
	})
}

func testExpiration(t *testing.T, cb func(), level2 objcache.NewStorageFn, so *objcache.ServiceOptions) {
	p, err := objcache.NewService(nil, level2, so)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, p.Close())
	}()

	key := strs.RandAlnum(30)
	if err := p.Put(context.TODO(), key, math.Pi, time.Second); err != nil {
		t.Fatalf("Key %q Error: %s", key, err)
	}

	var newVal float64
	err = p.Get(context.TODO(), key, &newVal)
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, math.Pi, newVal)

	cb()

	newVal = 0
	err = p.Get(context.TODO(), key, &newVal)
	assert.NoError(t, err, "%+v", err)
	assert.Empty(t, newVal)
}
