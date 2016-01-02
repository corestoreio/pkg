// +build ignore

package groupedproduct

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCheckoutCartGroupedProductImage => Grouped Product Image.
// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
var PathCheckoutCartGroupedProductImage = model.NewStr(`checkout/cart/grouped_product_image`, model.WithPkgCfg(PackageConfiguration))
