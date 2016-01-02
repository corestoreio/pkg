// +build ignore

package usps

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersUspsActive => Enabled for Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsActive = model.NewBool(`carriers/usps/active`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsActiveRma => Enabled for RMA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsActiveRma = model.NewBool(`carriers/usps/active_rma`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsGatewayUrl => Gateway URL.
var PathCarriersUspsGatewayUrl = model.NewStr(`carriers/usps/gateway_url`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsGatewaySecureUrl => Secure Gateway URL.
var PathCarriersUspsGatewaySecureUrl = model.NewStr(`carriers/usps/gateway_secure_url`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsTitle => Title.
var PathCarriersUspsTitle = model.NewStr(`carriers/usps/title`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsUserid => User ID.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUspsUserid = model.NewStr(`carriers/usps/userid`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsPassword => Password.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUspsPassword = model.NewStr(`carriers/usps/password`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsMode => Mode.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
var PathCarriersUspsMode = model.NewStr(`carriers/usps/mode`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsShipmentRequesttype => Packages Request Type.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
var PathCarriersUspsShipmentRequesttype = model.NewStr(`carriers/usps/shipment_requesttype`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsContainer => Container.
// SourceModel: Otnegam\Usps\Model\Source\Container
var PathCarriersUspsContainer = model.NewStr(`carriers/usps/container`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsSize => Size.
// SourceModel: Otnegam\Usps\Model\Source\Size
var PathCarriersUspsSize = model.NewStr(`carriers/usps/size`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsWidth => Width.
var PathCarriersUspsWidth = model.NewStr(`carriers/usps/width`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsLength => Length.
var PathCarriersUspsLength = model.NewStr(`carriers/usps/length`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsHeight => Height.
var PathCarriersUspsHeight = model.NewStr(`carriers/usps/height`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsGirth => Girth.
var PathCarriersUspsGirth = model.NewStr(`carriers/usps/girth`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsMachinable => Machinable.
// SourceModel: Otnegam\Usps\Model\Source\Machinable
var PathCarriersUspsMachinable = model.NewStr(`carriers/usps/machinable`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
var PathCarriersUspsMaxPackageWeight = model.NewStr(`carriers/usps/max_package_weight`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersUspsHandlingType = model.NewStr(`carriers/usps/handling_type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsHandlingAction => Handling Applied.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
var PathCarriersUspsHandlingAction = model.NewStr(`carriers/usps/handling_action`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsHandlingFee => Handling Fee.
var PathCarriersUspsHandlingFee = model.NewStr(`carriers/usps/handling_fee`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsAllowedMethods => Allowed Methods.
// SourceModel: Otnegam\Usps\Model\Source\Method
var PathCarriersUspsAllowedMethods = model.NewStringCSV(`carriers/usps/allowed_methods`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsFreeMethod => Free Method.
// SourceModel: Otnegam\Usps\Model\Source\Freemethod
var PathCarriersUspsFreeMethod = model.NewStr(`carriers/usps/free_method`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsFreeShippingEnable => Free Shipping Amount Threshold.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCarriersUspsFreeShippingEnable = model.NewBool(`carriers/usps/free_shipping_enable`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsFreeShippingSubtotal => Free Shipping Amount Threshold.
var PathCarriersUspsFreeShippingSubtotal = model.NewStr(`carriers/usps/free_shipping_subtotal`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsSpecificerrmsg => Displayed Error Message.
var PathCarriersUspsSpecificerrmsg = model.NewStr(`carriers/usps/specificerrmsg`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersUspsSallowspecific = model.NewStr(`carriers/usps/sallowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersUspsSpecificcountry = model.NewStringCSV(`carriers/usps/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsDebug = model.NewBool(`carriers/usps/debug`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUspsShowmethod = model.NewBool(`carriers/usps/showmethod`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersUspsSortOrder => Sort Order.
var PathCarriersUspsSortOrder = model.NewStr(`carriers/usps/sort_order`, model.WithPkgCfg(PackageConfiguration))
