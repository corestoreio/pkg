// +build ignore

package configurableproduct

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CheckoutCartConfigurableProductImage => Configurable Product Image.
	// Path: checkout/cart/configurable_product_image
	// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
	CheckoutCartConfigurableProductImage model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutCartConfigurableProductImage = model.NewStr(`checkout/cart/configurable_product_image`, model.WithConfigStructure(cfgStruct))

	return pp
}
