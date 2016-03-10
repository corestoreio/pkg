// +build ignore

package fedex

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CarriersFedexActive => Enabled for Checkout.
	// Path: carriers/fedex/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexActive cfgmodel.Bool

	// CarriersFedexActiveRma => Enabled for RMA.
	// Path: carriers/fedex/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexActiveRma cfgmodel.Bool

	// CarriersFedexTitle => Title.
	// Path: carriers/fedex/title
	CarriersFedexTitle cfgmodel.Str

	// CarriersFedexAccount => Account ID.
	// Please make sure to use only digits here. No dashes are allowed.
	// Path: carriers/fedex/account
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexAccount cfgmodel.Str

	// CarriersFedexMeterNumber => Meter Number.
	// Path: carriers/fedex/meter_number
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexMeterNumber cfgmodel.Str

	// CarriersFedexKey => Key.
	// Path: carriers/fedex/key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexKey cfgmodel.Str

	// CarriersFedexPassword => Password.
	// Path: carriers/fedex/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexPassword cfgmodel.Str

	// CarriersFedexSandboxMode => Sandbox Mode.
	// Path: carriers/fedex/sandbox_mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexSandboxMode cfgmodel.Bool

	// CarriersFedexProductionWebservicesUrl => Web-Services URL (Production).
	// Path: carriers/fedex/production_webservices_url
	CarriersFedexProductionWebservicesUrl cfgmodel.Str

	// CarriersFedexSandboxWebservicesUrl => Web-Services URL (Sandbox).
	// Path: carriers/fedex/sandbox_webservices_url
	CarriersFedexSandboxWebservicesUrl cfgmodel.Str

	// CarriersFedexShipmentRequesttype => Packages Request Type.
	// Path: carriers/fedex/shipment_requesttype
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
	CarriersFedexShipmentRequesttype cfgmodel.Str

	// CarriersFedexPackaging => Packaging.
	// Path: carriers/fedex/packaging
	// SourceModel: Magento\Fedex\Model\Source\Packaging
	CarriersFedexPackaging cfgmodel.Str

	// CarriersFedexDropoff => Dropoff.
	// Path: carriers/fedex/dropoff
	// SourceModel: Magento\Fedex\Model\Source\Dropoff
	CarriersFedexDropoff cfgmodel.Str

	// CarriersFedexUnitOfMeasure => Weight Unit.
	// Path: carriers/fedex/unit_of_measure
	// SourceModel: Magento\Fedex\Model\Source\Unitofmeasure
	CarriersFedexUnitOfMeasure cfgmodel.Str

	// CarriersFedexMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/fedex/max_package_weight
	CarriersFedexMaxPackageWeight cfgmodel.Str

	// CarriersFedexHandlingType => Calculate Handling Fee.
	// Path: carriers/fedex/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersFedexHandlingType cfgmodel.Str

	// CarriersFedexHandlingAction => Handling Applied.
	// Path: carriers/fedex/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
	CarriersFedexHandlingAction cfgmodel.Str

	// CarriersFedexHandlingFee => Handling Fee.
	// Path: carriers/fedex/handling_fee
	CarriersFedexHandlingFee cfgmodel.Str

	// CarriersFedexResidenceDelivery => Residential Delivery.
	// Path: carriers/fedex/residence_delivery
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexResidenceDelivery cfgmodel.Bool

	// CarriersFedexAllowedMethods => Allowed Methods.
	// Path: carriers/fedex/allowed_methods
	// SourceModel: Magento\Fedex\Model\Source\Method
	CarriersFedexAllowedMethods cfgmodel.StringCSV

	// CarriersFedexSmartpostHubid => Hub ID.
	// The field is applicable if the Smart Post method is selected.
	// Path: carriers/fedex/smartpost_hubid
	CarriersFedexSmartpostHubid cfgmodel.Str

	// CarriersFedexFreeMethod => Free Method.
	// Path: carriers/fedex/free_method
	// SourceModel: Magento\Fedex\Model\Source\Freemethod
	CarriersFedexFreeMethod cfgmodel.Str

	// CarriersFedexFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersFedexFreeShippingEnable cfgmodel.Bool

	// CarriersFedexFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_subtotal
	CarriersFedexFreeShippingSubtotal cfgmodel.Str

	// CarriersFedexSpecificerrmsg => Displayed Error Message.
	// Path: carriers/fedex/specificerrmsg
	CarriersFedexSpecificerrmsg cfgmodel.Str

	// CarriersFedexSallowspecific => Ship to Applicable Countries.
	// Path: carriers/fedex/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFedexSallowspecific cfgmodel.Str

	// CarriersFedexSpecificcountry => Ship to Specific Countries.
	// Path: carriers/fedex/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFedexSpecificcountry cfgmodel.StringCSV

	// CarriersFedexDebug => Debug.
	// Path: carriers/fedex/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexDebug cfgmodel.Bool

	// CarriersFedexShowmethod => Show Method if Not Applicable.
	// Path: carriers/fedex/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexShowmethod cfgmodel.Bool

	// CarriersFedexSortOrder => Sort Order.
	// Path: carriers/fedex/sort_order
	CarriersFedexSortOrder cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFedexActive = cfgmodel.NewBool(`carriers/fedex/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexActiveRma = cfgmodel.NewBool(`carriers/fedex/active_rma`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexTitle = cfgmodel.NewStr(`carriers/fedex/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexAccount = cfgmodel.NewStr(`carriers/fedex/account`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexMeterNumber = cfgmodel.NewStr(`carriers/fedex/meter_number`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexKey = cfgmodel.NewStr(`carriers/fedex/key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexPassword = cfgmodel.NewStr(`carriers/fedex/password`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSandboxMode = cfgmodel.NewBool(`carriers/fedex/sandbox_mode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexProductionWebservicesUrl = cfgmodel.NewStr(`carriers/fedex/production_webservices_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSandboxWebservicesUrl = cfgmodel.NewStr(`carriers/fedex/sandbox_webservices_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexShipmentRequesttype = cfgmodel.NewStr(`carriers/fedex/shipment_requesttype`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexPackaging = cfgmodel.NewStr(`carriers/fedex/packaging`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexDropoff = cfgmodel.NewStr(`carriers/fedex/dropoff`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexUnitOfMeasure = cfgmodel.NewStr(`carriers/fedex/unit_of_measure`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexMaxPackageWeight = cfgmodel.NewStr(`carriers/fedex/max_package_weight`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingType = cfgmodel.NewStr(`carriers/fedex/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingAction = cfgmodel.NewStr(`carriers/fedex/handling_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingFee = cfgmodel.NewStr(`carriers/fedex/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexResidenceDelivery = cfgmodel.NewBool(`carriers/fedex/residence_delivery`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexAllowedMethods = cfgmodel.NewStringCSV(`carriers/fedex/allowed_methods`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSmartpostHubid = cfgmodel.NewStr(`carriers/fedex/smartpost_hubid`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeMethod = cfgmodel.NewStr(`carriers/fedex/free_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeShippingEnable = cfgmodel.NewBool(`carriers/fedex/free_shipping_enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeShippingSubtotal = cfgmodel.NewStr(`carriers/fedex/free_shipping_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSpecificerrmsg = cfgmodel.NewStr(`carriers/fedex/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSallowspecific = cfgmodel.NewStr(`carriers/fedex/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSpecificcountry = cfgmodel.NewStringCSV(`carriers/fedex/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexDebug = cfgmodel.NewBool(`carriers/fedex/debug`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexShowmethod = cfgmodel.NewBool(`carriers/fedex/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSortOrder = cfgmodel.NewStr(`carriers/fedex/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
