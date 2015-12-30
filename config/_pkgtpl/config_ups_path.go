// +build ignore

package ups

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersUpsAccessLicenseNumber => Access License Number.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUpsAccessLicenseNumber = model.NewStr(`carriers/ups/access_license_number`)

// PathCarriersUpsActive => Enabled for Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsActive = model.NewBool(`carriers/ups/active`)

// PathCarriersUpsActiveRma => Enabled for RMA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsActiveRma = model.NewBool(`carriers/ups/active_rma`)

// PathCarriersUpsAllowedMethods => Allowed Methods.
// SourceModel: Otnegam\Ups\Model\Config\Source\Method
var PathCarriersUpsAllowedMethods = model.NewStringCSV(`carriers/ups/allowed_methods`)

// PathCarriersUpsShipmentRequesttype => Packages Request Type.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
var PathCarriersUpsShipmentRequesttype = model.NewStr(`carriers/ups/shipment_requesttype`)

// PathCarriersUpsContainer => Container.
// SourceModel: Otnegam\Ups\Model\Config\Source\Container
var PathCarriersUpsContainer = model.NewStr(`carriers/ups/container`)

// PathCarriersUpsFreeShippingEnable => Free Shipping Amount Threshold.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCarriersUpsFreeShippingEnable = model.NewBool(`carriers/ups/free_shipping_enable`)

// PathCarriersUpsFreeShippingSubtotal => Free Shipping Amount Threshold.
var PathCarriersUpsFreeShippingSubtotal = model.NewStr(`carriers/ups/free_shipping_subtotal`)

// PathCarriersUpsDestType => Destination Type.
// SourceModel: Otnegam\Ups\Model\Config\Source\DestType
var PathCarriersUpsDestType = model.NewStr(`carriers/ups/dest_type`)

// PathCarriersUpsFreeMethod => Free Method.
// SourceModel: Otnegam\Ups\Model\Config\Source\Freemethod
var PathCarriersUpsFreeMethod = model.NewStr(`carriers/ups/free_method`)

// PathCarriersUpsGatewayUrl => Gateway URL.
var PathCarriersUpsGatewayUrl = model.NewStr(`carriers/ups/gateway_url`)

// PathCarriersUpsGatewayXmlUrl => Gateway XML URL.
var PathCarriersUpsGatewayXmlUrl = model.NewStr(`carriers/ups/gateway_xml_url`)

// PathCarriersUpsHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersUpsHandlingType = model.NewStr(`carriers/ups/handling_type`)

// PathCarriersUpsHandlingAction => Handling Applied.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
var PathCarriersUpsHandlingAction = model.NewStr(`carriers/ups/handling_action`)

// PathCarriersUpsHandlingFee => Handling Fee.
var PathCarriersUpsHandlingFee = model.NewStr(`carriers/ups/handling_fee`)

// PathCarriersUpsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
var PathCarriersUpsMaxPackageWeight = model.NewStr(`carriers/ups/max_package_weight`)

// PathCarriersUpsMinPackageWeight => Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight).
var PathCarriersUpsMinPackageWeight = model.NewStr(`carriers/ups/min_package_weight`)

// PathCarriersUpsOriginShipment => Origin of the Shipment.
// SourceModel: Otnegam\Ups\Model\Config\Source\OriginShipment
var PathCarriersUpsOriginShipment = model.NewStr(`carriers/ups/origin_shipment`)

// PathCarriersUpsPassword => Password.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUpsPassword = model.NewStr(`carriers/ups/password`)

// PathCarriersUpsPickup => Pickup Method.
// SourceModel: Otnegam\Ups\Model\Config\Source\Pickup
var PathCarriersUpsPickup = model.NewStr(`carriers/ups/pickup`)

// PathCarriersUpsSortOrder => Sort Order.
var PathCarriersUpsSortOrder = model.NewStr(`carriers/ups/sort_order`)

// PathCarriersUpsTitle => Title.
var PathCarriersUpsTitle = model.NewStr(`carriers/ups/title`)

// PathCarriersUpsTrackingXmlUrl => Tracking XML URL.
var PathCarriersUpsTrackingXmlUrl = model.NewStr(`carriers/ups/tracking_xml_url`)

// PathCarriersUpsType => UPS Type.
// SourceModel: Otnegam\Ups\Model\Config\Source\Type
var PathCarriersUpsType = model.NewStr(`carriers/ups/type`)

// PathCarriersUpsIsAccountLive => Live Account.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsIsAccountLive = model.NewBool(`carriers/ups/is_account_live`)

// PathCarriersUpsUnitOfMeasure => Weight Unit.
// SourceModel: Otnegam\Ups\Model\Config\Source\Unitofmeasure
var PathCarriersUpsUnitOfMeasure = model.NewStr(`carriers/ups/unit_of_measure`)

// PathCarriersUpsUsername => User ID.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersUpsUsername = model.NewStr(`carriers/ups/username`)

// PathCarriersUpsNegotiatedActive => Enable Negotiated Rates.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsNegotiatedActive = model.NewBool(`carriers/ups/negotiated_active`)

// PathCarriersUpsShipperNumber => Shipper Number.
// Required for negotiated rates; 6-character UPS
var PathCarriersUpsShipperNumber = model.NewStr(`carriers/ups/shipper_number`)

// PathCarriersUpsSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersUpsSallowspecific = model.NewStr(`carriers/ups/sallowspecific`)

// PathCarriersUpsSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersUpsSpecificcountry = model.NewStringCSV(`carriers/ups/specificcountry`)

// PathCarriersUpsShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsShowmethod = model.NewBool(`carriers/ups/showmethod`)

// PathCarriersUpsSpecificerrmsg => Displayed Error Message.
var PathCarriersUpsSpecificerrmsg = model.NewStr(`carriers/ups/specificerrmsg`)

// PathCarriersUpsModeXml => Mode.
// This enables or disables SSL verification of the Otnegam server by UPS.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
var PathCarriersUpsModeXml = model.NewStr(`carriers/ups/mode_xml`)

// PathCarriersUpsDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersUpsDebug = model.NewBool(`carriers/ups/debug`)
