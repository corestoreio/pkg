// +build ignore

package ups

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
	// CarriersUpsAccessLicenseNumber => Access License Number.
	// Path: carriers/ups/access_license_number
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersUpsAccessLicenseNumber model.Str

	// CarriersUpsActive => Enabled for Checkout.
	// Path: carriers/ups/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsActive model.Bool

	// CarriersUpsActiveRma => Enabled for RMA.
	// Path: carriers/ups/active_rma
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsActiveRma model.Bool

	// CarriersUpsAllowedMethods => Allowed Methods.
	// Path: carriers/ups/allowed_methods
	// SourceModel: Otnegam\Ups\Model\Config\Source\Method
	CarriersUpsAllowedMethods model.StringCSV

	// CarriersUpsShipmentRequesttype => Packages Request Type.
	// Path: carriers/ups/shipment_requesttype
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
	CarriersUpsShipmentRequesttype model.Str

	// CarriersUpsContainer => Container.
	// Path: carriers/ups/container
	// SourceModel: Otnegam\Ups\Model\Config\Source\Container
	CarriersUpsContainer model.Str

	// CarriersUpsFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	CarriersUpsFreeShippingEnable model.Bool

	// CarriersUpsFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_subtotal
	CarriersUpsFreeShippingSubtotal model.Str

	// CarriersUpsDestType => Destination Type.
	// Path: carriers/ups/dest_type
	// SourceModel: Otnegam\Ups\Model\Config\Source\DestType
	CarriersUpsDestType model.Str

	// CarriersUpsFreeMethod => Free Method.
	// Path: carriers/ups/free_method
	// SourceModel: Otnegam\Ups\Model\Config\Source\Freemethod
	CarriersUpsFreeMethod model.Str

	// CarriersUpsGatewayUrl => Gateway URL.
	// Path: carriers/ups/gateway_url
	CarriersUpsGatewayUrl model.Str

	// CarriersUpsGatewayXmlUrl => Gateway XML URL.
	// Path: carriers/ups/gateway_xml_url
	CarriersUpsGatewayXmlUrl model.Str

	// CarriersUpsHandlingType => Calculate Handling Fee.
	// Path: carriers/ups/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
	CarriersUpsHandlingType model.Str

	// CarriersUpsHandlingAction => Handling Applied.
	// Path: carriers/ups/handling_action
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
	CarriersUpsHandlingAction model.Str

	// CarriersUpsHandlingFee => Handling Fee.
	// Path: carriers/ups/handling_fee
	CarriersUpsHandlingFee model.Str

	// CarriersUpsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/ups/max_package_weight
	CarriersUpsMaxPackageWeight model.Str

	// CarriersUpsMinPackageWeight => Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight).
	// Path: carriers/ups/min_package_weight
	CarriersUpsMinPackageWeight model.Str

	// CarriersUpsOriginShipment => Origin of the Shipment.
	// Path: carriers/ups/origin_shipment
	// SourceModel: Otnegam\Ups\Model\Config\Source\OriginShipment
	CarriersUpsOriginShipment model.Str

	// CarriersUpsPassword => Password.
	// Path: carriers/ups/password
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersUpsPassword model.Str

	// CarriersUpsPickup => Pickup Method.
	// Path: carriers/ups/pickup
	// SourceModel: Otnegam\Ups\Model\Config\Source\Pickup
	CarriersUpsPickup model.Str

	// CarriersUpsSortOrder => Sort Order.
	// Path: carriers/ups/sort_order
	CarriersUpsSortOrder model.Str

	// CarriersUpsTitle => Title.
	// Path: carriers/ups/title
	CarriersUpsTitle model.Str

	// CarriersUpsTrackingXmlUrl => Tracking XML URL.
	// Path: carriers/ups/tracking_xml_url
	CarriersUpsTrackingXmlUrl model.Str

	// CarriersUpsType => UPS Type.
	// Path: carriers/ups/type
	// SourceModel: Otnegam\Ups\Model\Config\Source\Type
	CarriersUpsType model.Str

	// CarriersUpsIsAccountLive => Live Account.
	// Path: carriers/ups/is_account_live
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsIsAccountLive model.Bool

	// CarriersUpsUnitOfMeasure => Weight Unit.
	// Path: carriers/ups/unit_of_measure
	// SourceModel: Otnegam\Ups\Model\Config\Source\Unitofmeasure
	CarriersUpsUnitOfMeasure model.Str

	// CarriersUpsUsername => User ID.
	// Path: carriers/ups/username
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersUpsUsername model.Str

	// CarriersUpsNegotiatedActive => Enable Negotiated Rates.
	// Path: carriers/ups/negotiated_active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsNegotiatedActive model.Bool

	// CarriersUpsShipperNumber => Shipper Number.
	// Required for negotiated rates; 6-character UPS
	// Path: carriers/ups/shipper_number
	CarriersUpsShipperNumber model.Str

	// CarriersUpsSallowspecific => Ship to Applicable Countries.
	// Path: carriers/ups/sallowspecific
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersUpsSallowspecific model.Str

	// CarriersUpsSpecificcountry => Ship to Specific Countries.
	// Path: carriers/ups/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersUpsSpecificcountry model.StringCSV

	// CarriersUpsShowmethod => Show Method if Not Applicable.
	// Path: carriers/ups/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsShowmethod model.Bool

	// CarriersUpsSpecificerrmsg => Displayed Error Message.
	// Path: carriers/ups/specificerrmsg
	CarriersUpsSpecificerrmsg model.Str

	// CarriersUpsModeXml => Mode.
	// This enables or disables SSL verification of the Otnegam server by UPS.
	// Path: carriers/ups/mode_xml
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
	CarriersUpsModeXml model.Str

	// CarriersUpsDebug => Debug.
	// Path: carriers/ups/debug
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUpsDebug model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersUpsAccessLicenseNumber = model.NewStr(`carriers/ups/access_license_number`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsActive = model.NewBool(`carriers/ups/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsActiveRma = model.NewBool(`carriers/ups/active_rma`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsAllowedMethods = model.NewStringCSV(`carriers/ups/allowed_methods`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsShipmentRequesttype = model.NewStr(`carriers/ups/shipment_requesttype`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsContainer = model.NewStr(`carriers/ups/container`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsFreeShippingEnable = model.NewBool(`carriers/ups/free_shipping_enable`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsFreeShippingSubtotal = model.NewStr(`carriers/ups/free_shipping_subtotal`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsDestType = model.NewStr(`carriers/ups/dest_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsFreeMethod = model.NewStr(`carriers/ups/free_method`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsGatewayUrl = model.NewStr(`carriers/ups/gateway_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsGatewayXmlUrl = model.NewStr(`carriers/ups/gateway_xml_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsHandlingType = model.NewStr(`carriers/ups/handling_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsHandlingAction = model.NewStr(`carriers/ups/handling_action`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsHandlingFee = model.NewStr(`carriers/ups/handling_fee`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsMaxPackageWeight = model.NewStr(`carriers/ups/max_package_weight`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsMinPackageWeight = model.NewStr(`carriers/ups/min_package_weight`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsOriginShipment = model.NewStr(`carriers/ups/origin_shipment`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsPassword = model.NewStr(`carriers/ups/password`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsPickup = model.NewStr(`carriers/ups/pickup`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsSortOrder = model.NewStr(`carriers/ups/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsTitle = model.NewStr(`carriers/ups/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsTrackingXmlUrl = model.NewStr(`carriers/ups/tracking_xml_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsType = model.NewStr(`carriers/ups/type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsIsAccountLive = model.NewBool(`carriers/ups/is_account_live`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsUnitOfMeasure = model.NewStr(`carriers/ups/unit_of_measure`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsUsername = model.NewStr(`carriers/ups/username`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsNegotiatedActive = model.NewBool(`carriers/ups/negotiated_active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsShipperNumber = model.NewStr(`carriers/ups/shipper_number`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsSallowspecific = model.NewStr(`carriers/ups/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsSpecificcountry = model.NewStringCSV(`carriers/ups/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsShowmethod = model.NewBool(`carriers/ups/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsSpecificerrmsg = model.NewStr(`carriers/ups/specificerrmsg`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsModeXml = model.NewStr(`carriers/ups/mode_xml`, model.WithPkgCfg(pkgCfg))
	pp.CarriersUpsDebug = model.NewBool(`carriers/ups/debug`, model.WithPkgCfg(pkgCfg))

	return pp
}
