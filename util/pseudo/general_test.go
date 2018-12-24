package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestGeneral(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		v := s.Password(4, 10, true, true, true)
		if v == "" {
			t.Errorf("Password failed with lang %s", lang)
		}

		v = s.SimplePassword()
		if v == "" {
			t.Errorf("SimplePassword failed with lang %s", lang)
		}

		v = s.Color()
		if v == "" {
			t.Errorf("Color failed with lang %s", lang)
		}

		v = s.HexColor()
		if v == "" {
			t.Errorf("HexColor failed with lang %s", lang)
		}

		v = s.HexColorShort()
		if v == "" {
			t.Errorf("HexColorShort failed with lang %s", lang)
		}

		v = s.DigitsN(2)
		if v == "" {
			t.Errorf("DigitsN failed with lang %s", lang)
		}

		v = s.Digits()
		if v == "" {
			t.Errorf("Digits failed with lang %s", lang)
		}
	}
}
