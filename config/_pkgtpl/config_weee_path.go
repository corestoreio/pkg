// +build ignore

package weee

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathTaxWeeeEnable => Enable FPT.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxWeeeEnable = model.NewBool(`tax/weee/enable`)

// PathTaxWeeeDisplayList => Display Prices In Product Lists.
// SourceModel: Otnegam\Weee\Model\Config\Source\Display
var PathTaxWeeeDisplayList = model.NewStr(`tax/weee/display_list`)

// PathTaxWeeeDisplay => Display Prices On Product View Page.
// SourceModel: Otnegam\Weee\Model\Config\Source\Display
var PathTaxWeeeDisplay = model.NewStr(`tax/weee/display`)

// PathTaxWeeeDisplaySales => Display Prices In Sales Modules.
// SourceModel: Otnegam\Weee\Model\Config\Source\Display
var PathTaxWeeeDisplaySales = model.NewStr(`tax/weee/display_sales`)

// PathTaxWeeeDisplayEmail => Display Prices In Emails.
// SourceModel: Otnegam\Weee\Model\Config\Source\Display
var PathTaxWeeeDisplayEmail = model.NewStr(`tax/weee/display_email`)

// PathTaxWeeeApplyVat => Apply Tax To FPT.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxWeeeApplyVat = model.NewBool(`tax/weee/apply_vat`)

// PathTaxWeeeIncludeInSubtotal => Include FPT In Subtotal.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxWeeeIncludeInSubtotal = model.NewBool(`tax/weee/include_in_subtotal`)

// PathSalesTotalsSortWeee => Fixed Product Tax.
var PathSalesTotalsSortWeee = model.NewStr(`sales/totals_sort/weee`)
