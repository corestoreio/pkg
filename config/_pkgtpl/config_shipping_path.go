// +build ignore

package shipping

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathShippingOriginCountryId => Country.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathShippingOriginCountryId = model.NewStr(`shipping/origin/country_id`)

// PathShippingOriginRegionId => Region/State.
var PathShippingOriginRegionId = model.NewStr(`shipping/origin/region_id`)

// PathShippingOriginPostcode => ZIP/Postal Code.
var PathShippingOriginPostcode = model.NewStr(`shipping/origin/postcode`)

// PathShippingOriginCity => City.
var PathShippingOriginCity = model.NewStr(`shipping/origin/city`)

// PathShippingOriginStreetLine1 => Street Address.
var PathShippingOriginStreetLine1 = model.NewStr(`shipping/origin/street_line1`)

// PathShippingOriginStreetLine2 => Street Address Line 2.
var PathShippingOriginStreetLine2 = model.NewStr(`shipping/origin/street_line2`)

// PathShippingShippingPolicyEnableShippingPolicy => Apply custom Shipping Policy.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathShippingShippingPolicyEnableShippingPolicy = model.NewBool(`shipping/shipping_policy/enable_shipping_policy`)

// PathShippingShippingPolicyShippingPolicyContent => Shipping Policy.
var PathShippingShippingPolicyShippingPolicyContent = model.NewStr(`shipping/shipping_policy/shipping_policy_content`)
