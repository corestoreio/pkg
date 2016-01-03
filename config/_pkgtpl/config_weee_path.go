// +build ignore

package weee

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
	// TaxWeeeEnable => Enable FPT.
	// Path: tax/weee/enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	TaxWeeeEnable model.Bool

	// TaxWeeeDisplayList => Display Prices In Product Lists.
	// Path: tax/weee/display_list
	// SourceModel: Otnegam\Weee\Model\Config\Source\Display
	TaxWeeeDisplayList model.Str

	// TaxWeeeDisplay => Display Prices On Product View Page.
	// Path: tax/weee/display
	// SourceModel: Otnegam\Weee\Model\Config\Source\Display
	TaxWeeeDisplay model.Str

	// TaxWeeeDisplaySales => Display Prices In Sales Modules.
	// Path: tax/weee/display_sales
	// SourceModel: Otnegam\Weee\Model\Config\Source\Display
	TaxWeeeDisplaySales model.Str

	// TaxWeeeDisplayEmail => Display Prices In Emails.
	// Path: tax/weee/display_email
	// SourceModel: Otnegam\Weee\Model\Config\Source\Display
	TaxWeeeDisplayEmail model.Str

	// TaxWeeeApplyVat => Apply Tax To FPT.
	// Path: tax/weee/apply_vat
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	TaxWeeeApplyVat model.Bool

	// TaxWeeeIncludeInSubtotal => Include FPT In Subtotal.
	// Path: tax/weee/include_in_subtotal
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	TaxWeeeIncludeInSubtotal model.Bool

	// SalesTotalsSortWeee => Fixed Product Tax.
	// Path: sales/totals_sort/weee
	SalesTotalsSortWeee model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
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
