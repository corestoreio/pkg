package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestAddresses(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))
		t.Logf("=== %q ===", lang)
		v := s.Continent()
		t.Logf("Continent %q", v)
		assert.LenBetween(t, v, 1, 35, "Continent failed with lang %s", lang)

		v = s.Country()
		t.Logf("Country %q", v)
		assert.LenBetween(t, v, 1, 50, "Country failed with lang %s", lang)

		v = s.City()
		t.Logf("City %q", v)
		assert.LenBetween(t, v, 1, 40, "City failed with lang %s", lang)

		v = s.State()
		t.Logf("State %q", v)
		assert.LenBetween(t, v, 1, 56, "State failed with lang %s", lang)

		v = s.StateAbbrev()
		t.Logf("StateAbbrev %q", v)
		assert.LenBetween(t, v, 0, 35, "StateAbbrev failed with lang %s", lang)

		v = s.Street()
		t.Logf("Street %q", v)
		assert.LenBetween(t, v, 1, 40, "Street failed with lang %s", lang)

		v = s.StreetAddress()
		t.Logf("StreetAddress %q", v)
		assert.LenBetween(t, v, 1, 43, "StreetAddress failed with lang %s", lang)

		v = s.Zip()
		t.Logf("Zip %q", v)
		assert.LenBetween(t, v, 1, 35, "Zip failed with lang %s", lang)

		v = s.Phone()
		t.Logf("Phone %q", v)
		assert.LenBetween(t, v, 1, 35, "Phone failed with lang %s", lang)

	}
}
