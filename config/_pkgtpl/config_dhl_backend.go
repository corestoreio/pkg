// +build ignore

package dhl

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
	// CarriersDhlActive => Enabled for Checkout.
	// Path: carriers/dhl/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersDhlActive cfgmodel.Bool

	// CarriersDhlActiveRma => Enabled for RMA.
	// Path: carriers/dhl/active_rma
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersDhlActiveRma cfgmodel.Bool

	// CarriersDhlGatewayUrl => Gateway URL.
	// Path: carriers/dhl/gateway_url
	CarriersDhlGatewayUrl cfgmodel.Str

	// CarriersDhlTitle => Title.
	// Path: carriers/dhl/title
	CarriersDhlTitle cfgmodel.Str

	// CarriersDhlId => Access ID.
	// Path: carriers/dhl/id
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersDhlId cfgmodel.Str

	// CarriersDhlPassword => Password.
	// Path: carriers/dhl/password
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	CarriersDhlPassword cfgmodel.Str

	// CarriersDhlAccount => Account Number.
	// Path: carriers/dhl/account
	CarriersDhlAccount cfgmodel.Str

	// CarriersDhlContentType => Content Type.
	// Path: carriers/dhl/content_type
	// SourceModel: Magento\Dhl\Model\Source\Contenttype
	CarriersDhlContentType cfgmodel.Str

	// CarriersDhlHandlingType => Calculate Handling Fee.
	// Path: carriers/dhl/handling_type
	// SourceModel: Magento\Shipping\Model\Source\HandlingType
	CarriersDhlHandlingType cfgmodel.Str

	// CarriersDhlHandlingAction => Handling Applied.
	// "Per Order" allows a single handling fee for the entire order. "Per
	// Package" allows an individual handling fee for each package.
	// Path: carriers/dhl/handling_action
	// SourceModel: Magento\Shipping\Model\Source\HandlingAction
	CarriersDhlHandlingAction cfgmodel.Str

	// CarriersDhlHandlingFee => Handling Fee.
	// Path: carriers/dhl/handling_fee
	CarriersDhlHandlingFee cfgmodel.Str

	// CarriersDhlDivideOrderWeight => Divide Order Weight.
	// Select this to allow DHL to optimize shipping charges by splitting the
	// order if it exceeds 70 kg.
	// Path: carriers/dhl/divide_order_weight
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersDhlDivideOrderWeight cfgmodel.Bool

	// CarriersDhlUnitOfMeasure => Weight Unit.
	// Path: carriers/dhl/unit_of_measure
	// SourceModel: Magento\Dhl\Model\Source\Method\Unitofmeasure
	CarriersDhlUnitOfMeasure cfgmodel.Str

	// CarriersDhlSize => Size.
	// Path: carriers/dhl/size
	// SourceModel: Magento\Dhl\Model\Source\Method\Size
	CarriersDhlSize cfgmodel.Str

	// CarriersDhlHeight => Height.
	// Path: carriers/dhl/height
	CarriersDhlHeight cfgmodel.Str

	// CarriersDhlDepth => Depth.
	// Path: carriers/dhl/depth
	CarriersDhlDepth cfgmodel.Str

	// CarriersDhlWidth => Width.
	// Path: carriers/dhl/width
	CarriersDhlWidth cfgmodel.Str

	// CarriersDhlDocMethods => Allowed Methods.
	// Path: carriers/dhl/doc_methods
	// SourceModel: Magento\Dhl\Model\Source\Method\Doc
	CarriersDhlDocMethods cfgmodel.StringCSV

	// CarriersDhlNondocMethods => Allowed Methods.
	// Path: carriers/dhl/nondoc_methods
	// SourceModel: Magento\Dhl\Model\Source\Method\Nondoc
	CarriersDhlNondocMethods cfgmodel.StringCSV

	// CarriersDhlReadyTime => Ready time.
	// Package ready time after order submission (in hours)
	// Path: carriers/dhl/ready_time
	CarriersDhlReadyTime cfgmodel.Str

	// CarriersDhlSpecificerrmsg => Displayed Error Message.
	// Path: carriers/dhl/specificerrmsg
	CarriersDhlSpecificerrmsg cfgmodel.Str

	// CarriersDhlFreeMethodDoc => Free Method.
	// Path: carriers/dhl/free_method_doc
	// SourceModel: Magento\Dhl\Model\Source\Method\Freedoc
	CarriersDhlFreeMethodDoc cfgmodel.Str

	// CarriersDhlFreeMethodNondoc => Free Method.
	// Path: carriers/dhl/free_method_nondoc
	// SourceModel: Magento\Dhl\Model\Source\Method\Freenondoc
	CarriersDhlFreeMethodNondoc cfgmodel.Str

	// CarriersDhlFreeShippingEnable => Free Shipping Amount Threshold.
	// Path: carriers/dhl/free_shipping_enable
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CarriersDhlFreeShippingEnable cfgmodel.Bool

	// CarriersDhlFreeShippingSubtotal => Free Shipping Amount Threshold.
	// Path: carriers/dhl/free_shipping_subtotal
	CarriersDhlFreeShippingSubtotal cfgmodel.Str

	// CarriersDhlSallowspecific => Ship to Applicable Countries.
	// Path: carriers/dhl/sallowspecific
	// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
	CarriersDhlSallowspecific cfgmodel.Str

	// CarriersDhlSpecificcountry => Ship to Specific Countries.
	// Path: carriers/dhl/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	CarriersDhlSpecificcountry cfgmodel.StringCSV

	// CarriersDhlShowmethod => Show Method if Not Applicable.
	// Path: carriers/dhl/showmethod
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersDhlShowmethod cfgmodel.Bool

	// CarriersDhlSortOrder => Sort Order.
	// Path: carriers/dhl/sort_order
	CarriersDhlSortOrder cfgmodel.Str

	// CarriersDhlDebug => Debug.
	// Path: carriers/dhl/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CarriersDhlDebug cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersDhlActive = cfgmodel.NewBool(`carriers/dhl/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlActiveRma = cfgmodel.NewBool(`carriers/dhl/active_rma`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlGatewayUrl = cfgmodel.NewStr(`carriers/dhl/gateway_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlTitle = cfgmodel.NewStr(`carriers/dhl/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlId = cfgmodel.NewStr(`carriers/dhl/id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlPassword = cfgmodel.NewStr(`carriers/dhl/password`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlAccount = cfgmodel.NewStr(`carriers/dhl/account`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlContentType = cfgmodel.NewStr(`carriers/dhl/content_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlHandlingType = cfgmodel.NewStr(`carriers/dhl/handling_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlHandlingAction = cfgmodel.NewStr(`carriers/dhl/handling_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlHandlingFee = cfgmodel.NewStr(`carriers/dhl/handling_fee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlDivideOrderWeight = cfgmodel.NewBool(`carriers/dhl/divide_order_weight`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlUnitOfMeasure = cfgmodel.NewStr(`carriers/dhl/unit_of_measure`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlSize = cfgmodel.NewStr(`carriers/dhl/size`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlHeight = cfgmodel.NewStr(`carriers/dhl/height`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlDepth = cfgmodel.NewStr(`carriers/dhl/depth`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlWidth = cfgmodel.NewStr(`carriers/dhl/width`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlDocMethods = cfgmodel.NewStringCSV(`carriers/dhl/doc_methods`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlNondocMethods = cfgmodel.NewStringCSV(`carriers/dhl/nondoc_methods`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlReadyTime = cfgmodel.NewStr(`carriers/dhl/ready_time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlSpecificerrmsg = cfgmodel.NewStr(`carriers/dhl/specificerrmsg`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlFreeMethodDoc = cfgmodel.NewStr(`carriers/dhl/free_method_doc`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlFreeMethodNondoc = cfgmodel.NewStr(`carriers/dhl/free_method_nondoc`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlFreeShippingEnable = cfgmodel.NewBool(`carriers/dhl/free_shipping_enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlFreeShippingSubtotal = cfgmodel.NewStr(`carriers/dhl/free_shipping_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlSallowspecific = cfgmodel.NewStr(`carriers/dhl/sallowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlSpecificcountry = cfgmodel.NewStringCSV(`carriers/dhl/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlShowmethod = cfgmodel.NewBool(`carriers/dhl/showmethod`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlSortOrder = cfgmodel.NewStr(`carriers/dhl/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CarriersDhlDebug = cfgmodel.NewBool(`carriers/dhl/debug`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
