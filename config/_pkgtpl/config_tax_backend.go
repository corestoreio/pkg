// +build ignore

package tax

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
	// TaxClassesShippingTaxClass => Tax Class for Shipping.
	// Path: tax/classes/shipping_tax_class
	// SourceModel: Magento\Tax\Model\TaxClass\Source\Product
	TaxClassesShippingTaxClass cfgmodel.Str

	// TaxClassesDefaultProductTaxClass => Default Tax Class for Product.
	// Path: tax/classes/default_product_tax_class
	// BackendModel: Magento\Tax\Model\Config\TaxClass
	// SourceModel: Magento\Tax\Model\TaxClass\Source\Product
	TaxClassesDefaultProductTaxClass cfgmodel.Str

	// TaxClassesDefaultCustomerTaxClass => Default Tax Class for Customer.
	// Path: tax/classes/default_customer_tax_class
	// SourceModel: Magento\Tax\Model\TaxClass\Source\Customer
	TaxClassesDefaultCustomerTaxClass cfgmodel.Str

	// TaxCalculationAlgorithm => Tax Calculation Method Based On.
	// Path: tax/calculation/algorithm
	// SourceModel: Magento\Tax\Model\System\Config\Source\Algorithm
	TaxCalculationAlgorithm cfgmodel.Str

	// TaxCalculationBasedOn => Tax Calculation Based On.
	// Path: tax/calculation/based_on
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\Config\Source\Basedon
	TaxCalculationBasedOn cfgmodel.Str

	// TaxCalculationPriceIncludesTax => Catalog Prices.
	// This sets whether catalog prices entered from Magento Admin include tax.
	// Path: tax/calculation/price_includes_tax
	// BackendModel: Magento\Tax\Model\Config\Price\IncludePrice
	// SourceModel: Magento\Tax\Model\System\Config\Source\PriceType
	TaxCalculationPriceIncludesTax cfgmodel.Str

	// TaxCalculationShippingIncludesTax => Shipping Prices.
	// This sets whether shipping amounts entered from Magento Admin or obtained
	// from gateways include tax.
	// Path: tax/calculation/shipping_includes_tax
	// BackendModel: Magento\Tax\Model\Config\Price\IncludePrice
	// SourceModel: Magento\Tax\Model\System\Config\Source\PriceType
	TaxCalculationShippingIncludesTax cfgmodel.Str

	// TaxCalculationApplyAfterDiscount => Apply Customer Tax.
	// Path: tax/calculation/apply_after_discount
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Apply
	TaxCalculationApplyAfterDiscount cfgmodel.Str

	// TaxCalculationDiscountTax => Apply Discount On Prices.
	// Apply discount on price including tax is calculated based on store tax if
	// "Apply Tax after Discount" is selected.
	// Path: tax/calculation/discount_tax
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\PriceType
	TaxCalculationDiscountTax cfgmodel.Str

	// TaxCalculationApplyTaxOn => Apply Tax On.
	// Path: tax/calculation/apply_tax_on
	// SourceModel: Magento\Tax\Model\Config\Source\Apply\On
	TaxCalculationApplyTaxOn cfgmodel.Str

	// TaxCalculationCrossBorderTradeEnabled => Enable Cross Border Trade.
	// When catalog price includes tax, enable this setting to fix the price no
	// matter what the customer's tax rate.
	// Path: tax/calculation/cross_border_trade_enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxCalculationCrossBorderTradeEnabled cfgmodel.Bool

	// TaxDefaultsCountry => Default Country.
	// Path: tax/defaults/country
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Country
	TaxDefaultsCountry cfgmodel.Str

	// TaxDefaultsRegion => Default State.
	// Path: tax/defaults/region
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Region
	TaxDefaultsRegion cfgmodel.Str

	// TaxDefaultsPostcode => Default Post Code.
	// Path: tax/defaults/postcode
	TaxDefaultsPostcode cfgmodel.Str

	// TaxDisplayType => Display Product Prices In Catalog.
	// Path: tax/display/type
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxDisplayType cfgmodel.Str

	// TaxDisplayShipping => Display Shipping Prices.
	// Path: tax/display/shipping
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxDisplayShipping cfgmodel.Str

	// TaxCartDisplayPrice => Display Prices.
	// Path: tax/cart_display/price
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxCartDisplayPrice cfgmodel.Str

	// TaxCartDisplaySubtotal => Display Subtotal.
	// Path: tax/cart_display/subtotal
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxCartDisplaySubtotal cfgmodel.Str

	// TaxCartDisplayShipping => Display Shipping Amount.
	// Path: tax/cart_display/shipping
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxCartDisplayShipping cfgmodel.Str

	// TaxCartDisplayGrandtotal => Include Tax In Order Total.
	// Path: tax/cart_display/grandtotal
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxCartDisplayGrandtotal cfgmodel.Bool

	// TaxCartDisplayFullSummary => Display Full Tax Summary.
	// Path: tax/cart_display/full_summary
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxCartDisplayFullSummary cfgmodel.Bool

	// TaxCartDisplayZeroTax => Display Zero Tax Subtotal.
	// Path: tax/cart_display/zero_tax
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxCartDisplayZeroTax cfgmodel.Bool

	// TaxSalesDisplayPrice => Display Prices.
	// Path: tax/sales_display/price
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxSalesDisplayPrice cfgmodel.Str

	// TaxSalesDisplaySubtotal => Display Subtotal.
	// Path: tax/sales_display/subtotal
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxSalesDisplaySubtotal cfgmodel.Str

	// TaxSalesDisplayShipping => Display Shipping Amount.
	// Path: tax/sales_display/shipping
	// BackendModel: Magento\Tax\Model\Config\Notification
	// SourceModel: Magento\Tax\Model\System\Config\Source\Tax\Display\Type
	TaxSalesDisplayShipping cfgmodel.Str

	// TaxSalesDisplayGrandtotal => Include Tax In Order Total.
	// Path: tax/sales_display/grandtotal
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxSalesDisplayGrandtotal cfgmodel.Bool

	// TaxSalesDisplayFullSummary => Display Full Tax Summary.
	// Path: tax/sales_display/full_summary
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxSalesDisplayFullSummary cfgmodel.Bool

	// TaxSalesDisplayZeroTax => Display Zero Tax Subtotal.
	// Path: tax/sales_display/zero_tax
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	TaxSalesDisplayZeroTax cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.TaxClassesShippingTaxClass = cfgmodel.NewStr(`tax/classes/shipping_tax_class`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxClassesDefaultProductTaxClass = cfgmodel.NewStr(`tax/classes/default_product_tax_class`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxClassesDefaultCustomerTaxClass = cfgmodel.NewStr(`tax/classes/default_customer_tax_class`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationAlgorithm = cfgmodel.NewStr(`tax/calculation/algorithm`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationBasedOn = cfgmodel.NewStr(`tax/calculation/based_on`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationPriceIncludesTax = cfgmodel.NewStr(`tax/calculation/price_includes_tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationShippingIncludesTax = cfgmodel.NewStr(`tax/calculation/shipping_includes_tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationApplyAfterDiscount = cfgmodel.NewStr(`tax/calculation/apply_after_discount`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationDiscountTax = cfgmodel.NewStr(`tax/calculation/discount_tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationApplyTaxOn = cfgmodel.NewStr(`tax/calculation/apply_tax_on`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCalculationCrossBorderTradeEnabled = cfgmodel.NewBool(`tax/calculation/cross_border_trade_enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxDefaultsCountry = cfgmodel.NewStr(`tax/defaults/country`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxDefaultsRegion = cfgmodel.NewStr(`tax/defaults/region`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxDefaultsPostcode = cfgmodel.NewStr(`tax/defaults/postcode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxDisplayType = cfgmodel.NewStr(`tax/display/type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxDisplayShipping = cfgmodel.NewStr(`tax/display/shipping`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplayPrice = cfgmodel.NewStr(`tax/cart_display/price`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplaySubtotal = cfgmodel.NewStr(`tax/cart_display/subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplayShipping = cfgmodel.NewStr(`tax/cart_display/shipping`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplayGrandtotal = cfgmodel.NewBool(`tax/cart_display/grandtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplayFullSummary = cfgmodel.NewBool(`tax/cart_display/full_summary`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxCartDisplayZeroTax = cfgmodel.NewBool(`tax/cart_display/zero_tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplayPrice = cfgmodel.NewStr(`tax/sales_display/price`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplaySubtotal = cfgmodel.NewStr(`tax/sales_display/subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplayShipping = cfgmodel.NewStr(`tax/sales_display/shipping`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplayGrandtotal = cfgmodel.NewBool(`tax/sales_display/grandtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplayFullSummary = cfgmodel.NewBool(`tax/sales_display/full_summary`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TaxSalesDisplayZeroTax = cfgmodel.NewBool(`tax/sales_display/zero_tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
