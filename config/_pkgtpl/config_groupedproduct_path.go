// +build ignore

package groupedproduct

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
	// CheckoutCartGroupedProductImage => Grouped Product Image.
	// Path: checkout/cart/grouped_product_image
	// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
	CheckoutCartGroupedProductImage model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutCartGroupedProductImage = model.NewStr(`checkout/cart/grouped_product_image`, model.WithPkgCfg(pkgCfg))

	return pp
}
