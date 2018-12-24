package pseudo

import (
	"testing"
)

func TestNames(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		s.SetLang(lang)

		v := s.MaleFirstName()
		if v == "" {
			t.Errorf("MaleFirstName failed with lang %s", lang)
		}

		v = s.FemaleFirstName()
		if v == "" {
			t.Errorf("FemaleFirstName failed with lang %s", lang)
		}

		v = s.FirstName()
		if v == "" {
			t.Errorf("FirstName failed with lang %s", lang)
		}

		v = s.MaleLastName()
		if v == "" {
			t.Errorf("MaleLastName failed with lang %s", lang)
		}

		v = s.FemaleLastName()
		if v == "" {
			t.Errorf("FemaleLastName failed with lang %s", lang)
		}

		v = s.LastName()
		if v == "" {
			t.Errorf("LastName failed with lang %s", lang)
		}

		v = s.MalePatronymic()
		if v == "" {
			t.Errorf("MalePatronymic failed with lang %s", lang)
		}

		v = s.FemalePatronymic()
		if v == "" {
			t.Errorf("FemalePatronymic failed with lang %s", lang)
		}

		v = s.Patronymic()
		if v == "" {
			t.Errorf("Patronymic failed with lang %s", lang)
		}

		v = s.MaleFullNameWithPrefix()
		if v == "" {
			t.Errorf("MaleFullNameWithPrefix failed with lang %s", lang)
		}

		v = s.FemaleFullNameWithPrefix()
		if v == "" {
			t.Errorf("FemaleFullNameWithPrefix failed with lang %s", lang)
		}

		v = s.FullNameWithPrefix()
		if v == "" {
			t.Errorf("FullNameWithPrefix failed with lang %s", lang)
		}

		v = s.MaleFullNameWithSuffix()
		if v == "" {
			t.Errorf("MaleFullNameWithSuffix failed with lang %s", lang)
		}

		v = s.FemaleFullNameWithSuffix()
		if v == "" {
			t.Errorf("FemaleFullNameWithSuffix failed with lang %s", lang)
		}

		v = s.FullNameWithSuffix()
		if v == "" {
			t.Errorf("FullNameWithSuffix failed with lang %s", lang)
		}

		v = s.MaleFullName()
		if v == "" {
			t.Errorf("MaleFullName failed with lang %s", lang)
		}

		v = s.FemaleFullName()
		if v == "" {
			t.Errorf("FemaleFullName failed with lang %s", lang)
		}

		v = s.FullName()
		if v == "" {
			t.Errorf("FullName failed with lang %s", lang)
		}
	}
}
