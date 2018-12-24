package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestDates(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		v := s.WeekDay()
		if v == "" {
			t.Errorf("WeekDay failed with lang %s", lang)
		}

		v = s.WeekDayShort()
		if v == "" {
			t.Errorf("WeekDayShort failed with lang %s", lang)
		}

		n := s.WeekdayNum()
		if n < 0 || n > 7 {
			t.Errorf("WeekdayNum failed with lang %s", lang)
		}

		v = s.Month()
		if v == "" {
			t.Errorf("Month failed with lang %s", lang)
		}

		v = s.MonthShort()
		if v == "" {
			t.Errorf("MonthShort failed with lang %s", lang)
		}

		n = s.MonthNum()
		if n < 0 || n > 31 {
			t.Errorf("MonthNum failed with lang %s", lang)
		}

		n = s.Year(1950, 2020)
		if n < 1950 || n > 2020 {
			t.Errorf("Year failed with lang %s", lang)
		}
	}
}
