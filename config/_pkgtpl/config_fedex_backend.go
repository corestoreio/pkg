// +build ignore

package fedex

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
	// CarriersFedexActive => Enabled for Checkout.
	// Path: carriers/fedex/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexActive model.Bool

	// CarriersFedexActiveRma => Enabled for RMA.
	// Path: carriers/fedex/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexActiveRma model.Bool

	// CarriersFedexTitle => Title.
	// Path: carriers/fedex/title
	CarriersFedexTitle model.Str

	// CarriersFedexAccount => Account ID.
	// Please make sure to use only digits here. No dashes are allowed.
	// Path: carriers/fedex/account
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexAccount model.Str

	// CarriersFedexMeterNumber => Meter Number.
	// Path: carriers/fedex/meter_number
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexMeterNumber model.Str

	// CarriersFedexKey => Key.
	// Path: carriers/fedex/key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexKey model.Str

	// CarriersFedexPassword => Password.
	// Path: carriers/fedex/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersFedexPassword model.Str

	// CarriersFedexSandboxMode => Sandbox Mode.
	// Path: carriers/fedex/sandbox_mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexSandboxMode model.Bool

	// CarriersFedexProductionWebservicesUrl => Web-Services URL (Production).
	// Path: carriers/fedex/production_webservices_url
	CarriersFedexProductionWebservicesUrl model.Str

	// CarriersFedexSandboxWebservicesUrl => Web-Services URL (Sandbox).
	// Path: carriers/fedex/sandbox_webservices_url
	CarriersFedexSandboxWebservicesUrl model.Str

	// CarriersFedexShipmentRequesttype => Packages Request Type.
	// Path: carriers/fedex/shipment_requesttype
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
	CarriersFedexShipmentRequesttype model.Str

	// CarriersFedexPackaging => Packaging.
	// Path: carriers/fedex/packaging
	// SourceModel: Magento\Fedex\Model\Source\Packaging
	CarriersFedexPackaging model.Str

	// CarriersFedexDropoff => Dropoff.
	// Path: carriers/fedex/dropoff
	// SourceModel: Magento\Fedex\Model\Source\Dropoff
	CarriersFedexDropoff model.Str

	// CarriersFedexUnitOfMeasure => Weight Unit.
	// Path: carriers/fedex/unit_of_measure
	// SourceModel: Magento\Fedex\Model\Source\Unitofmeasure
	CarriersFedexUnitOfMeasure model.Str

	// CarriersFedexMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/fedex/max_package_weight
	CarriersFedexMaxPackageWeight model.Str

	// CarriersFedexHandlingType => Calculate Handling Fee.
	// Path: carriers/fedex/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersFedexHandlingType model.Str

	// CarriersFedexHandlingAction => Handling Applied.
	// Path: carriers/fedex/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
	CarriersFedexHandlingAction model.Str

	// CarriersFedexHandlingFee => Handling Fee.
	// Path: carriers/fedex/handling_fee
	CarriersFedexHandlingFee model.Str

	// CarriersFedexResidenceDelivery => Residential Delivery.
	// Path: carriers/fedex/residence_delivery
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexResidenceDelivery model.Bool

	// CarriersFedexAllowedMethods => Allowed Methods.
	// Path: carriers/fedex/allowed_methods
	// SourceModel: Magento\Fedex\Model\Source\Method
	CarriersFedexAllowedMethods model.StringCSV

	// CarriersFedexSmartpostHubid => Hub ID.
	// The field is applicable if the Smart Post method is selected.
	// Path: carriers/fedex/smartpost_hubid
	CarriersFedexSmartpostHubid model.Str

	// CarriersFedexFreeMethod => Free Method.
	// Path: carriers/fedex/free_method
	// SourceModel: Magento\Fedex\Model\Source\Freemethod
	CarriersFedexFreeMethod model.Str

	// CarriersFedexFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersFedexFreeShippingEnable model.Bool

	// CarriersFedexFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_subtotal
	CarriersFedexFreeShippingSubtotal model.Str

	// CarriersFedexSpecificerrmsg => Displayed Error Message.
	// Path: carriers/fedex/specificerrmsg
	CarriersFedexSpecificerrmsg model.Str

	// CarriersFedexSallowspecific => Ship to Applicable Countries.
	// Path: carriers/fedex/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFedexSallowspecific model.Str

	// CarriersFedexSpecificcountry => Ship to Specific Countries.
	// Path: carriers/fedex/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersFedexSpecificcountry model.StringCSV

	// CarriersFedexDebug => Debug.
	// Path: carriers/fedex/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexDebug model.Bool

	// CarriersFedexShowmethod => Show Method if Not Applicable.
	// Path: carriers/fedex/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersFedexShowmethod model.Bool

	// CarriersFedexSortOrder => Sort Order.
	// Path: carriers/fedex/sort_order
	CarriersFedexSortOrder model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFedexActive = model.NewBool(`carriers/fedex/active`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexActiveRma = model.NewBool(`carriers/fedex/active_rma`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexTitle = model.NewStr(`carriers/fedex/title`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexAccount = model.NewStr(`carriers/fedex/account`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexMeterNumber = model.NewStr(`carriers/fedex/meter_number`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexKey = model.NewStr(`carriers/fedex/key`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexPassword = model.NewStr(`carriers/fedex/password`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSandboxMode = model.NewBool(`carriers/fedex/sandbox_mode`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexProductionWebservicesUrl = model.NewStr(`carriers/fedex/production_webservices_url`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSandboxWebservicesUrl = model.NewStr(`carriers/fedex/sandbox_webservices_url`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexShipmentRequesttype = model.NewStr(`carriers/fedex/shipment_requesttype`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexPackaging = model.NewStr(`carriers/fedex/packaging`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexDropoff = model.NewStr(`carriers/fedex/dropoff`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexUnitOfMeasure = model.NewStr(`carriers/fedex/unit_of_measure`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexMaxPackageWeight = model.NewStr(`carriers/fedex/max_package_weight`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingType = model.NewStr(`carriers/fedex/handling_type`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingAction = model.NewStr(`carriers/fedex/handling_action`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexHandlingFee = model.NewStr(`carriers/fedex/handling_fee`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexResidenceDelivery = model.NewBool(`carriers/fedex/residence_delivery`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexAllowedMethods = model.NewStringCSV(`carriers/fedex/allowed_methods`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSmartpostHubid = model.NewStr(`carriers/fedex/smartpost_hubid`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeMethod = model.NewStr(`carriers/fedex/free_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeShippingEnable = model.NewBool(`carriers/fedex/free_shipping_enable`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexFreeShippingSubtotal = model.NewStr(`carriers/fedex/free_shipping_subtotal`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSpecificerrmsg = model.NewStr(`carriers/fedex/specificerrmsg`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSallowspecific = model.NewStr(`carriers/fedex/sallowspecific`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSpecificcountry = model.NewStringCSV(`carriers/fedex/specificcountry`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexDebug = model.NewBool(`carriers/fedex/debug`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexShowmethod = model.NewBool(`carriers/fedex/showmethod`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersFedexSortOrder = model.NewStr(`carriers/fedex/sort_order`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
