package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestCreditCards(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))
		tp := s.CreditCardType()
		assert.NotEmpty(t, tp, "s.CreditCardType()")
		n1 := s.CreditCardNum("")
		assert.Regexp(t, "^[0-9]{13,}$", n1, "s.CreditCardNum()")
		nVisa := s.CreditCardNum("visa")
		assert.Regexp(t, "^[0-9]{13,}$", nVisa, "s.CreditCardNum(visa)")
		// t.Logf("%q %q Visa: %q", tp, n1, nVisa)
	}
}
