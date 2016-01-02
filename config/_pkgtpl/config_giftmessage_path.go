// +build ignore

package giftmessage

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSalesGiftOptionsAllowOrder => Allow Gift Messages on Order Level.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesGiftOptionsAllowOrder = model.NewBool(`sales/gift_options/allow_order`, model.WithPkgCfg(PackageConfiguration))

// PathSalesGiftOptionsAllowItems => Allow Gift Messages for Order Items.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesGiftOptionsAllowItems = model.NewBool(`sales/gift_options/allow_items`, model.WithPkgCfg(PackageConfiguration))
