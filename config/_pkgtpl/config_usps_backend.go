// +build ignore

package usps

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
	// CarriersUspsActive => Enabled for Checkout.
	// Path: carriers/usps/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUspsActive cfgmodel.Bool

	// CarriersUspsActiveRma => Enabled for RMA.
	// Path: carriers/usps/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUspsActiveRma cfgmodel.Bool

	// CarriersUspsGatewayUrl => Gateway URL.
	// Path: carriers/usps/gateway_url
	CarriersUspsGatewayUrl cfgmodel.Str

	// CarriersUspsGatewaySecureUrl => Secure Gateway URL.
	// Path: carriers/usps/gateway_secure_url
	CarriersUspsGatewaySecureUrl cfgmodel.Str

	// CarriersUspsTitle => Title.
	// Path: carriers/usps/title
	CarriersUspsTitle cfgmodel.Str

	// CarriersUspsUserid => User ID.
	// Path: carriers/usps/userid
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUspsUserid cfgmodel.Str

	// CarriersUspsPassword => Password.
	// Path: carriers/usps/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersUspsPassword cfgmodel.Str

	// CarriersUspsMode => Mode.
	// Path: carriers/usps/mode
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Mode
	CarriersUspsMode cfgmodel.Str

	// CarriersUspsShipmentRequesttype => Packages Request Type.
	// Path: carriers/usps/shipment_requesttype
	// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
	CarriersUspsShipmentRequesttype cfgmodel.Str

	// CarriersUspsContainer => Container.
	// Path: carriers/usps/container
	// SourceModel: Magento\Usps\Model\Source\Container
	CarriersUspsContainer cfgmodel.Str

	// CarriersUspsSize => Size.
	// Path: carriers/usps/size
	// SourceModel: Magento\Usps\Model\Source\Size
	CarriersUspsSize cfgmodel.Str

	// CarriersUspsWidth => Width.
	// Path: carriers/usps/width
	CarriersUspsWidth cfgmodel.Str

	// CarriersUspsLength => Length.
	// Path: carriers/usps/length
	CarriersUspsLength cfgmodel.Str

	// CarriersUspsHeight => Height.
	// Path: carriers/usps/height
	CarriersUspsHeight cfgmodel.Str

	// CarriersUspsGirth => Girth.
	// Path: carriers/usps/girth
	CarriersUspsGirth cfgmodel.Str

	// CarriersUspsMachinable => Machinable.
	// Path: carriers/usps/machinable
	// SourceModel: Magento\Usps\Model\Source\Machinable
	CarriersUspsMachinable cfgmodel.Str

	// CarriersUspsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/usps/max_package_weight
	CarriersUspsMaxPackageWeight cfgmodel.Str

	// CarriersUspsHandlingType => Calculate Handling Fee.
	// Path: carriers/usps/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersUspsHandlingType cfgmodel.Str

	// CarriersUspsHandlingAction => Handling Applied.
	// Path: carriers/usps/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
	CarriersUspsHandlingAction cfgmodel.Str

	// CarriersUspsHandlingFee => Handling Fee.
	// Path: carriers/usps/handling_fee
	CarriersUspsHandlingFee cfgmodel.Str

	// CarriersUspsAllowedMethods => Allowed Methods.
	// Path: carriers/usps/allowed_methods
	// SourceModel: Magento\Usps\Model\Source\Method
	CarriersUspsAllowedMethods cfgmodel.StringCSV

	// CarriersUspsFreeMethod => Free Method.
	// Path: carriers/usps/free_method
	// SourceModel: Magento\Usps\Model\Source\Freemethod
	CarriersUspsFreeMethod cfgmodel.Str

	// CarriersUspsFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/usps/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersUspsFreeShippingEnable cfgmodel.Bool

	// CarriersUspsFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/usps/free_shipping_subtotal
	CarriersUspsFreeShippingSubtotal cfgmodel.Str

	// CarriersUspsSpecificerrmsg => Displayed Error Message.
	// Path: carriers/usps/specificerrmsg
	CarriersUspsSpecificerrmsg cfgmodel.Str

	// CarriersUspsSallowspecific => Ship to Applicable Countries.
	// Path: carriers/usps/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersUspsSallowspecific cfgmodel.Str

	// CarriersUspsSpecificcountry => Ship to Specific Countries.
	// Path: carriers/usps/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersUspsSpecificcountry cfgmodel.StringCSV

	// CarriersUspsDebug => Debug.
	// Path: carriers/usps/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUspsDebug cfgmodel.Bool

	// CarriersUspsShowmethod => Show Method if Not Applicable.
	// Path: carriers/usps/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersUspsShowmethod cfgmodel.Bool

	// CarriersUspsSortOrder => Sort Order.
	// Path: carriers/usps/sort_order
	CarriersUspsSortOrder cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersUspsActive = cfgmodel.NewBool(`carriers/usps/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsActiveRma = cfgmodel.NewBool(`carriers/usps/active_rma`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsGatewayUrl = cfgmodel.NewStr(`carriers/usps/gateway_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsGatewaySecureUrl = cfgmodel.NewStr(`carriers/usps/gateway_secure_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsTitle = cfgmodel.NewStr(`carriers/usps/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsUserid = cfgmodel.NewStr(`carriers/usps/userid`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsPassword = cfgmodel.NewStr(`carriers/usps/password`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsMode = cfgmodel.NewStr(`carriers/usps/mode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsShipmentRequesttype = cfgmodel.NewStr(`carriers/usps/shipment_requesttype`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsContainer = cfgmodel.NewStr(`carriers/usps/container`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsSize = cfgmodel.NewStr(`carriers/usps/size`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsWidth = cfgmodel.NewStr(`carriers/usps/width`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsLength = cfgmodel.NewStr(`carriers/usps/length`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsHeight = cfgmodel.NewStr(`carriers/usps/height`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsGirth = cfgmodel.NewStr(`carriers/usps/girth`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsMachinable = cfgmodel.NewStr(`carriers/usps/machinable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsMaxPackageWeight = cfgmodel.NewStr(`carriers/usps/max_package_weight`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsHandlingType = cfgmodel.NewStr(`carriers/usps/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsHandlingAction = cfgmodel.NewStr(`carriers/usps/handling_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsHandlingFee = cfgmodel.NewStr(`carriers/usps/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsAllowedMethods = cfgmodel.NewStringCSV(`carriers/usps/allowed_methods`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsFreeMethod = cfgmodel.NewStr(`carriers/usps/free_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsFreeShippingEnable = cfgmodel.NewBool(`carriers/usps/free_shipping_enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsFreeShippingSubtotal = cfgmodel.NewStr(`carriers/usps/free_shipping_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsSpecificerrmsg = cfgmodel.NewStr(`carriers/usps/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsSallowspecific = cfgmodel.NewStr(`carriers/usps/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsSpecificcountry = cfgmodel.NewStringCSV(`carriers/usps/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsDebug = cfgmodel.NewBool(`carriers/usps/debug`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsShowmethod = cfgmodel.NewBool(`carriers/usps/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersUspsSortOrder = cfgmodel.NewStr(`carriers/usps/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
