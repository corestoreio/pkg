// +build ignore

package fedex

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
	// CarriersFedexActive => Enabled for Checkout.
	// Path: carriers/fedex/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexActive model.Bool

	// CarriersFedexActiveRma => Enabled for RMA.
	// Path: carriers/fedex/active_rma
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexActiveRma model.Bool

	// CarriersFedexTitle => Title.
	// Path: carriers/fedex/title
	CarriersFedexTitle model.Str

	// CarriersFedexAccount => Account ID.
	// Please make sure to use only digits here. No dashes are allowed.
	// Path: carriers/fedex/account
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersFedexAccount model.Str

	// CarriersFedexMeterNumber => Meter Number.
	// Path: carriers/fedex/meter_number
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersFedexMeterNumber model.Str

	// CarriersFedexKey => Key.
	// Path: carriers/fedex/key
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersFedexKey model.Str

	// CarriersFedexPassword => Password.
	// Path: carriers/fedex/password
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersFedexPassword model.Str

	// CarriersFedexSandboxMode => Sandbox Mode.
	// Path: carriers/fedex/sandbox_mode
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexSandboxMode model.Bool

	// CarriersFedexProductionWebservicesUrl => Web-Services URL (Production).
	// Path: carriers/fedex/production_webservices_url
	CarriersFedexProductionWebservicesUrl model.Str

	// CarriersFedexSandboxWebservicesUrl => Web-Services URL (Sandbox).
	// Path: carriers/fedex/sandbox_webservices_url
	CarriersFedexSandboxWebservicesUrl model.Str

	// CarriersFedexShipmentRequesttype => Packages Request Type.
	// Path: carriers/fedex/shipment_requesttype
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
	CarriersFedexShipmentRequesttype model.Str

	// CarriersFedexPackaging => Packaging.
	// Path: carriers/fedex/packaging
	// SourceModel: Otnegam\Fedex\Model\Source\Packaging
	CarriersFedexPackaging model.Str

	// CarriersFedexDropoff => Dropoff.
	// Path: carriers/fedex/dropoff
	// SourceModel: Otnegam\Fedex\Model\Source\Dropoff
	CarriersFedexDropoff model.Str

	// CarriersFedexUnitOfMeasure => Weight Unit.
	// Path: carriers/fedex/unit_of_measure
	// SourceModel: Otnegam\Fedex\Model\Source\Unitofmeasure
	CarriersFedexUnitOfMeasure model.Str

	// CarriersFedexMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/fedex/max_package_weight
	CarriersFedexMaxPackageWeight model.Str

	// CarriersFedexHandlingType => Calculate Handling Fee.
	// Path: carriers/fedex/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
	CarriersFedexHandlingType model.Str

	// CarriersFedexHandlingAction => Handling Applied.
	// Path: carriers/fedex/handling_action
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
	CarriersFedexHandlingAction model.Str

	// CarriersFedexHandlingFee => Handling Fee.
	// Path: carriers/fedex/handling_fee
	CarriersFedexHandlingFee model.Str

	// CarriersFedexResidenceDelivery => Residential Delivery.
	// Path: carriers/fedex/residence_delivery
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexResidenceDelivery model.Bool

	// CarriersFedexAllowedMethods => Allowed Methods.
	// Path: carriers/fedex/allowed_methods
	// SourceModel: Otnegam\Fedex\Model\Source\Method
	CarriersFedexAllowedMethods model.StringCSV

	// CarriersFedexSmartpostHubid => Hub ID.
	// The field is applicable if the Smart Post method is selected.
	// Path: carriers/fedex/smartpost_hubid
	CarriersFedexSmartpostHubid model.Str

	// CarriersFedexFreeMethod => Free Method.
	// Path: carriers/fedex/free_method
	// SourceModel: Otnegam\Fedex\Model\Source\Freemethod
	CarriersFedexFreeMethod model.Str

	// CarriersFedexFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	CarriersFedexFreeShippingEnable model.Bool

	// CarriersFedexFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/fedex/free_shipping_subtotal
	CarriersFedexFreeShippingSubtotal model.Str

	// CarriersFedexSpecificerrmsg => Displayed Error Message.
	// Path: carriers/fedex/specificerrmsg
	CarriersFedexSpecificerrmsg model.Str

	// CarriersFedexSallowspecific => Ship to Applicable Countries.
	// Path: carriers/fedex/sallowspecific
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersFedexSallowspecific model.Str

	// CarriersFedexSpecificcountry => Ship to Specific Countries.
	// Path: carriers/fedex/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersFedexSpecificcountry model.StringCSV

	// CarriersFedexDebug => Debug.
	// Path: carriers/fedex/debug
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexDebug model.Bool

	// CarriersFedexShowmethod => Show Method if Not Applicable.
	// Path: carriers/fedex/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersFedexShowmethod model.Bool

	// CarriersFedexSortOrder => Sort Order.
	// Path: carriers/fedex/sort_order
	CarriersFedexSortOrder model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersFedexActive = model.NewBool(`carriers/fedex/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexActiveRma = model.NewBool(`carriers/fedex/active_rma`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexTitle = model.NewStr(`carriers/fedex/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexAccount = model.NewStr(`carriers/fedex/account`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexMeterNumber = model.NewStr(`carriers/fedex/meter_number`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexKey = model.NewStr(`carriers/fedex/key`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexPassword = model.NewStr(`carriers/fedex/password`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSandboxMode = model.NewBool(`carriers/fedex/sandbox_mode`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexProductionWebservicesUrl = model.NewStr(`carriers/fedex/production_webservices_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSandboxWebservicesUrl = model.NewStr(`carriers/fedex/sandbox_webservices_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexShipmentRequesttype = model.NewStr(`carriers/fedex/shipment_requesttype`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexPackaging = model.NewStr(`carriers/fedex/packaging`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexDropoff = model.NewStr(`carriers/fedex/dropoff`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexUnitOfMeasure = model.NewStr(`carriers/fedex/unit_of_measure`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexMaxPackageWeight = model.NewStr(`carriers/fedex/max_package_weight`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexHandlingType = model.NewStr(`carriers/fedex/handling_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexHandlingAction = model.NewStr(`carriers/fedex/handling_action`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexHandlingFee = model.NewStr(`carriers/fedex/handling_fee`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexResidenceDelivery = model.NewBool(`carriers/fedex/residence_delivery`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexAllowedMethods = model.NewStringCSV(`carriers/fedex/allowed_methods`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSmartpostHubid = model.NewStr(`carriers/fedex/smartpost_hubid`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexFreeMethod = model.NewStr(`carriers/fedex/free_method`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexFreeShippingEnable = model.NewBool(`carriers/fedex/free_shipping_enable`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexFreeShippingSubtotal = model.NewStr(`carriers/fedex/free_shipping_subtotal`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSpecificerrmsg = model.NewStr(`carriers/fedex/specificerrmsg`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSallowspecific = model.NewStr(`carriers/fedex/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSpecificcountry = model.NewStringCSV(`carriers/fedex/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexDebug = model.NewBool(`carriers/fedex/debug`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexShowmethod = model.NewBool(`carriers/fedex/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersFedexSortOrder = model.NewStr(`carriers/fedex/sort_order`, model.WithPkgCfg(pkgCfg))

	return pp
}
