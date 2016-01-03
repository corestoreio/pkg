// +build ignore

package shipping

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
	// ShippingOriginCountryId => Country.
	// Path: shipping/origin/country_id
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	ShippingOriginCountryId model.Str

	// ShippingOriginRegionId => Region/State.
	// Path: shipping/origin/region_id
	ShippingOriginRegionId model.Str

	// ShippingOriginPostcode => ZIP/Postal Code.
	// Path: shipping/origin/postcode
	ShippingOriginPostcode model.Str

	// ShippingOriginCity => City.
	// Path: shipping/origin/city
	ShippingOriginCity model.Str

	// ShippingOriginStreetLine1 => Street Address.
	// Path: shipping/origin/street_line1
	ShippingOriginStreetLine1 model.Str

	// ShippingOriginStreetLine2 => Street Address Line 2.
	// Path: shipping/origin/street_line2
	ShippingOriginStreetLine2 model.Str

	// ShippingShippingPolicyEnableShippingPolicy => Apply custom Shipping Policy.
	// Path: shipping/shipping_policy/enable_shipping_policy
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	ShippingShippingPolicyEnableShippingPolicy model.Bool

	// ShippingShippingPolicyShippingPolicyContent => Shipping Policy.
	// Path: shipping/shipping_policy/shipping_policy_content
	ShippingShippingPolicyShippingPolicyContent model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.ShippingOriginCountryId = model.NewStr(`shipping/origin/country_id`, model.WithConfigStructure(cfgStruct))
	pp.ShippingOriginRegionId = model.NewStr(`shipping/origin/region_id`, model.WithConfigStructure(cfgStruct))
	pp.ShippingOriginPostcode = model.NewStr(`shipping/origin/postcode`, model.WithConfigStructure(cfgStruct))
	pp.ShippingOriginCity = model.NewStr(`shipping/origin/city`, model.WithConfigStructure(cfgStruct))
	pp.ShippingOriginStreetLine1 = model.NewStr(`shipping/origin/street_line1`, model.WithConfigStructure(cfgStruct))
	pp.ShippingOriginStreetLine2 = model.NewStr(`shipping/origin/street_line2`, model.WithConfigStructure(cfgStruct))
	pp.ShippingShippingPolicyEnableShippingPolicy = model.NewBool(`shipping/shipping_policy/enable_shipping_policy`, model.WithConfigStructure(cfgStruct))
	pp.ShippingShippingPolicyShippingPolicyContent = model.NewStr(`shipping/shipping_policy/shipping_policy_content`, model.WithConfigStructure(cfgStruct))

	return pp
}
