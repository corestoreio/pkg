// +build ignore

package fedex

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersFedexActive => Enabled for Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexActive = model.NewBool(`carriers/fedex/active`)

// PathCarriersFedexActiveRma => Enabled for RMA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexActiveRma = model.NewBool(`carriers/fedex/active_rma`)

// PathCarriersFedexTitle => Title.
var PathCarriersFedexTitle = model.NewStr(`carriers/fedex/title`)

// PathCarriersFedexAccount => Account ID.
// Please make sure to use only digits here. No dashes are allowed.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersFedexAccount = model.NewStr(`carriers/fedex/account`)

// PathCarriersFedexMeterNumber => Meter Number.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersFedexMeterNumber = model.NewStr(`carriers/fedex/meter_number`)

// PathCarriersFedexKey => Key.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersFedexKey = model.NewStr(`carriers/fedex/key`)

// PathCarriersFedexPassword => Password.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersFedexPassword = model.NewStr(`carriers/fedex/password`)

// PathCarriersFedexSandboxMode => Sandbox Mode.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexSandboxMode = model.NewBool(`carriers/fedex/sandbox_mode`)

// PathCarriersFedexProductionWebservicesUrl => Web-Services URL (Production).
var PathCarriersFedexProductionWebservicesUrl = model.NewStr(`carriers/fedex/production_webservices_url`)

// PathCarriersFedexSandboxWebservicesUrl => Web-Services URL (Sandbox).
var PathCarriersFedexSandboxWebservicesUrl = model.NewStr(`carriers/fedex/sandbox_webservices_url`)

// PathCarriersFedexShipmentRequesttype => Packages Request Type.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
var PathCarriersFedexShipmentRequesttype = model.NewStr(`carriers/fedex/shipment_requesttype`)

// PathCarriersFedexPackaging => Packaging.
// SourceModel: Otnegam\Fedex\Model\Source\Packaging
var PathCarriersFedexPackaging = model.NewStr(`carriers/fedex/packaging`)

// PathCarriersFedexDropoff => Dropoff.
// SourceModel: Otnegam\Fedex\Model\Source\Dropoff
var PathCarriersFedexDropoff = model.NewStr(`carriers/fedex/dropoff`)

// PathCarriersFedexUnitOfMeasure => Weight Unit.
// SourceModel: Otnegam\Fedex\Model\Source\Unitofmeasure
var PathCarriersFedexUnitOfMeasure = model.NewStr(`carriers/fedex/unit_of_measure`)

// PathCarriersFedexMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
var PathCarriersFedexMaxPackageWeight = model.NewStr(`carriers/fedex/max_package_weight`)

// PathCarriersFedexHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersFedexHandlingType = model.NewStr(`carriers/fedex/handling_type`)

// PathCarriersFedexHandlingAction => Handling Applied.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
var PathCarriersFedexHandlingAction = model.NewStr(`carriers/fedex/handling_action`)

// PathCarriersFedexHandlingFee => Handling Fee.
var PathCarriersFedexHandlingFee = model.NewStr(`carriers/fedex/handling_fee`)

// PathCarriersFedexResidenceDelivery => Residential Delivery.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexResidenceDelivery = model.NewBool(`carriers/fedex/residence_delivery`)

// PathCarriersFedexAllowedMethods => Allowed Methods.
// SourceModel: Otnegam\Fedex\Model\Source\Method
var PathCarriersFedexAllowedMethods = model.NewStringCSV(`carriers/fedex/allowed_methods`)

// PathCarriersFedexSmartpostHubid => Hub ID.
// The field is applicable if the Smart Post method is selected.
var PathCarriersFedexSmartpostHubid = model.NewStr(`carriers/fedex/smartpost_hubid`)

// PathCarriersFedexFreeMethod => Free Method.
// SourceModel: Otnegam\Fedex\Model\Source\Freemethod
var PathCarriersFedexFreeMethod = model.NewStr(`carriers/fedex/free_method`)

// PathCarriersFedexFreeShippingEnable => Free Shipping Amount Threshold.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCarriersFedexFreeShippingEnable = model.NewBool(`carriers/fedex/free_shipping_enable`)

// PathCarriersFedexFreeShippingSubtotal => Free Shipping Amount Threshold.
var PathCarriersFedexFreeShippingSubtotal = model.NewStr(`carriers/fedex/free_shipping_subtotal`)

// PathCarriersFedexSpecificerrmsg => Displayed Error Message.
var PathCarriersFedexSpecificerrmsg = model.NewStr(`carriers/fedex/specificerrmsg`)

// PathCarriersFedexSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersFedexSallowspecific = model.NewStr(`carriers/fedex/sallowspecific`)

// PathCarriersFedexSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersFedexSpecificcountry = model.NewStringCSV(`carriers/fedex/specificcountry`)

// PathCarriersFedexDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexDebug = model.NewBool(`carriers/fedex/debug`)

// PathCarriersFedexShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersFedexShowmethod = model.NewBool(`carriers/fedex/showmethod`)

// PathCarriersFedexSortOrder => Sort Order.
var PathCarriersFedexSortOrder = model.NewStr(`carriers/fedex/sort_order`)
