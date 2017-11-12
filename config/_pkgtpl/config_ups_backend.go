// +build ignore

package ups

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
	// CarriersUpsAccessLicenseNumber => Access License Number.
	// Path: carriers/ups/access_license_number
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsAccessLicenseNumber cfgmodel.Str

	// CarriersUpsActive => Enabled for Checkout.
	// Path: carriers/ups/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsActive cfgmodel.Bool

	// CarriersUpsActiveRma => Enabled for RMA.
	// Path: carriers/ups/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsActiveRma cfgmodel.Bool

	// CarriersUpsAllowedMethods => Allowed Methods.
	// Path: carriers/ups/allowed_methods
	// SourceModel: Magento\Ups\Model\Config\Source\Method
	CarriersUpsAllowedMethods cfgmodel.StringCSV

	// CarriersUpsShipmentRequesttype => Packages Request Type.
	// Path: carriers/ups/shipment_requesttype
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
	CarriersUpsShipmentRequesttype cfgmodel.Str

	// CarriersUpsContainer => Container.
	// Path: carriers/ups/container
	// SourceModel: Magento\Ups\Model\Config\Source\Container
	CarriersUpsContainer cfgmodel.Str

	// CarriersUpsFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersUpsFreeShippingEnable cfgmodel.Bool

	// CarriersUpsFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/ups/free_shipping_subtotal
	CarriersUpsFreeShippingSubtotal cfgmodel.Str

	// CarriersUpsDestType => Destination Type.
	// Path: carriers/ups/dest_type
	// SourceModel: Magento\Ups\Model\Config\Source\DestType
	CarriersUpsDestType cfgmodel.Str

	// CarriersUpsFreeMethod => Free Method.
	// Path: carriers/ups/free_method
	// SourceModel: Magento\Ups\Model\Config\Source\Freemethod
	CarriersUpsFreeMethod cfgmodel.Str

	// CarriersUpsGatewayUrl => Gateway URL.
	// Path: carriers/ups/gateway_url
	CarriersUpsGatewayUrl cfgmodel.Str

	// CarriersUpsGatewayXmlUrl => Gateway XML URL.
	// Path: carriers/ups/gateway_xml_url
	CarriersUpsGatewayXmlUrl cfgmodel.Str

	// CarriersUpsHandlingType => Calculate Handling Fee.
	// Path: carriers/ups/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersUpsHandlingType cfgmodel.Str

	// CarriersUpsHandlingAction => Handling Applied.
	// Path: carriers/ups/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
	CarriersUpsHandlingAction cfgmodel.Str

	// CarriersUpsHandlingFee => Handling Fee.
	// Path: carriers/ups/handling_fee
	CarriersUpsHandlingFee cfgmodel.Str

	// CarriersUpsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/ups/max_package_weight
	CarriersUpsMaxPackageWeight cfgmodel.Str

	// CarriersUpsMinPackageWeight => Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight).
	// Path: carriers/ups/min_package_weight
	CarriersUpsMinPackageWeight cfgmodel.Str

	// CarriersUpsOriginShipment => Origin of the Shipment.
	// Path: carriers/ups/origin_shipment
	// SourceModel: Magento\Ups\Model\Config\Source\OriginShipment
	CarriersUpsOriginShipment cfgmodel.Str

	// CarriersUpsPassword => Password.
	// Path: carriers/ups/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsPassword cfgmodel.Str

	// CarriersUpsPickup => Pickup Method.
	// Path: carriers/ups/pickup
	// SourceModel: Magento\Ups\Model\Config\Source\Pickup
	CarriersUpsPickup cfgmodel.Str

	// CarriersUpsSortOrder => Sort Order.
	// Path: carriers/ups/sort_order
	CarriersUpsSortOrder cfgmodel.Str

	// CarriersUpsTitle => Title.
	// Path: carriers/ups/title
	CarriersUpsTitle cfgmodel.Str

	// CarriersUpsTrackingXmlUrl => Tracking XML URL.
	// Path: carriers/ups/tracking_xml_url
	CarriersUpsTrackingXmlUrl cfgmodel.Str

	// CarriersUpsType => UPS Type.
	// Path: carriers/ups/type
	// SourceModel: Magento\Ups\Model\Config\Source\Type
	CarriersUpsType cfgmodel.Str

	// CarriersUpsIsAccountLive => Live Account.
	// Path: carriers/ups/is_account_live
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsIsAccountLive cfgmodel.Bool

	// CarriersUpsUnitOfMeasure => Weight Unit.
	// Path: carriers/ups/unit_of_measure
	// SourceModel: Magento\Ups\Model\Config\Source\Unitofmeasure
	CarriersUpsUnitOfMeasure cfgmodel.Str

	// CarriersUpsUsername => User ID.
	// Path: carriers/ups/username
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUpsUsername cfgmodel.Str

	// CarriersUpsNegotiatedActive => Enable Negotiated Rates.
	// Path: carriers/ups/negotiated_active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsNegotiatedActive cfgmodel.Bool

	// CarriersUpsShipperNumber => Shipper Number.
	// Required for negotiated rates; 6-character UPS
	// Path: carriers/ups/shipper_number
	CarriersUpsShipperNumber cfgmodel.Str

	// CarriersUpsSallowspecific => Ship to Applicable Countries.
	// Path: carriers/ups/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersUpsSallowspecific cfgmodel.Str

	// CarriersUpsSpecificcountry => Ship to Specific Countries.
	// Path: carriers/ups/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersUpsSpecificcountry cfgmodel.StringCSV

	// CarriersUpsShowmethod => Show Method if Not Applicable.
	// Path: carriers/ups/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsShowmethod cfgmodel.Bool

	// CarriersUpsSpecificerrmsg => Displayed Error Message.
	// Path: carriers/ups/specificerrmsg
	CarriersUpsSpecificerrmsg cfgmodel.Str

	// CarriersUpsModeXml => Mode.
	// This enables or disables SSL verification of the Magento server by UPS.
	// Path: carriers/ups/mode_xml
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Mode
	CarriersUpsModeXml cfgmodel.Str

	// CarriersUpsDebug => Debug.
	// Path: carriers/ups/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUpsDebug cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersUpsAccessLicenseNumber = cfgmodel.NewStr(`carriers/ups/access_license_number`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsActive = cfgmodel.NewBool(`carriers/ups/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsActiveRma = cfgmodel.NewBool(`carriers/ups/active_rma`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsAllowedMethods = cfgmodel.NewStringCSV(`carriers/ups/allowed_methods`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsShipmentRequesttype = cfgmodel.NewStr(`carriers/ups/shipment_requesttype`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsContainer = cfgmodel.NewStr(`carriers/ups/container`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsFreeShippingEnable = cfgmodel.NewBool(`carriers/ups/free_shipping_enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsFreeShippingSubtotal = cfgmodel.NewStr(`carriers/ups/free_shipping_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsDestType = cfgmodel.NewStr(`carriers/ups/dest_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsFreeMethod = cfgmodel.NewStr(`carriers/ups/free_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsGatewayUrl = cfgmodel.NewStr(`carriers/ups/gateway_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsGatewayXmlUrl = cfgmodel.NewStr(`carriers/ups/gateway_xml_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsHandlingType = cfgmodel.NewStr(`carriers/ups/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsHandlingAction = cfgmodel.NewStr(`carriers/ups/handling_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsHandlingFee = cfgmodel.NewStr(`carriers/ups/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsMaxPackageWeight = cfgmodel.NewStr(`carriers/ups/max_package_weight`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsMinPackageWeight = cfgmodel.NewStr(`carriers/ups/min_package_weight`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsOriginShipment = cfgmodel.NewStr(`carriers/ups/origin_shipment`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsPassword = cfgmodel.NewStr(`carriers/ups/password`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsPickup = cfgmodel.NewStr(`carriers/ups/pickup`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsSortOrder = cfgmodel.NewStr(`carriers/ups/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsTitle = cfgmodel.NewStr(`carriers/ups/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsTrackingXmlUrl = cfgmodel.NewStr(`carriers/ups/tracking_xml_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsType = cfgmodel.NewStr(`carriers/ups/type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsIsAccountLive = cfgmodel.NewBool(`carriers/ups/is_account_live`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsUnitOfMeasure = cfgmodel.NewStr(`carriers/ups/unit_of_measure`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsUsername = cfgmodel.NewStr(`carriers/ups/username`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsNegotiatedActive = cfgmodel.NewBool(`carriers/ups/negotiated_active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsShipperNumber = cfgmodel.NewStr(`carriers/ups/shipper_number`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsSallowspecific = cfgmodel.NewStr(`carriers/ups/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsSpecificcountry = cfgmodel.NewStringCSV(`carriers/ups/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsShowmethod = cfgmodel.NewBool(`carriers/ups/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsSpecificerrmsg = cfgmodel.NewStr(`carriers/ups/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsModeXml = cfgmodel.NewStr(`carriers/ups/mode_xml`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUpsDebug = cfgmodel.NewBool(`carriers/ups/debug`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
