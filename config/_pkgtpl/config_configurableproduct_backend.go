// +build ignore

package configurableproduct

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CheckoutCartConfigurableProductImage => Configurable Product Image.
	// Path: checkout/cart/configurable_product_image
	// SourceModel: Magento\Catalog\Model\Config\Source\Product\Thumbnail
	CheckoutCartConfigurableProductImage cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutCartConfigurableProductImage = cfgmodel.NewStr(`checkout/cart/configurable_product_image`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
