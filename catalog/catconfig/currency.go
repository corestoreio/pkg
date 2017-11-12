package catconfig

import (
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/directory"
	"github.com/corestoreio/errors"
)

// BaseCurrency returns the base currency code of a website.
// 	1st argument should be a path to catalog/price/scope
// 	2nd argument should be a path to currency/options/base
func BaseCurrency(cr config.Getter, sg config.Scoped, ps PriceScope, cc directory.ConfigCurrency) (directory.Currency, error) {
	// TODO, and also see test: TestWebsiteBaseCurrency
	isGlobal, err := ps.IsGlobal(sg)
	if err != nil {
		return directory.Currency{}, errors.Wrap(err, "asdf")
	}
	if isGlobal {
		return cc.GetDefault(cr) // default scope
	}
	return cc.Get(sg) // website scope
}
