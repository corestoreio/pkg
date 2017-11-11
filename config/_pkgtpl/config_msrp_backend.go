// +build ignore

package msrp

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SalesMsrpEnabled => Enable MAP.
	// Warning! Enabling MAP by default will hide all product prices on
	// Storefront.
	// Path: sales/msrp/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMsrpEnabled cfgmodel.Bool

	// SalesMsrpDisplayPriceType => Display Actual Price.
	// Path: sales/msrp/display_price_type
	// SourceModel: Magento\Msrp\Model\Product\Attribute\Source\Type
	SalesMsrpDisplayPriceType cfgmodel.Str

	// SalesMsrpExplanationMessage => Default Popup Text Message.
	// Path: sales/msrp/explanation_message
	SalesMsrpExplanationMessage cfgmodel.Str

	// SalesMsrpExplanationMessageWhatsThis => Default "What's This" Text Message.
	// Path: sales/msrp/explanation_message_whats_this
	SalesMsrpExplanationMessageWhatsThis cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesMsrpEnabled = cfgmodel.NewBool(`sales/msrp/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpDisplayPriceType = cfgmodel.NewStr(`sales/msrp/display_price_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpExplanationMessage = cfgmodel.NewStr(`sales/msrp/explanation_message`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMsrpExplanationMessageWhatsThis = cfgmodel.NewStr(`sales/msrp/explanation_message_whats_this`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
