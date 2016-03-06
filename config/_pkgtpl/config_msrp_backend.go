// +build ignore

package msrp

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
	// SalesMsrpEnabled => Enable MAP.
	// Warning! Enabling MAP by default will hide all product prices on
	// Storefront.
	// Path: sales/msrp/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMsrpEnabled model.Bool

	// SalesMsrpDisplayPriceType => Display Actual Price.
	// Path: sales/msrp/display_price_type
	// SourceModel: Magento\Msrp\Model\Product\Attribute\Source\Type
	SalesMsrpDisplayPriceType model.Str

	// SalesMsrpExplanationMessage => Default Popup Text Message.
	// Path: sales/msrp/explanation_message
	SalesMsrpExplanationMessage model.Str

	// SalesMsrpExplanationMessageWhatsThis => Default "What's This" Text Message.
	// Path: sales/msrp/explanation_message_whats_this
	SalesMsrpExplanationMessageWhatsThis model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesMsrpEnabled = model.NewBool(`sales/msrp/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpDisplayPriceType = model.NewStr(`sales/msrp/display_price_type`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpExplanationMessage = model.NewStr(`sales/msrp/explanation_message`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpExplanationMessageWhatsThis = model.NewStr(`sales/msrp/explanation_message_whats_this`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
