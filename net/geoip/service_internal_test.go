package geoip

import (
	"bytes"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ io.Closer = (*Service)(nil)

func deferClose(t *testing.T, c io.Closer) {
	assert.NoError(t, c.Close())
}

func mustGetTestService(opts ...Option) *Service {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	return MustNew(append(opts, WithGeoIP2File(maxMindDB))...)
}

func TestMustNew(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.IsNotFound(err), "Error: %s", err)
			} else {
				t.Fatal("Expecting an error")
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	s := MustNew(WithGeoIP2File("not found"))
	assert.Nil(t, s)
}

func TestNewServiceErrorWithoutOptions(t *testing.T) {
	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.geoIP)
}

func TestNewService_WithGeoIP2File_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := New(
		WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
		WithGeoIP2File(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")),
	)
	defer deferClose(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.geoIP)
	for i := 0; i < 3; i++ {
		if err := s.Options(WithGeoIP2File(filepath.Join("testdata", "GeoIP2-Country-Test.mmdb"))); err != nil {
			t.Fatal(err)
		}
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `geoip.WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewService_WithGeoIP2Webservice_Atomic(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s, err := New(
		WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
		WithGeoIP2Webservice(nil, "a", "b", 1),
	)
	defer deferClose(t, s)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.geoIP)
	for i := 0; i < 3; i++ {
		assert.NoError(t, s.Options(WithGeoIP2Webservice(nil, "d", "e", 1)))
	}
	assert.True(t, 3 == strings.Count(logBuf.String(), `WithGeoIP.geoIPDone done: 1`), logBuf.String())
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := New(WithGeoIP2File("Walhalla/GeoIP2-Country-Test.mmdb"))
	assert.Nil(t, s)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func TestNewServiceWithGeoIP2Reader(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s)
	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland

	assert.NoError(t, err)
	haveCty, err := s.geoIP.Country(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", haveCty.Country.IsoCode)
}
