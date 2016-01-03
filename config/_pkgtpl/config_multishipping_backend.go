// +build ignore

package multishipping

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
	// MultishippingOptionsCheckoutMultiple => Allow Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	MultishippingOptionsCheckoutMultiple model.Bool

	// MultishippingOptionsCheckoutMultipleMaximumQty => Maximum Qty Allowed for Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple_maximum_qty
	MultishippingOptionsCheckoutMultipleMaximumQty model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.MultishippingOptionsCheckoutMultiple = model.NewBool(`multishipping/options/checkout_multiple`, model.WithConfigStructure(cfgStruct))
	pp.MultishippingOptionsCheckoutMultipleMaximumQty = model.NewStr(`multishipping/options/checkout_multiple_maximum_qty`, model.WithConfigStructure(cfgStruct))

	return pp
}
