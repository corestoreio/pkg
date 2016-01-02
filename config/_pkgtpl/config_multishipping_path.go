// +build ignore

package multishipping

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
	// MultishippingOptionsCheckoutMultiple => Allow Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	MultishippingOptionsCheckoutMultiple model.Bool

	// MultishippingOptionsCheckoutMultipleMaximumQty => Maximum Qty Allowed for Shipping to Multiple Addresses.
	// Path: multishipping/options/checkout_multiple_maximum_qty
	MultishippingOptionsCheckoutMultipleMaximumQty model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.MultishippingOptionsCheckoutMultiple = model.NewBool(`multishipping/options/checkout_multiple`, model.WithPkgCfg(pkgCfg))
	pp.MultishippingOptionsCheckoutMultipleMaximumQty = model.NewStr(`multishipping/options/checkout_multiple_maximum_qty`, model.WithPkgCfg(pkgCfg))

	return pp
}
