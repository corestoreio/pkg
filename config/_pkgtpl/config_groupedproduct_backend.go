// +build ignore

package groupedproduct

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CheckoutCartGroupedProductImage => Grouped Product Image.
	// Path: checkout/cart/grouped_product_image
	// SourceModel: Magento\Catalog\Model\Config\Source\Product\Thumbnail
	CheckoutCartGroupedProductImage model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutCartGroupedProductImage = model.NewStr(`checkout/cart/grouped_product_image`, model.WithConfigStructure(cfgStruct))

	return pp
}
