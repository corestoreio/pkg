// +build ignore

package multishipping

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathMultishippingOptionsCheckoutMultiple => Allow Shipping to Multiple Addresses.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathMultishippingOptionsCheckoutMultiple = model.NewBool(`multishipping/options/checkout_multiple`, model.WithPkgCfg(PackageConfiguration))

// PathMultishippingOptionsCheckoutMultipleMaximumQty => Maximum Qty Allowed for Shipping to Multiple Addresses.
var PathMultishippingOptionsCheckoutMultipleMaximumQty = model.NewStr(`multishipping/options/checkout_multiple_maximum_qty`, model.WithPkgCfg(PackageConfiguration))
