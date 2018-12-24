package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestGeo(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		f := s.Latitude()
		if f < -90 || f > 90 {
			t.Errorf("Latitude failed with lang %s", lang)
		}

		i := s.LatitudeDegrees()
		if i < -180 || i > 180 {
			t.Errorf("LatitudeDegrees failed with lang %s", lang)
		}

		i = s.LatitudeMinutes()
		if i < 0 || i >= 60 {
			t.Errorf("LatitudeMinutes failed with lang %s", lang)
		}

		i = s.LatitudeSeconds()
		if i < 0 || i >= 60 {
			t.Errorf("LatitudeSeconds failed with lang %s", lang)
		}

		ld := s.LatitudeDirection()
		if ld != "N" && ld != "S" {
			t.Errorf("LatitudeDirection failed with lang %s", lang)
		}

		f = s.Longitude()
		if f < -180 || f > 180 {
			t.Errorf("Longitude failed with lang %s", lang)
		}

		i = s.LongitudeDegrees()
		if i < -180 || i > 180 {
			t.Errorf("LongitudeDegrees failed with lang %s", lang)
		}

		i = s.LongitudeMinutes()
		if i < 0 || i >= 60 {
			t.Errorf("LongitudeMinutes failed with lang %s", lang)
		}

		i = s.LongitudeSeconds()
		if i < 0 || i >= 60 {
			t.Errorf("LongitudeSeconds failed with lang %s", lang)
		}

		ld = s.LongitudeDirection()
		if ld != "W" && ld != "E" {
			t.Errorf("LongitudeDirection failed with lang %s", lang)
		}
	}
}
