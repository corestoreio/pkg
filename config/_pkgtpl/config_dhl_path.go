// +build ignore

package dhl

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
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CarriersDhlActive = model.NewBool(`carriers/dhl/active`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlActiveRma = model.NewBool(`carriers/dhl/active_rma`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlGatewayUrl = model.NewStr(`carriers/dhl/gateway_url`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlTitle = model.NewStr(`carriers/dhl/title`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlId = model.NewStr(`carriers/dhl/id`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlPassword = model.NewStr(`carriers/dhl/password`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlAccount = model.NewStr(`carriers/dhl/account`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlContentType = model.NewStr(`carriers/dhl/content_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlHandlingType = model.NewStr(`carriers/dhl/handling_type`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlHandlingAction = model.NewStr(`carriers/dhl/handling_action`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlHandlingFee = model.NewStr(`carriers/dhl/handling_fee`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlDivideOrderWeight = model.NewBool(`carriers/dhl/divide_order_weight`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlUnitOfMeasure = model.NewStr(`carriers/dhl/unit_of_measure`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlSize = model.NewStr(`carriers/dhl/size`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlHeight = model.NewStr(`carriers/dhl/height`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlDepth = model.NewStr(`carriers/dhl/depth`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlWidth = model.NewStr(`carriers/dhl/width`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlDocMethods = model.NewStringCSV(`carriers/dhl/doc_methods`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlNondocMethods = model.NewStringCSV(`carriers/dhl/nondoc_methods`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlReadyTime = model.NewStr(`carriers/dhl/ready_time`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlSpecificerrmsg = model.NewStr(`carriers/dhl/specificerrmsg`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlFreeMethodDoc = model.NewStr(`carriers/dhl/free_method_doc`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlFreeMethodNondoc = model.NewStr(`carriers/dhl/free_method_nondoc`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlFreeShippingEnable = model.NewBool(`carriers/dhl/free_shipping_enable`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlFreeShippingSubtotal = model.NewStr(`carriers/dhl/free_shipping_subtotal`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlSallowspecific = model.NewStr(`carriers/dhl/sallowspecific`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlSpecificcountry = model.NewStringCSV(`carriers/dhl/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlShowmethod = model.NewBool(`carriers/dhl/showmethod`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlSortOrder = model.NewStr(`carriers/dhl/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.CarriersDhlDebug = model.NewBool(`carriers/dhl/debug`, model.WithPkgCfg(pkgCfg))

	return pp
}
