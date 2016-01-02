// +build ignore

package paypal

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPaypalGeneralMerchantCountry => Merchant Country.
// If not specified, Default Country from General Config will be used
// BackendModel: Otnegam\Paypal\Model\System\Config\Backend\MerchantCountry
// SourceModel: Otnegam\Paypal\Model\System\Config\Source\MerchantCountry
var PathPaypalGeneralMerchantCountry = model.NewStr(`paypal/general/merchant_country`, model.WithPkgCfg(PackageConfiguration))
