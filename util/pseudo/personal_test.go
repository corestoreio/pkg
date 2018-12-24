package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestPersonal(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))
		t.Logf("=== %q ===", lang)

		v := s.Gender()
		t.Logf("Gender %q %q", lang, v)
		assert.LenBetween(t, v, 1, 40, "Gender failed with lang %s", lang)

		v = s.GenderAbbrev()
		t.Logf("GenderAbbrev %q %q", lang, v)
		assert.LenBetween(t, v, 1, 40, "GenderAbbrev failed with lang %s", lang)

		v = s.Language()
		t.Logf("Language %q %q", lang, v)
		assert.LenBetween(t, v, 1, 40, "Language failed with lang %s", lang)

	}
}
