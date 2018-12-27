package pseudo

import (
	"testing"
)

func TestProducts(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		s.SetLang(lang)

		v := s.Brand()
		if v == "" {
			t.Errorf("Brand failed with lang %s", lang)
		}

		v = s.ProductName()
		if v == "" {
			t.Errorf("ProductName failed with lang %s", lang)
		}

		v = s.Product()
		if v == "" {
			t.Errorf("Product failed with lang %s", lang)
		}

		v = s.Model()
		if v == "" {
			t.Errorf("Model failed with lang %s", lang)
		}
	}
}
