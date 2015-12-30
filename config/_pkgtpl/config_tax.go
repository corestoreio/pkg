// +build ignore

package tax

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID:        "tax",
		Label:     `Tax`,
		SortOrder: 303,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Tax::config_tax
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "classes",
				Label:     `Tax Classes`,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/classes/shipping_tax_class
						ID:        "shipping_tax_class",
						Label:     `Tax Class for Shipping`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Product
					},

					&element.Field{
						// Path: tax/classes/default_product_tax_class
						ID:        "default_product_tax_class",
						Label:     `Default Tax Class for Product`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   2,
						// BackendModel: Otnegam\Tax\Model\Config\TaxClass
						// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Product
					},

					&element.Field{
						// Path: tax/classes/default_customer_tax_class
						ID:        "default_customer_tax_class",
						Label:     `Default Tax Class for Customer`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   3,
						// SourceModel: Otnegam\Tax\Model\TaxClass\Source\Customer
					},
				),
			},

			&element.Group{
				ID:        "calculation",
				Label:     `Calculation Settings`,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/calculation/algorithm
						ID:        "algorithm",
						Label:     `Tax Calculation Method Based On`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `TOTAL_BASE_CALCULATION`,
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Algorithm
					},

					&element.Field{
						// Path: tax/calculation/based_on
						ID:        "based_on",
						Label:     `Tax Calculation Based On`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `shipping`,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\Config\Source\Basedon
					},

					&element.Field{
						// Path: tax/calculation/price_includes_tax
						ID:        "price_includes_tax",
						Label:     `Catalog Prices`,
						Comment:   element.LongText(`This sets whether catalog prices entered from Otnegam Admin include tax.`),
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// BackendModel: Otnegam\Tax\Model\Config\Price\IncludePrice
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
					},

					&element.Field{
						// Path: tax/calculation/shipping_includes_tax
						ID:        "shipping_includes_tax",
						Label:     `Shipping Prices`,
						Comment:   element.LongText(`This sets whether shipping amounts entered from Otnegam Admin or obtained from gateways include tax.`),
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// BackendModel: Otnegam\Tax\Model\Config\Price\IncludePrice
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
					},

					&element.Field{
						// Path: tax/calculation/apply_after_discount
						ID:        "apply_after_discount",
						Label:     `Apply Customer Tax`,
						Type:      element.TypeSelect,
						SortOrder: 40,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Apply
					},

					&element.Field{
						// Path: tax/calculation/discount_tax
						ID:        "discount_tax",
						Label:     `Apply Discount On Prices`,
						Comment:   element.LongText(`Apply discount on price including tax is calculated based on store tax if "Apply Tax after Discount" is selected.`),
						Type:      element.TypeSelect,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `["0","0"]`,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\PriceType
					},

					&element.Field{
						// Path: tax/calculation/apply_tax_on
						ID:        "apply_tax_on",
						Label:     `Apply Tax On`,
						Type:      element.TypeSelect,
						SortOrder: 60,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Tax\Model\Config\Source\Apply\On
					},

					&element.Field{
						// Path: tax/calculation/cross_border_trade_enabled
						ID:        "cross_border_trade_enabled",
						Label:     `Enable Cross Border Trade`,
						Comment:   element.LongText(`When catalog price includes tax, enable this setting to fix the price no matter what the customer's tax rate.`),
						Type:      element.TypeSelect,
						SortOrder: 70,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&element.Group{
				ID:        "defaults",
				Label:     `Default Tax Destination Calculation`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/defaults/country
						ID:        "country",
						Label:     `Default Country`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `US`,
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Country
					},

					&element.Field{
						// Path: tax/defaults/region
						ID:        "region",
						Label:     `Default State`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Region
					},

					&element.Field{
						// Path: tax/defaults/postcode
						ID:        "postcode",
						Label:     `Default Post Code`,
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},

			&element.Group{
				ID:        "display",
				Label:     `Price Display Settings`,
				SortOrder: 40,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/display/type
						ID:        "type",
						Label:     `Display Product Prices In Catalog`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/display/shipping
						ID:        "shipping",
						Label:     `Display Shipping Prices`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},
				),
			},

			&element.Group{
				ID:        "cart_display",
				Label:     `Shopping Cart Display Settings`,
				SortOrder: 50,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/cart_display/price
						ID:        "price",
						Label:     `Display Prices`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/cart_display/subtotal
						ID:        "subtotal",
						Label:     `Display Subtotal`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/cart_display/shipping
						ID:        "shipping",
						Label:     `Display Shipping Amount`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/cart_display/grandtotal
						ID:        "grandtotal",
						Label:     `Include Tax In Order Total`,
						Type:      element.TypeSelect,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: tax/cart_display/full_summary
						ID:        "full_summary",
						Label:     `Display Full Tax Summary`,
						Type:      element.TypeSelect,
						SortOrder: 60,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: tax/cart_display/zero_tax
						ID:        "zero_tax",
						Label:     `Display Zero Tax Subtotal`,
						Type:      element.TypeSelect,
						SortOrder: 120,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&element.Group{
				ID:        "sales_display",
				Label:     `Orders, Invoices, Credit Memos Display Settings`,
				SortOrder: 60,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/sales_display/price
						ID:        "price",
						Label:     `Display Prices`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/sales_display/subtotal
						ID:        "subtotal",
						Label:     `Display Subtotal`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/sales_display/shipping
						ID:        "shipping",
						Label:     `Display Shipping Amount`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Tax\Model\Config\Notification
						// SourceModel: Otnegam\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&element.Field{
						// Path: tax/sales_display/grandtotal
						ID:        "grandtotal",
						Label:     `Include Tax In Order Total`,
						Type:      element.TypeSelect,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: tax/sales_display/full_summary
						ID:        "full_summary",
						Label:     `Display Full Tax Summary`,
						Type:      element.TypeSelect,
						SortOrder: 60,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: tax/sales_display/zero_tax
						ID:        "zero_tax",
						Label:     `Display Zero Tax Subtotal`,
						Type:      element.TypeSelect,
						SortOrder: 120,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&element.Section{
		ID: "tax",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "cart_display",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/cart_display/discount
						ID:      `discount`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: true,
					},
				),
			},

			&element.Group{
				ID: "sales_display",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/sales_display/discount
						ID:      `discount`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: true,
					},
				),
			},

			&element.Group{
				ID: "notification",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: tax/notification/info_url
						ID:      `info_url`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: `http://docs.magento.com/m2/ce/user_guide/tax/warning-messages.html`,
					},
				),
			},
		),
	},
)
