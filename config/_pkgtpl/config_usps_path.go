// +build ignore

package usps

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersUspsActive => Enabled for Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsActive = model.NewBool(`carriers/usps/active`)

// PathCarriersUspsActiveRma => Enabled for RMA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsActiveRma = model.NewBool(`carriers/usps/active_rma`)

// PathCarriersUspsGatewayUrl => Gateway URL.
var PathCarriersUspsGatewayUrl = model.NewStr(`carriers/usps/gateway_url`)

// PathCarriersUspsGatewaySecureUrl => Secure Gateway URL.
var PathCarriersUspsGatewaySecureUrl = model.NewStr(`carriers/usps/gateway_secure_url`)

// PathCarriersUspsTitle => Title.
var PathCarriersUspsTitle = model.NewStr(`carriers/usps/title`)

// PathCarriersUspsUserid => User ID.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUspsUserid = model.NewStr(`carriers/usps/userid`)

// PathCarriersUspsPassword => Password.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUspsPassword = model.NewStr(`carriers/usps/password`)

// PathCarriersUspsMode => Mode.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
var PathCarriersUspsMode = model.NewStr(`carriers/usps/mode`)

// PathCarriersUspsShipmentRequesttype => Packages Request Type.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
var PathCarriersUspsShipmentRequesttype = model.NewStr(`carriers/usps/shipment_requesttype`)

// PathCarriersUspsContainer => Container.
// SourceModel: Otnegam\Usps\Model\Source\Container
var PathCarriersUspsContainer = model.NewStr(`carriers/usps/container`)

// PathCarriersUspsSize => Size.
// SourceModel: Otnegam\Usps\Model\Source\Size
var PathCarriersUspsSize = model.NewStr(`carriers/usps/size`)

// PathCarriersUspsWidth => Width.
var PathCarriersUspsWidth = model.NewStr(`carriers/usps/width`)

// PathCarriersUspsLength => Length.
var PathCarriersUspsLength = model.NewStr(`carriers/usps/length`)

// PathCarriersUspsHeight => Height.
var PathCarriersUspsHeight = model.NewStr(`carriers/usps/height`)

// PathCarriersUspsGirth => Girth.
var PathCarriersUspsGirth = model.NewStr(`carriers/usps/girth`)

// PathCarriersUspsMachinable => Machinable.
// SourceModel: Otnegam\Usps\Model\Source\Machinable
var PathCarriersUspsMachinable = model.NewStr(`carriers/usps/machinable`)

// PathCarriersUspsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
var PathCarriersUspsMaxPackageWeight = model.NewStr(`carriers/usps/max_package_weight`)

// PathCarriersUspsHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersUspsHandlingType = model.NewStr(`carriers/usps/handling_type`)

// PathCarriersUspsHandlingAction => Handling Applied.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
var PathCarriersUspsHandlingAction = model.NewStr(`carriers/usps/handling_action`)

// PathCarriersUspsHandlingFee => Handling Fee.
var PathCarriersUspsHandlingFee = model.NewStr(`carriers/usps/handling_fee`)

// PathCarriersUspsAllowedMethods => Allowed Methods.
// SourceModel: Otnegam\Usps\Model\Source\Method
var PathCarriersUspsAllowedMethods = model.NewStringCSV(`carriers/usps/allowed_methods`)

// PathCarriersUspsFreeMethod => Free Method.
// SourceModel: Otnegam\Usps\Model\Source\Freemethod
var PathCarriersUspsFreeMethod = model.NewStr(`carriers/usps/free_method`)

// PathCarriersUspsFreeShippingEnable => Free Shipping Amount Threshold.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCarriersUspsFreeShippingEnable = model.NewBool(`carriers/usps/free_shipping_enable`)

// PathCarriersUspsFreeShippingSubtotal => Free Shipping Amount Threshold.
var PathCarriersUspsFreeShippingSubtotal = model.NewStr(`carriers/usps/free_shipping_subtotal`)

// PathCarriersUspsSpecificerrmsg => Displayed Error Message.
var PathCarriersUspsSpecificerrmsg = model.NewStr(`carriers/usps/specificerrmsg`)

// PathCarriersUspsSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersUspsSallowspecific = model.NewStr(`carriers/usps/sallowspecific`)

// PathCarriersUspsSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersUspsSpecificcountry = model.NewStringCSV(`carriers/usps/specificcountry`)

// PathCarriersUspsDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsDebug = model.NewBool(`carriers/usps/debug`)

// PathCarriersUspsShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsShowmethod = model.NewBool(`carriers/usps/showmethod`)

// PathCarriersUspsSortOrder => Sort Order.
var PathCarriersUspsSortOrder = model.NewStr(`carriers/usps/sort_order`)
