// +build ignore

package multishipping

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// MultishippingOptionsCheckoutMultiple => Allow Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	MultishippingOptionsCheckoutMultiple cfgmodel.Bool

	// MultishippingOptionsCheckoutMultipleMaximumQty => Maximum Qty Allowed for Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple_maximum_qty
	MultishippingOptionsCheckoutMultipleMaximumQty cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.MultishippingOptionsCheckoutMultiple = cfgmodel.NewBool(`multishipping/options/checkout_multiple`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.MultishippingOptionsCheckoutMultipleMaximumQty = cfgmodel.NewStr(`multishipping/options/checkout_multiple_maximum_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
