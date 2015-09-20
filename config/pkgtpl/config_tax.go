// +build ignore

package tax

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "tax",
		Label:     "Tax",
		SortOrder: 303,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "classes",
				Label:     `Tax Classes`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/classes/shipping_tax_class`,
						ID:           "shipping_tax_class",
						Label:        `Tax Class for Shipping`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\TaxClass\Source\Product
					},

					&config.Field{
						// Path: `tax/classes/default_product_tax_class`,
						ID:           "default_product_tax_class",
						Label:        `Default Tax Class for Product`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      2,
						BackendModel: nil, // Magento\Tax\Model\Config\TaxClass
						SourceModel:  nil, // Magento\Tax\Model\TaxClass\Source\Product
					},

					&config.Field{
						// Path: `tax/classes/default_customer_tax_class`,
						ID:           "default_customer_tax_class",
						Label:        `Default Tax Class for Customer`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      3,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\TaxClass\Source\Customer
					},
				},
			},

			&config.Group{
				ID:        "calculation",
				Label:     `Calculation Settings`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/calculation/algorithm`,
						ID:           "algorithm",
						Label:        `Tax Calculation Method Based On`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `TOTAL_BASE_CALCULATION`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Algorithm
					},

					&config.Field{
						// Path: `tax/calculation/based_on`,
						ID:           "based_on",
						Label:        `Tax Calculation Based On`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `shipping`,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\Config\Source\Basedon
					},

					&config.Field{
						// Path: `tax/calculation/price_includes_tax`,
						ID:           "price_includes_tax",
						Label:        `Catalog Prices`,
						Comment:      `This sets whether catalog prices entered from Magento Admin include tax.`,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil, // Magento\Tax\Model\Config\Price\IncludePrice
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\PriceType
					},

					&config.Field{
						// Path: `tax/calculation/shipping_includes_tax`,
						ID:           "shipping_includes_tax",
						Label:        `Shipping Prices`,
						Comment:      `This sets whether shipping amounts entered from Magento Admin or obtained from gateways include tax.`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil, // Magento\Tax\Model\Config\Price\IncludePrice
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\PriceType
					},

					&config.Field{
						// Path: `tax/calculation/apply_after_discount`,
						ID:           "apply_after_discount",
						Label:        `Apply Customer Tax`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Apply
					},

					&config.Field{
						// Path: `tax/calculation/discount_tax`,
						ID:           "discount_tax",
						Label:        `Apply Discount On Prices`,
						Comment:      `Apply discount on price including tax is calculated based on store tax if "Apply Tax after Discount" is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `["0","0"]`,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\PriceType
					},

					&config.Field{
						// Path: `tax/calculation/apply_tax_on`,
						ID:           "apply_tax_on",
						Label:        `Apply Tax On`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\Config\Source\Apply\On
					},

					&config.Field{
						// Path: `tax/calculation/cross_border_trade_enabled`,
						ID:           "cross_border_trade_enabled",
						Label:        `Enable Cross Border Trade`,
						Comment:      `When catalog price includes tax, enable this setting to fix the price no matter what the customer's tax rate.`,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "defaults",
				Label:     `Default Tax Destination Calculation`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/defaults/country`,
						ID:           "country",
						Label:        `Default Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `US`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Country
					},

					&config.Field{
						// Path: `tax/defaults/region`,
						ID:           "region",
						Label:        `Default State`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Region
					},

					&config.Field{
						// Path: `tax/defaults/postcode`,
						ID:           "postcode",
						Label:        `Default Post Code`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `*`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "display",
				Label:     `Price Display Settings`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/display/type`,
						ID:           "type",
						Label:        `Display Product Prices In Catalog`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/display/shipping`,
						ID:           "shipping",
						Label:        `Display Shipping Prices`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},
				},
			},

			&config.Group{
				ID:        "cart_display",
				Label:     `Shopping Cart Display Settings`,
				Comment:   ``,
				SortOrder: 50,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/cart_display/price`,
						ID:           "price",
						Label:        `Display Prices`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/cart_display/subtotal`,
						ID:           "subtotal",
						Label:        `Display Subtotal`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/cart_display/shipping`,
						ID:           "shipping",
						Label:        `Display Shipping Amount`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/cart_display/grandtotal`,
						ID:           "grandtotal",
						Label:        `Include Tax In Grand Total`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/cart_display/full_summary`,
						ID:           "full_summary",
						Label:        `Display Full Tax Summary`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/cart_display/zero_tax`,
						ID:           "zero_tax",
						Label:        `Display Zero Tax Subtotal`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "sales_display",
				Label:     `Orders, Invoices, Credit Memos Display Settings`,
				Comment:   ``,
				SortOrder: 60,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/sales_display/price`,
						ID:           "price",
						Label:        `Display Prices`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/sales_display/subtotal`,
						ID:           "subtotal",
						Label:        `Display Subtotal`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/sales_display/shipping`,
						ID:           "shipping",
						Label:        `Display Shipping Amount`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Tax\Model\Config\Notification
						SourceModel:  nil, // Magento\Tax\Model\System\Config\Source\Tax\Display\Type
					},

					&config.Field{
						// Path: `tax/sales_display/grandtotal`,
						ID:           "grandtotal",
						Label:        `Include Tax In Grand Total`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/sales_display/full_summary`,
						ID:           "full_summary",
						Label:        `Display Full Tax Summary`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/sales_display/zero_tax`,
						ID:           "zero_tax",
						Label:        `Display Zero Tax Subtotal`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "tax",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "cart_display",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/cart_display/discount`,
						ID:      "discount",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: true,
					},
				},
			},

			&config.Group{
				ID: "sales_display",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/sales_display/discount`,
						ID:      "discount",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: true,
					},
				},
			},

			&config.Group{
				ID: "notification",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/notification/url`,
						ID:      "url",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `http://www.magentocommerce.com/knowledge-base/entry/magento-ce-18-ee-113-tax-calc`,
					},
				},
			},
		},
	},
)
