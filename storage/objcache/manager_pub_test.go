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
	"net"
	"sync"
	"testing"

	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
)

var _ io.Closer = (*objcache.Manager)(nil)

func TestNewProcessor_EncoderError(t *testing.T) {
	p, err := objcache.NewManager(objcache.WithPooledEncoder(gobCodec{}))
	if err != nil {
		t.Fatal(err)
	}

	ch := struct {
		ErrChan chan error
	}{
		ErrChan: make(chan error),
	}
	err = p.Set(context.TODO(), "key1", ch, nil)
	assert.EqualError(t, err, "[objcache] With key \"key1\": gob: type struct { ErrChan chan error } has no exported fields", "Error: %s", err)
}

const iterations = 30

func testCountry(t *testing.T, wg *sync.WaitGroup, p *objcache.Manager, key string) {
	defer wg.Done()

	var val = getTestCountry(t)

	if err := p.Set(context.TODO(), key, val, nil); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < iterations; i++ {
		var newVal = new(Country)
		if err := p.Get(context.TODO(), key, newVal, nil); err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, val, newVal)
	}

	if err := p.Set(context.TODO(), key, Country{}, nil); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < iterations; i++ {
		if err := p.Set(context.TODO(), key, val, nil); err != nil {
			t.Fatal(err)
		}
	}
	var newVal = new(Country)
	if err := p.Get(context.TODO(), key, newVal, nil); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, val, newVal)

}

func testStoreSlice(t *testing.T, wg *sync.WaitGroup, p *objcache.Manager, key string) {
	defer wg.Done()

	var val = getTestStores()
	if err := p.Set(context.TODO(), key, val, nil); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < iterations; i++ {
		var newVal TableStoreSlice
		if err := p.Get(context.TODO(), key, &newVal, nil); err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, val, newVal)
	}
	if err := p.Set(context.TODO(), key, TableStoreSlice{}, nil); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < iterations; i++ {
		if err := p.Set(context.TODO(), key, val, nil); err != nil {
			t.Fatal(err)
		}
	}
	var newVal TableStoreSlice
	if err := p.Get(context.TODO(), key, &newVal, nil); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, val, newVal)
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

func getTestCountry(t interface {
	Fatal(...interface{})
}) *Country {
	td, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	c := new(Country)
	if err := json.Unmarshal(td, c); err != nil {
		t.Fatal(err)
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

func newTestNewProcessor(t *testing.T, opts ...objcache.Option) {
	p, err := objcache.NewManager(append(opts, objcache.WithPooledEncoder(gobCodec{}, Country{}, TableStoreSlice{}))...)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var wg sync.WaitGroup

	// to detect race conditions run with -race
	wg.Add(1)
	go testCountry(t, &wg, p, "country_one")

	wg.Add(1)
	go testStoreSlice(t, &wg, p, "stores_one")

	wg.Add(1)
	go testCountry(t, &wg, p, "country_two")

	wg.Add(1)
	go testStoreSlice(t, &wg, p, "stores_two")

	wg.Add(1)
	go testStoreSlice(t, &wg, p, "stores_three")

	wg.Add(1)
	go testCountry(t, &wg, p, "country_three")

	wg.Wait()
}
