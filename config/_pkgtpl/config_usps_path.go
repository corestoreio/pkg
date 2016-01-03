// +build ignore

package usps

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CarriersUspsActive => Enabled for Checkout.
	// Path: carriers/usps/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUspsActive model.Bool

	// CarriersUspsActiveRma => Enabled for RMA.
	// Path: carriers/usps/active_rma
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUspsActiveRma model.Bool

	// CarriersUspsGatewayUrl => Gateway URL.
	// Path: carriers/usps/gateway_url
	CarriersUspsGatewayUrl model.Str

	// CarriersUspsGatewaySecureUrl => Secure Gateway URL.
	// Path: carriers/usps/gateway_secure_url
	CarriersUspsGatewaySecureUrl model.Str

	// CarriersUspsTitle => Title.
	// Path: carriers/usps/title
	CarriersUspsTitle model.Str

	// CarriersUspsUserid => User ID.
	// Path: carriers/usps/userid
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersUspsUserid model.Str

	// CarriersUspsPassword => Password.
	// Path: carriers/usps/password
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersUspsPassword model.Str

	// CarriersUspsMode => Mode.
	// Path: carriers/usps/mode
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
	CarriersUspsMode model.Str

	// CarriersUspsShipmentRequesttype => Packages Request Type.
	// Path: carriers/usps/shipment_requesttype
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
	CarriersUspsShipmentRequesttype model.Str

	// CarriersUspsContainer => Container.
	// Path: carriers/usps/container
	// SourceModel: Otnegam\Usps\Model\Source\Container
	CarriersUspsContainer model.Str

	// CarriersUspsSize => Size.
	// Path: carriers/usps/size
	// SourceModel: Otnegam\Usps\Model\Source\Size
	CarriersUspsSize model.Str

	// CarriersUspsWidth => Width.
	// Path: carriers/usps/width
	CarriersUspsWidth model.Str

	// CarriersUspsLength => Length.
	// Path: carriers/usps/length
	CarriersUspsLength model.Str

	// CarriersUspsHeight => Height.
	// Path: carriers/usps/height
	CarriersUspsHeight model.Str

	// CarriersUspsGirth => Girth.
	// Path: carriers/usps/girth
	CarriersUspsGirth model.Str

	// CarriersUspsMachinable => Machinable.
	// Path: carriers/usps/machinable
	// SourceModel: Otnegam\Usps\Model\Source\Machinable
	CarriersUspsMachinable model.Str

	// CarriersUspsMaxPackageWeight => Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight).
	// Path: carriers/usps/max_package_weight
	CarriersUspsMaxPackageWeight model.Str

	// CarriersUspsHandlingType => Calculate Handling Fee.
	// Path: carriers/usps/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
	CarriersUspsHandlingType model.Str

	// CarriersUspsHandlingAction => Handling Applied.
	// Path: carriers/usps/handling_action
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
	CarriersUspsHandlingAction model.Str

	// CarriersUspsHandlingFee => Handling Fee.
	// Path: carriers/usps/handling_fee
	CarriersUspsHandlingFee model.Str

	// CarriersUspsAllowedMethods => Allowed Methods.
	// Path: carriers/usps/allowed_methods
	// SourceModel: Otnegam\Usps\Model\Source\Method
	CarriersUspsAllowedMethods model.StringCSV

	// CarriersUspsFreeMethod => Free Method.
	// Path: carriers/usps/free_method
	// SourceModel: Otnegam\Usps\Model\Source\Freemethod
	CarriersUspsFreeMethod model.Str

	// CarriersUspsFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/usps/free_shipping_enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	CarriersUspsFreeShippingEnable model.Bool

	// CarriersUspsFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/usps/free_shipping_subtotal
	CarriersUspsFreeShippingSubtotal model.Str

	// CarriersUspsSpecificerrmsg => Displayed Error Message.
	// Path: carriers/usps/specificerrmsg
	CarriersUspsSpecificerrmsg model.Str

	// CarriersUspsSallowspecific => Ship to Applicable Countries.
	// Path: carriers/usps/sallowspecific
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersUspsSallowspecific model.Str

	// CarriersUspsSpecificcountry => Ship to Specific Countries.
	// Path: carriers/usps/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersUspsSpecificcountry model.StringCSV

	// CarriersUspsDebug => Debug.
	// Path: carriers/usps/debug
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUspsDebug model.Bool

	// CarriersUspsShowmethod => Show Method if Not Applicable.
	// Path: carriers/usps/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersUspsShowmethod model.Bool

	// CarriersUspsSortOrder => Sort Order.
	// Path: carriers/usps/sort_order
	CarriersUspsSortOrder model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersUspsActive = model.NewBool(`carriers/usps/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsActiveRma = model.NewBool(`carriers/usps/active_rma`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsGatewayUrl = model.NewStr(`carriers/usps/gateway_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsGatewaySecureUrl = model.NewStr(`carriers/usps/gateway_secure_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsTitle = model.NewStr(`carriers/usps/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsUserid = model.NewStr(`carriers/usps/userid`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsPassword = model.NewStr(`carriers/usps/password`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsMode = model.NewStr(`carriers/usps/mode`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsShipmentRequesttype = model.NewStr(`carriers/usps/shipment_requesttype`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsContainer = model.NewStr(`carriers/usps/container`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsSize = model.NewStr(`carriers/usps/size`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsWidth = model.NewStr(`carriers/usps/width`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsLength = model.NewStr(`carriers/usps/length`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsHeight = model.NewStr(`carriers/usps/height`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsGirth = model.NewStr(`carriers/usps/girth`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsMachinable = model.NewStr(`carriers/usps/machinable`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsMaxPackageWeight = model.NewStr(`carriers/usps/max_package_weight`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsHandlingType = model.NewStr(`carriers/usps/handling_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsHandlingAction = model.NewStr(`carriers/usps/handling_action`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsHandlingFee = model.NewStr(`carriers/usps/handling_fee`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsAllowedMethods = model.NewStringCSV(`carriers/usps/allowed_methods`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsFreeMethod = model.NewStr(`carriers/usps/free_method`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsFreeShippingEnable = model.NewBool(`carriers/usps/free_shipping_enable`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsFreeShippingSubtotal = model.NewStr(`carriers/usps/free_shipping_subtotal`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsSpecificerrmsg = model.NewStr(`carriers/usps/specificerrmsg`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsSallowspecific = model.NewStr(`carriers/usps/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsSpecificcountry = model.NewStringCSV(`carriers/usps/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsDebug = model.NewBool(`carriers/usps/debug`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsShowmethod = model.NewBool(`carriers/usps/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersUspsSortOrder = model.NewStr(`carriers/usps/sort_order`, model.WithConfigStructure(cfgStruct))

	return pp
}
