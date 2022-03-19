package internal

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/strs"
	"golang.org/x/sync/errgroup"
)

func TestExpiration(t *testing.T, cb func(), level2 objcache.NewStorageFn[string], so *objcache.ServiceOptions) {
	p, err := objcache.NewService[string](nil, level2, so)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		assert.NoError(t, p.Close())
	}()

	key := strs.RandAlnum(30)
	if err := p.Set(context.TODO(), key, math.Pi, time.Second); err != nil {
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

func NewSrvOpt(c objcache.Codecer, primeObjects ...any) *objcache.ServiceOptions {
	return &objcache.ServiceOptions{
		Codec:        c,
		PrimeObjects: primeObjects,
	}
}

var _ objcache.Codecer = (*JSONCodec)(nil)

type JSONCodec struct{}

func (c JSONCodec) NewEncoder(w io.Writer) objcache.Encoder {
	return json.NewEncoder(w)
}

func (c JSONCodec) NewDecoder(r io.Reader) objcache.Decoder {
	return json.NewDecoder(r)
}

var _ objcache.Codecer = GobCodec{}

type GobCodec struct{}

func (c GobCodec) NewEncoder(w io.Writer) objcache.Encoder {
	return gob.NewEncoder(w)
}

func (c GobCodec) NewDecoder(r io.Reader) objcache.Decoder {
	return gob.NewDecoder(r)
}

func LookupRedisEnv(t testing.TB) string {
	redConURL := os.Getenv("CS_REDIS_TEST")
	if redConURL == "" {
		t.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/?db=3"
		`)
	}
	return redConURL
}

func NewServiceComplexParallelTest(t *testing.T, level2 objcache.NewStorageFn[string], so *objcache.ServiceOptions) {
	if so == nil {
		so = NewSrvOpt(GobCodec{}, Country{}, TableStoreSlice{})
	}
	p, err := objcache.NewService(objcache.NewBlackHoleClient[string](nil), level2, so)
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

func NewTestServiceDelete(t *testing.T, level2 objcache.NewStorageFn[string]) {
	p, err := objcache.NewService(objcache.NewBlackHoleClient[string](nil), level2, NewSrvOpt(JSONCodec{}))
	assert.NoError(t, err)
	defer func() { assert.NoError(t, p.Close()) }()

	t.Run("single key", func(t *testing.T) {
		err = p.Set(context.TODO(), "bc_delete", 1970, 0)
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
		vals := []any{&bcInt1, &bcInt2}
		err = p.SetMulti(context.TODO(), keys, vals, nil)
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

const iterations = 30

func testCountry(p *objcache.Service[string], key string) func() error {
	val := mustGetTestCountry()
	return func() error {
		if err := p.Set(context.TODO(), key, val, 0); err != nil {
			return fmt.Errorf("errorID:1644526078036 error: %w", err)
		}

		for i := 0; i < iterations; i++ {
			newVal := new(Country)
			if err := p.Get(context.TODO(), key, newVal); err != nil {
				return fmt.Errorf("errorID:1644612622397 error: %w", err)
			}
			if want, have := val.IP.String(), newVal.IP.String(); want != have {
				return fmt.Errorf("errorID:1644612686938 %q != %q", want, have)
			}
			if !reflect.DeepEqual(val, newVal) {
				return fmt.Errorf("errorID:1644612699145 %#v\n!=\n%#v", val, newVal)
			}
		}

		if err := p.Set(context.TODO(), key, Country{}, 0); err != nil {
			return fmt.Errorf("errorID:1644612713924 error: %w", err)
		}

		for i := 0; i < iterations; i++ {
			if err := p.Set(context.TODO(), key, val, 0); err != nil {
				return fmt.Errorf("errorID:1644612719012 error: %w", err)
			}
		}
		newVal := new(Country)
		if err := p.Get(context.TODO(), key, newVal); err != nil {
			return fmt.Errorf("errorID:1644612724908 error: %w", err)
		}
		if !reflect.DeepEqual(val, newVal) {
			return fmt.Errorf("%#v\n!=\n%#v", val, newVal)
		}
		return nil
	}
}

func testStoreSlice(p *objcache.Service[string], key string) func() error {
	return func() error {
		val := getTestStores()
		if err := p.Set(context.TODO(), key, val, 0); err != nil {
			return fmt.Errorf("errorID: 1644612742610 error: %w", err)
		}
		for i := 0; i < iterations; i++ {
			var newVal TableStoreSlice
			if err := p.Get(context.TODO(), key, &newVal); err != nil {
				return fmt.Errorf("errorID: 1644612740544 error: %w", err)
			}
			if !reflect.DeepEqual(val, newVal) {
				return fmt.Errorf("%#v\n!=\n%#v", val, newVal)
			}
		}
		if err := p.Set(context.TODO(), key, TableStoreSlice{}, 0); err != nil {
			return fmt.Errorf("errorID: 1644612738151 error: %w", err)
		}

		for i := 0; i < iterations; i++ {
			if err := p.Set(context.TODO(), key, val, 0); err != nil {
				return fmt.Errorf("errorID: 1644612733176 error: %w", err)
			}
		}
		var newVal TableStoreSlice
		if err := p.Get(context.TODO(), key, &newVal); err != nil {
			return fmt.Errorf("errorID: 1644612730733 error: %w", err)
		}

		if !reflect.DeepEqual(val, newVal) {
			return fmt.Errorf("%#v\n!=\n%#v", val, newVal)
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
