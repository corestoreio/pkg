// +build ignore

package msrp

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
	// SalesMsrpEnabled => Enable MAP.
	// Warning! Enabling MAP by default will hide all product prices on
	// Storefront.
	// Path: sales/msrp/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesMsrpEnabled model.Bool

	// SalesMsrpDisplayPriceType => Display Actual Price.
	// Path: sales/msrp/display_price_type
	// SourceModel: Otnegam\Msrp\Model\Product\Attribute\Source\Type
	SalesMsrpDisplayPriceType model.Str

	// SalesMsrpExplanationMessage => Default Popup Text Message.
	// Path: sales/msrp/explanation_message
	SalesMsrpExplanationMessage model.Str

	// SalesMsrpExplanationMessageWhatsThis => Default "What's This" Text Message.
	// Path: sales/msrp/explanation_message_whats_this
	SalesMsrpExplanationMessageWhatsThis model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesMsrpEnabled = model.NewBool(`sales/msrp/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SalesMsrpDisplayPriceType = model.NewStr(`sales/msrp/display_price_type`, model.WithConfigStructure(cfgStruct))
	pp.SalesMsrpExplanationMessage = model.NewStr(`sales/msrp/explanation_message`, model.WithConfigStructure(cfgStruct))
	pp.SalesMsrpExplanationMessageWhatsThis = model.NewStr(`sales/msrp/explanation_message_whats_this`, model.WithConfigStructure(cfgStruct))

	return pp
}
