// +build ignore

package shipping

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
	// ShippingOriginCountryId => Country.
	// Path: shipping/origin/country_id
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	ShippingOriginCountryId cfgmodel.Str

	// ShippingOriginRegionId => Region/State.
	// Path: shipping/origin/region_id
	ShippingOriginRegionId cfgmodel.Str

	// ShippingOriginPostcode => ZIP/Postal Code.
	// Path: shipping/origin/postcode
	ShippingOriginPostcode cfgmodel.Str

	// ShippingOriginCity => City.
	// Path: shipping/origin/city
	ShippingOriginCity cfgmodel.Str

	// ShippingOriginStreetLine1 => Street Address.
	// Path: shipping/origin/street_line1
	ShippingOriginStreetLine1 cfgmodel.Str

	// ShippingOriginStreetLine2 => Street Address Line 2.
	// Path: shipping/origin/street_line2
	ShippingOriginStreetLine2 cfgmodel.Str

	// ShippingShippingPolicyEnableShippingPolicy => Apply custom Shipping Policy.
	// Path: shipping/shipping_policy/enable_shipping_policy
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	ShippingShippingPolicyEnableShippingPolicy cfgmodel.Bool

	// ShippingShippingPolicyShippingPolicyContent => Shipping Policy.
	// Path: shipping/shipping_policy/shipping_policy_content
	ShippingShippingPolicyShippingPolicyContent cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.ShippingOriginCountryId = cfgmodel.NewStr(`shipping/origin/country_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingOriginRegionId = cfgmodel.NewStr(`shipping/origin/region_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingOriginPostcode = cfgmodel.NewStr(`shipping/origin/postcode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingOriginCity = cfgmodel.NewStr(`shipping/origin/city`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingOriginStreetLine1 = cfgmodel.NewStr(`shipping/origin/street_line1`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingOriginStreetLine2 = cfgmodel.NewStr(`shipping/origin/street_line2`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingShippingPolicyEnableShippingPolicy = cfgmodel.NewBool(`shipping/shipping_policy/enable_shipping_policy`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ShippingShippingPolicyShippingPolicyContent = cfgmodel.NewStr(`shipping/shipping_policy/shipping_policy_content`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
