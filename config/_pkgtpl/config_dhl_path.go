// +build ignore

package dhl

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCarriersDhlActive => Enabled for Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersDhlActive = model.NewBool(`carriers/dhl/active`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlActiveRma => Enabled for RMA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersDhlActiveRma = model.NewBool(`carriers/dhl/active_rma`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlGatewayUrl => Gateway URL.
var PathCarriersDhlGatewayUrl = model.NewStr(`carriers/dhl/gateway_url`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlTitle => Title.
var PathCarriersDhlTitle = model.NewStr(`carriers/dhl/title`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlId => Access ID.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersDhlId = model.NewStr(`carriers/dhl/id`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlPassword => Password.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathCarriersDhlPassword = model.NewStr(`carriers/dhl/password`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlAccount => Account Number.
var PathCarriersDhlAccount = model.NewStr(`carriers/dhl/account`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlContentType => Content Type.
// SourceModel: Otnegam\Dhl\Model\Source\Contenttype
var PathCarriersDhlContentType = model.NewStr(`carriers/dhl/content_type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlHandlingType => Calculate Handling Fee.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
var PathCarriersDhlHandlingType = model.NewStr(`carriers/dhl/handling_type`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlHandlingAction => Handling Applied.
// "Per Order" allows a single handling fee for the entire order. "Per
// Package" allows an individual handling fee for each package.
// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
var PathCarriersDhlHandlingAction = model.NewStr(`carriers/dhl/handling_action`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlHandlingFee => Handling Fee.
var PathCarriersDhlHandlingFee = model.NewStr(`carriers/dhl/handling_fee`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlDivideOrderWeight => Divide Order Weight.
// Select this to allow DHL to optimize shipping charges by splitting the
// order if it exceeds 70 kg.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersDhlDivideOrderWeight = model.NewBool(`carriers/dhl/divide_order_weight`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlUnitOfMeasure => Weight Unit.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Unitofmeasure
var PathCarriersDhlUnitOfMeasure = model.NewStr(`carriers/dhl/unit_of_measure`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlSize => Size.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Size
var PathCarriersDhlSize = model.NewStr(`carriers/dhl/size`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlHeight => Height.
var PathCarriersDhlHeight = model.NewStr(`carriers/dhl/height`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlDepth => Depth.
var PathCarriersDhlDepth = model.NewStr(`carriers/dhl/depth`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlWidth => Width.
var PathCarriersDhlWidth = model.NewStr(`carriers/dhl/width`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlDocMethods => Allowed Methods.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Doc
var PathCarriersDhlDocMethods = model.NewStringCSV(`carriers/dhl/doc_methods`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlNondocMethods => Allowed Methods.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Nondoc
var PathCarriersDhlNondocMethods = model.NewStringCSV(`carriers/dhl/nondoc_methods`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlReadyTime => Ready time.
// Package ready time after order submission (in hours)
var PathCarriersDhlReadyTime = model.NewStr(`carriers/dhl/ready_time`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlSpecificerrmsg => Displayed Error Message.
var PathCarriersDhlSpecificerrmsg = model.NewStr(`carriers/dhl/specificerrmsg`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlFreeMethodDoc => Free Method.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Freedoc
var PathCarriersDhlFreeMethodDoc = model.NewStr(`carriers/dhl/free_method_doc`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlFreeMethodNondoc => Free Method.
// SourceModel: Otnegam\Dhl\Model\Source\Method\Freenondoc
var PathCarriersDhlFreeMethodNondoc = model.NewStr(`carriers/dhl/free_method_nondoc`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlFreeShippingEnable => Free Shipping Amount Threshold.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathCarriersDhlFreeShippingEnable = model.NewBool(`carriers/dhl/free_shipping_enable`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlFreeShippingSubtotal => Free Shipping Amount Threshold.
var PathCarriersDhlFreeShippingSubtotal = model.NewStr(`carriers/dhl/free_shipping_subtotal`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlSallowspecific => Ship to Applicable Countries.
// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
var PathCarriersDhlSallowspecific = model.NewStr(`carriers/dhl/sallowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlSpecificcountry => Ship to Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathCarriersDhlSpecificcountry = model.NewStringCSV(`carriers/dhl/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlShowmethod => Show Method if Not Applicable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersDhlShowmethod = model.NewBool(`carriers/dhl/showmethod`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlSortOrder => Sort Order.
var PathCarriersDhlSortOrder = model.NewStr(`carriers/dhl/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathCarriersDhlDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCarriersDhlDebug = model.NewBool(`carriers/dhl/debug`, model.WithPkgCfg(PackageConfiguration))
