// +build ignore

package msrp

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSalesMsrpEnabled => Enable MAP.
// Warning! Enabling MAP by default will hide all product prices on
// Storefront.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMsrpEnabled = model.NewBool(`sales/msrp/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMsrpDisplayPriceType => Display Actual Price.
// SourceModel: Otnegam\Msrp\Model\Product\Attribute\Source\Type
var PathSalesMsrpDisplayPriceType = model.NewStr(`sales/msrp/display_price_type`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMsrpExplanationMessage => Default Popup Text Message.
var PathSalesMsrpExplanationMessage = model.NewStr(`sales/msrp/explanation_message`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMsrpExplanationMessageWhatsThis => Default "What's This" Text Message.
var PathSalesMsrpExplanationMessageWhatsThis = model.NewStr(`sales/msrp/explanation_message_whats_this`, model.WithPkgCfg(PackageConfiguration))
