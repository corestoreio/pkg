// +build ignore

package tax

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathTaxClassesShippingTaxClass => Tax Class for Shipping.
// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Product
var PathTaxClassesShippingTaxClass = model.NewStr(`tax/classes/shipping_tax_class`)

// PathTaxClassesDefaultProductTaxClass => Default Tax Class for Product.
// BackendModel: Otnegam\Tax\Model\Config\TaxClass
// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Product
var PathTaxClassesDefaultProductTaxClass = model.NewStr(`tax/classes/default_product_tax_class`)

// PathTaxClassesDefaultCustomerTaxClass => Default Tax Class for Customer.
// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Customer
var PathTaxClassesDefaultCustomerTaxClass = model.NewStr(`tax/classes/default_customer_tax_class`)

// PathTaxCalculationAlgorithm => Tax Calculation Method Based On.
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Algorithm
var PathTaxCalculationAlgorithm = model.NewStr(`tax/calculation/algorithm`)

// PathTaxCalculationBasedOn => Tax Calculation Based On.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\Config\Source\Basedon
var PathTaxCalculationBasedOn = model.NewStr(`tax/calculation/based_on`)

// PathTaxCalculationPriceIncludesTax => Catalog Prices.
// This sets whether catalog prices entered from Otnegam Admin include tax.
// BackendModel: Otnegam\Tax\Model\Config\Price\IncludePrice
// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
var PathTaxCalculationPriceIncludesTax = model.NewStr(`tax/calculation/price_includes_tax`)

// PathTaxCalculationShippingIncludesTax => Shipping Prices.
// This sets whether shipping amounts entered from Otnegam Admin or obtained
// from gateways include tax.
// BackendModel: Otnegam\Tax\Model\Config\Price\IncludePrice
// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
var PathTaxCalculationShippingIncludesTax = model.NewStr(`tax/calculation/shipping_includes_tax`)

// PathTaxCalculationApplyAfterDiscount => Apply Customer Tax.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Apply
var PathTaxCalculationApplyAfterDiscount = model.NewStr(`tax/calculation/apply_after_discount`)

// PathTaxCalculationDiscountTax => Apply Discount On Prices.
// Apply discount on price including tax is calculated based on store tax if
// "Apply Tax after Discount" is selected.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
var PathTaxCalculationDiscountTax = model.NewStr(`tax/calculation/discount_tax`)

// PathTaxCalculationApplyTaxOn => Apply Tax On.
// SourceModel: Otnegam\Tax\Model\Config\Source\Apply\On
var PathTaxCalculationApplyTaxOn = model.NewStr(`tax/calculation/apply_tax_on`)

// PathTaxCalculationCrossBorderTradeEnabled => Enable Cross Border Trade.
// When catalog price includes tax, enable this setting to fix the price no
// matter what the customer's tax rate.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxCalculationCrossBorderTradeEnabled = model.NewBool(`tax/calculation/cross_border_trade_enabled`)

// PathTaxDefaultsCountry => Default Country.
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Country
var PathTaxDefaultsCountry = model.NewStr(`tax/defaults/country`)

// PathTaxDefaultsRegion => Default State.
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Region
var PathTaxDefaultsRegion = model.NewStr(`tax/defaults/region`)

// PathTaxDefaultsPostcode => Default Post Code.
var PathTaxDefaultsPostcode = model.NewStr(`tax/defaults/postcode`)

// PathTaxDisplayType => Display Product Prices In Catalog.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxDisplayType = model.NewStr(`tax/display/type`)

// PathTaxDisplayShipping => Display Shipping Prices.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxDisplayShipping = model.NewStr(`tax/display/shipping`)

// PathTaxCartDisplayPrice => Display Prices.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxCartDisplayPrice = model.NewStr(`tax/cart_display/price`)

// PathTaxCartDisplaySubtotal => Display Subtotal.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxCartDisplaySubtotal = model.NewStr(`tax/cart_display/subtotal`)

// PathTaxCartDisplayShipping => Display Shipping Amount.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxCartDisplayShipping = model.NewStr(`tax/cart_display/shipping`)

// PathTaxCartDisplayGrandtotal => Include Tax In Order Total.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxCartDisplayGrandtotal = model.NewBool(`tax/cart_display/grandtotal`)

// PathTaxCartDisplayFullSummary => Display Full Tax Summary.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxCartDisplayFullSummary = model.NewBool(`tax/cart_display/full_summary`)

// PathTaxCartDisplayZeroTax => Display Zero Tax Subtotal.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxCartDisplayZeroTax = model.NewBool(`tax/cart_display/zero_tax`)

// PathTaxSalesDisplayPrice => Display Prices.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxSalesDisplayPrice = model.NewStr(`tax/sales_display/price`)

// PathTaxSalesDisplaySubtotal => Display Subtotal.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxSalesDisplaySubtotal = model.NewStr(`tax/sales_display/subtotal`)

// PathTaxSalesDisplayShipping => Display Shipping Amount.
// BackendModel: Otnegam\Tax\Model\Config\Notification
// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
var PathTaxSalesDisplayShipping = model.NewStr(`tax/sales_display/shipping`)

// PathTaxSalesDisplayGrandtotal => Include Tax In Order Total.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxSalesDisplayGrandtotal = model.NewBool(`tax/sales_display/grandtotal`)

// PathTaxSalesDisplayFullSummary => Display Full Tax Summary.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxSalesDisplayFullSummary = model.NewBool(`tax/sales_display/full_summary`)

// PathTaxSalesDisplayZeroTax => Display Zero Tax Subtotal.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathTaxSalesDisplayZeroTax = model.NewBool(`tax/sales_display/zero_tax`)
