// +build ignore

package swatches

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogFrontendSwatchesPerProduct => Swatches per Product.
var PathCatalogFrontendSwatchesPerProduct = model.NewStr(`catalog/frontend/swatches_per_product`, model.WithPkgCfg(PackageConfiguration))
