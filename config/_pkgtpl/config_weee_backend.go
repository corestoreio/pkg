// +build ignore

package weee

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
	// TaxWeeeEnable => Enable FPT.
	// Path: tax/weee/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeEnable cfgmodel.Bool

	// TaxWeeeDisplayList => Display Prices In Product Lists.
	// Path: tax/weee/display_list
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplayList cfgmodel.Str

	// TaxWeeeDisplay => Display Prices On Product View Page.
	// Path: tax/weee/display
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplay cfgmodel.Str

	// TaxWeeeDisplaySales => Display Prices In Sales Modules.
	// Path: tax/weee/display_sales
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplaySales cfgmodel.Str

	// TaxWeeeDisplayEmail => Display Prices In Emails.
	// Path: tax/weee/display_email
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplayEmail cfgmodel.Str

	// TaxWeeeApplyVat => Apply Tax To FPT.
	// Path: tax/weee/apply_vat
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeApplyVat cfgmodel.Bool

	// TaxWeeeIncludeInSubtotal => Include FPT In Subtotal.
	// Path: tax/weee/include_in_subtotal
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeIncludeInSubtotal cfgmodel.Bool

	// SalesTotalsSortWeee => Fixed Product Tax.
	// Path: sales/totals_sort/weee
	SalesTotalsSortWeee cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.TaxWeeeEnable = cfgmodel.NewBool(`tax/weee/enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeDisplayList = cfgmodel.NewStr(`tax/weee/display_list`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeDisplay = cfgmodel.NewStr(`tax/weee/display`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeDisplaySales = cfgmodel.NewStr(`tax/weee/display_sales`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeDisplayEmail = cfgmodel.NewStr(`tax/weee/display_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeApplyVat = cfgmodel.NewBool(`tax/weee/apply_vat`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxWeeeIncludeInSubtotal = cfgmodel.NewBool(`tax/weee/include_in_subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortWeee = cfgmodel.NewStr(`sales/totals_sort/weee`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
