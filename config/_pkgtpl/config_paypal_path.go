// +build ignore

package paypal

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// PaypalGeneralMerchantCountry => Merchant Country.
	// If not specified, Default Country from General Config will be used
	// Path: paypal/general/merchant_country
	// BackendModel: Otnegam\Paypal\Model\System\Config\Backend\MerchantCountry
	// SourceModel: Otnegam\Paypal\Model\System\Config\Source\MerchantCountry
	PaypalGeneralMerchantCountry model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PaypalGeneralMerchantCountry = model.NewStr(`paypal/general/merchant_country`, model.WithPkgCfg(pkgCfg))

	return pp
}
