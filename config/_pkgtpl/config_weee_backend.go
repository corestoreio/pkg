// +build ignore

package weee

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
	// TaxWeeeEnable => Enable FPT.
	// Path: tax/weee/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeEnable model.Bool

	// TaxWeeeDisplayList => Display Prices In Product Lists.
	// Path: tax/weee/display_list
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplayList model.Str

	// TaxWeeeDisplay => Display Prices On Product View Page.
	// Path: tax/weee/display
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplay model.Str

	// TaxWeeeDisplaySales => Display Prices In Sales Modules.
	// Path: tax/weee/display_sales
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplaySales model.Str

	// TaxWeeeDisplayEmail => Display Prices In Emails.
	// Path: tax/weee/display_email
	// SourceModel: Magento\Weee\Model\Config\Source\Display
	TaxWeeeDisplayEmail model.Str

	// TaxWeeeApplyVat => Apply Tax To FPT.
	// Path: tax/weee/apply_vat
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeApplyVat model.Bool

	// TaxWeeeIncludeInSubtotal => Include FPT In Subtotal.
	// Path: tax/weee/include_in_subtotal
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxWeeeIncludeInSubtotal model.Bool

	// SalesTotalsSortWeee => Fixed Product Tax.
	// Path: sales/totals_sort/weee
	SalesTotalsSortWeee model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.TaxWeeeEnable = model.NewBool(`tax/weee/enable`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeDisplayList = model.NewStr(`tax/weee/display_list`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeDisplay = model.NewStr(`tax/weee/display`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeDisplaySales = model.NewStr(`tax/weee/display_sales`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeDisplayEmail = model.NewStr(`tax/weee/display_email`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeApplyVat = model.NewBool(`tax/weee/apply_vat`, model.WithConfigStructure(cfgStruct))
	pp.TaxWeeeIncludeInSubtotal = model.NewBool(`tax/weee/include_in_subtotal`, model.WithConfigStructure(cfgStruct))
	pp.SalesTotalsSortWeee = model.NewStr(`sales/totals_sort/weee`, model.WithConfigStructure(cfgStruct))

	return pp
}
