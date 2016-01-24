// +build ignore

package ups

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
	// CarriersUpsAccessLicenseNumber => Access License Number.
	// Path: carriers/ups/access_license_number
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsAccessLicenseNumber model.Str

	// CarriersUpsActive => Enabled for Checkout.
	// Path: carriers/ups/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsActive model.Bool

	// CarriersUpsActiveRma => Enabled for RMA.
	// Path: carriers/ups/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsActiveRma model.Bool

	// CarriersUpsAllowedMethods => Allowed Methods.
	// Path: carriers/ups/allowed_methods
	// SourceModel: Magento\Ups\Model\Config\Source\Method
	CarriersUpsAllowedMethods model.StringCSV

	// CarriersUpsShipmentRequesttype => Packages Request Type.
	// Path: carriers/ups/shipment_requesttype
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
	CarriersUpsShipmentRequesttype model.Str

	// CarriersUpsContainer => Container.
	// Path: carriers/ups/container
	// SourceModel: Magento\Ups\Model\Config\Source\Container
	CarriersUpsContainer model.Str

	// CarriersUpsFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersUpsFreeShippingEnable model.Bool

	// CarriersUpsFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_subtotal
	CarriersUpsFreeShippingSubtotal model.Str

	// CarriersUpsDestType => Destination Type.
	// Path: carriers/ups/dest_type
	// SourceModel: Magento\Ups\Model\Config\Source\DestType
	CarriersUpsDestType model.Str

	// CarriersUpsFreeMethod => Free Method.
	// Path: carriers/ups/free_method
	// SourceModel: Magento\Ups\Model\Config\Source\Freemethod
	CarriersUpsFreeMethod model.Str

	// CarriersUpsGatewayUrl => Gateway URL.
	// Path: carriers/ups/gateway_url
	CarriersUpsGatewayUrl model.Str

	// CarriersUpsGatewayXmlUrl => Gateway XML URL.
	// Path: carriers/ups/gateway_xml_url
	CarriersUpsGatewayXmlUrl model.Str

	// CarriersUpsHandlingType => Calculate Handling Fee.
	// Path: carriers/ups/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersUpsHandlingType model.Str

	// CarriersUpsHandlingAction => Handling Applied.
	// Path: carriers/ups/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
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
	// SourceModel: Magento\Ups\Model\Config\Source\OriginShipment
	CarriersUpsOriginShipment model.Str

	// CarriersUpsPassword => Password.
	// Path: carriers/ups/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsPassword model.Str

	// CarriersUpsPickup => Pickup Method.
	// Path: carriers/ups/pickup
	// SourceModel: Magento\Ups\Model\Config\Source\Pickup
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
	// SourceModel: Magento\Ups\Model\Config\Source\Type
	CarriersUpsType model.Str

	// CarriersUpsIsAccountLive => Live Account.
	// Path: carriers/ups/is_account_live
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsIsAccountLive model.Bool

	// CarriersUpsUnitOfMeasure => Weight Unit.
	// Path: carriers/ups/unit_of_measure
	// SourceModel: Magento\Ups\Model\Config\Source\Unitofmeasure
	CarriersUpsUnitOfMeasure model.Str

	// CarriersUpsUsername => User ID.
	// Path: carriers/ups/username
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsUsername model.Str

	// CarriersUpsNegotiatedActive => Enable Negotiated Rates.
	// Path: carriers/ups/negotiated_active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsNegotiatedActive model.Bool

	// CarriersUpsShipperNumber => Shipper Number.
	// Required for negotiated rates; 6-character UPS
	// Path: carriers/ups/shipper_number
	CarriersUpsShipperNumber model.Str

	// CarriersUpsSallowspecific => Ship to Applicable Countries.
	// Path: carriers/ups/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersUpsSallowspecific model.Str

	// CarriersUpsSpecificcountry => Ship to Specific Countries.
	// Path: carriers/ups/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersUpsSpecificcountry model.StringCSV

	// CarriersUpsShowmethod => Show Method if Not Applicable.
	// Path: carriers/ups/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsShowmethod model.Bool

	// CarriersUpsSpecificerrmsg => Displayed Error Message.
	// Path: carriers/ups/specificerrmsg
	CarriersUpsSpecificerrmsg model.Str

	// CarriersUpsModeXml => Mode.
	// This enables or disables SSL verification of the Magento server by UPS.
	// Path: carriers/ups/mode_xml
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Mode
	CarriersUpsModeXml model.Str

	// CarriersUpsDebug => Debug.
	// Path: carriers/ups/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsDebug model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersUpsAccessLicenseNumber = model.NewStr(`carriers/ups/access_license_number`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsActive = model.NewBool(`carriers/ups/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsActiveRma = model.NewBool(`carriers/ups/active_rma`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsAllowedMethods = model.NewStringCSV(`carriers/ups/allowed_methods`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsShipmentRequesttype = model.NewStr(`carriers/ups/shipment_requesttype`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsContainer = model.NewStr(`carriers/ups/container`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsFreeShippingEnable = model.NewBool(`carriers/ups/free_shipping_enable`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsFreeShippingSubtotal = model.NewStr(`carriers/ups/free_shipping_subtotal`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsDestType = model.NewStr(`carriers/ups/dest_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsFreeMethod = model.NewStr(`carriers/ups/free_method`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsGatewayUrl = model.NewStr(`carriers/ups/gateway_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsGatewayXmlUrl = model.NewStr(`carriers/ups/gateway_xml_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsHandlingType = model.NewStr(`carriers/ups/handling_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsHandlingAction = model.NewStr(`carriers/ups/handling_action`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsHandlingFee = model.NewStr(`carriers/ups/handling_fee`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsMaxPackageWeight = model.NewStr(`carriers/ups/max_package_weight`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsMinPackageWeight = model.NewStr(`carriers/ups/min_package_weight`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsOriginShipment = model.NewStr(`carriers/ups/origin_shipment`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsPassword = model.NewStr(`carriers/ups/password`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsPickup = model.NewStr(`carriers/ups/pickup`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsSortOrder = model.NewStr(`carriers/ups/sort_order`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsTitle = model.NewStr(`carriers/ups/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsTrackingXmlUrl = model.NewStr(`carriers/ups/tracking_xml_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsType = model.NewStr(`carriers/ups/type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsIsAccountLive = model.NewBool(`carriers/ups/is_account_live`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsUnitOfMeasure = model.NewStr(`carriers/ups/unit_of_measure`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsUsername = model.NewStr(`carriers/ups/username`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsNegotiatedActive = model.NewBool(`carriers/ups/negotiated_active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsShipperNumber = model.NewStr(`carriers/ups/shipper_number`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsSallowspecific = model.NewStr(`carriers/ups/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsSpecificcountry = model.NewStringCSV(`carriers/ups/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsShowmethod = model.NewBool(`carriers/ups/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsSpecificerrmsg = model.NewStr(`carriers/ups/specificerrmsg`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsModeXml = model.NewStr(`carriers/ups/mode_xml`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUpsDebug = model.NewBool(`carriers/ups/debug`, model.WithConfigStructure(cfgStruct))

	return pp
}
