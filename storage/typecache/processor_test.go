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
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/storage/typecache/tcbigcache"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewProcessor_EncoderError(t *testing.T) {
	p, err := typecache.NewProcessor()
	if err != nil {
		t.Fatal(err)
	}

	ch := struct {
		ErrChan chan error
	}{
		ErrChan: make(chan error),
	}
	err = p.Set([]byte("key1"), ch)
	assert.True(t, errors.IsFatal(err), "Error: %s", err)
}

func TestNewProcessor_DecoderError(t *testing.T) {
	p, err := typecache.NewProcessor(tcbigcache.With())
	if err != nil {
		t.Fatal(err)
	}
	key := []byte("key1")
	val1 := struct {
		Val string
	}{
		Val: "Gopher",
	}
	assert.NoError(t, p.Set(key, val1))

	var val2 struct {
		Val2 string
	}
	err = p.Get(key, &val2)
	assert.True(t, errors.IsFatal(err), "Error: %s", err)
}

func TestNewProcessor_GetError(t *testing.T) {
	p, err := typecache.NewProcessor(tcbigcache.With())
	if err != nil {
		t.Fatal(err)
	}
	key := []byte("key1")
	var ch struct {
		ErrChan string
	}
	err = p.Get(key, ch)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func testCountry(t *testing.T, wg *sync.WaitGroup, p *typecache.Processor, key []byte) {
	defer wg.Done()

	var val = getTestCountry(t)
	if err := p.Set(key, val); err != nil {
		t.Fatal(err)
	}

	var newVal = new(Country)
	if err := p.Get(key, newVal); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, val, newVal)
}

func testStoreSlice(t *testing.T, wg *sync.WaitGroup, p *typecache.Processor, key []byte) {
	defer wg.Done()

	var val = getTestStores()
	if err := p.Set(key, val); err != nil {
		t.Fatal(err)
	}

	var newVal TableStoreSlice
	if err := p.Get(key, &newVal); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, val, newVal)
}

func init() {
	gob.Register(Country{})
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

func init() {
	gob.Register(TableStoreSlice{})
	gob.Register(TableStore{})
}

// TableStoreSlice represents a collection type for DB table store
// Generated via tableToStruct.
type TableStoreSlice []*TableStore

// TableStore represents a type for DB table store
// Generated via tableToStruct.
type TableStore struct {
	StoreID   int64          `db:"store_id" json:",omitempty"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment
	Code      dbr.NullString `db:"code" json:",omitempty"`       // code varchar(32) NULL UNI
	WebsiteID int64          `db:"website_id" json:",omitempty"` // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	GroupID   int64          `db:"group_id" json:",omitempty"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'
	Name      string         `db:"name" json:",omitempty"`       // name varchar(255) NOT NULL
	SortOrder int64          `db:"sort_order" json:",omitempty"` // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'
	IsActive  bool           `db:"is_active" json:",omitempty"`  // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'
}

func getTestStores() TableStoreSlice {
	return TableStoreSlice{
		&TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		&TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
		&TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		&TableStore{IsActive: false, StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
	}
}
