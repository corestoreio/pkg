// +build ignore

package dhl

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
	// CarriersDhlActive => Enabled for Checkout.
	// Path: carriers/dhl/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersDhlActive model.Bool

	// CarriersDhlActiveRma => Enabled for RMA.
	// Path: carriers/dhl/active_rma
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersDhlActiveRma model.Bool

	// CarriersDhlGatewayUrl => Gateway URL.
	// Path: carriers/dhl/gateway_url
	CarriersDhlGatewayUrl model.Str

	// CarriersDhlTitle => Title.
	// Path: carriers/dhl/title
	CarriersDhlTitle model.Str

	// CarriersDhlId => Access ID.
	// Path: carriers/dhl/id
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersDhlId model.Str

	// CarriersDhlPassword => Password.
	// Path: carriers/dhl/password
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	CarriersDhlPassword model.Str

	// CarriersDhlAccount => Account Number.
	// Path: carriers/dhl/account
	CarriersDhlAccount model.Str

	// CarriersDhlContentType => Content Type.
	// Path: carriers/dhl/content_type
	// SourceModel: Otnegam\Dhl\Model\Source\Contenttype
	CarriersDhlContentType model.Str

	// CarriersDhlHandlingType => Calculate Handling Fee.
	// Path: carriers/dhl/handling_type
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
	CarriersDhlHandlingType model.Str

	// CarriersDhlHandlingAction => Handling Applied.
	// "Per Order" allows a single handling fee for the entire order. "Per
	// Package" allows an individual handling fee for each package.
	// Path: carriers/dhl/handling_action
	// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
	CarriersDhlHandlingAction model.Str

	// CarriersDhlHandlingFee => Handling Fee.
	// Path: carriers/dhl/handling_fee
	CarriersDhlHandlingFee model.Str

	// CarriersDhlDivideOrderWeight => Divide Order Weight.
	// Select this to allow DHL to optimize shipping charges by splitting the
	// order if it exceeds 70 kg.
	// Path: carriers/dhl/divide_order_weight
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersDhlDivideOrderWeight model.Bool

	// CarriersDhlUnitOfMeasure => Weight Unit.
	// Path: carriers/dhl/unit_of_measure
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Unitofmeasure
	CarriersDhlUnitOfMeasure model.Str

	// CarriersDhlSize => Size.
	// Path: carriers/dhl/size
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Size
	CarriersDhlSize model.Str

	// CarriersDhlHeight => Height.
	// Path: carriers/dhl/height
	CarriersDhlHeight model.Str

	// CarriersDhlDepth => Depth.
	// Path: carriers/dhl/depth
	CarriersDhlDepth model.Str

	// CarriersDhlWidth => Width.
	// Path: carriers/dhl/width
	CarriersDhlWidth model.Str

	// CarriersDhlDocMethods => Allowed Methods.
	// Path: carriers/dhl/doc_methods
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Doc
	CarriersDhlDocMethods model.StringCSV

	// CarriersDhlNondocMethods => Allowed Methods.
	// Path: carriers/dhl/nondoc_methods
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Nondoc
	CarriersDhlNondocMethods model.StringCSV

	// CarriersDhlReadyTime => Ready time.
	// Package ready time after order submission (in hours)
	// Path: carriers/dhl/ready_time
	CarriersDhlReadyTime model.Str

	// CarriersDhlSpecificerrmsg => Displayed Error Message.
	// Path: carriers/dhl/specificerrmsg
	CarriersDhlSpecificerrmsg model.Str

	// CarriersDhlFreeMethodDoc => Free Method.
	// Path: carriers/dhl/free_method_doc
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Freedoc
	CarriersDhlFreeMethodDoc model.Str

	// CarriersDhlFreeMethodNondoc => Free Method.
	// Path: carriers/dhl/free_method_nondoc
	// SourceModel: Otnegam\Dhl\Model\Source\Method\Freenondoc
	CarriersDhlFreeMethodNondoc model.Str

	// CarriersDhlFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/dhl/free_shipping_enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	CarriersDhlFreeShippingEnable model.Bool

	// CarriersDhlFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/dhl/free_shipping_subtotal
	CarriersDhlFreeShippingSubtotal model.Str

	// CarriersDhlSallowspecific => Ship to Applicable Countries.
	// Path: carriers/dhl/sallowspecific
	// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
	CarriersDhlSallowspecific model.Str

	// CarriersDhlSpecificcountry => Ship to Specific Countries.
	// Path: carriers/dhl/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	CarriersDhlSpecificcountry model.StringCSV

	// CarriersDhlShowmethod => Show Method if Not Applicable.
	// Path: carriers/dhl/showmethod
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersDhlShowmethod model.Bool

	// CarriersDhlSortOrder => Sort Order.
	// Path: carriers/dhl/sort_order
	CarriersDhlSortOrder model.Str

	// CarriersDhlDebug => Debug.
	// Path: carriers/dhl/debug
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CarriersDhlDebug model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersDhlActive = model.NewBool(`carriers/dhl/active`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlActiveRma = model.NewBool(`carriers/dhl/active_rma`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlGatewayUrl = model.NewStr(`carriers/dhl/gateway_url`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlTitle = model.NewStr(`carriers/dhl/title`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlId = model.NewStr(`carriers/dhl/id`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlPassword = model.NewStr(`carriers/dhl/password`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlAccount = model.NewStr(`carriers/dhl/account`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlContentType = model.NewStr(`carriers/dhl/content_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlHandlingType = model.NewStr(`carriers/dhl/handling_type`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlHandlingAction = model.NewStr(`carriers/dhl/handling_action`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlHandlingFee = model.NewStr(`carriers/dhl/handling_fee`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlDivideOrderWeight = model.NewBool(`carriers/dhl/divide_order_weight`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlUnitOfMeasure = model.NewStr(`carriers/dhl/unit_of_measure`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlSize = model.NewStr(`carriers/dhl/size`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlHeight = model.NewStr(`carriers/dhl/height`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlDepth = model.NewStr(`carriers/dhl/depth`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlWidth = model.NewStr(`carriers/dhl/width`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlDocMethods = model.NewStringCSV(`carriers/dhl/doc_methods`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlNondocMethods = model.NewStringCSV(`carriers/dhl/nondoc_methods`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlReadyTime = model.NewStr(`carriers/dhl/ready_time`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlSpecificerrmsg = model.NewStr(`carriers/dhl/specificerrmsg`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlFreeMethodDoc = model.NewStr(`carriers/dhl/free_method_doc`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlFreeMethodNondoc = model.NewStr(`carriers/dhl/free_method_nondoc`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlFreeShippingEnable = model.NewBool(`carriers/dhl/free_shipping_enable`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlFreeShippingSubtotal = model.NewStr(`carriers/dhl/free_shipping_subtotal`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlSallowspecific = model.NewStr(`carriers/dhl/sallowspecific`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlSpecificcountry = model.NewStringCSV(`carriers/dhl/specificcountry`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlShowmethod = model.NewBool(`carriers/dhl/showmethod`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlSortOrder = model.NewStr(`carriers/dhl/sort_order`, model.WithConfigStructure(cfgStruct))
	pp.CarriersDhlDebug = model.NewBool(`carriers/dhl/debug`, model.WithConfigStructure(cfgStruct))

	return pp
}
