// +build ignore

package paypal

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// PaypalGeneralMerchantCountry => Merchant Country.
	// If not specified, Default Country from General Config will be used
	// Path: paypal/general/merchant_country
	// BackendModel: Magento\Paypal\Model\System\Config\Backend\MerchantCountry
	// SourceModel: Magento\Paypal\Model\System\Config\Source\MerchantCountry
	PaypalGeneralMerchantCountry cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaypalGeneralMerchantCountry = cfgmodel.NewStr(`paypal/general/merchant_country`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
