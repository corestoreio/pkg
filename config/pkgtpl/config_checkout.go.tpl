package checkout

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     "Checkout",
		SortOrder: 305,
		Scope:     config.IDScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     `Checkout Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/options/onepage_checkout_enabled`,
						ID:           "onepage_checkout_enabled",
						Label:        `Enable Onepage Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `checkout/options/guest_checkout`,
						ID:           "guest_checkout",
						Label:        `Allow Guest Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `checkout/options/customer_must_be_logged`,
						ID:           "customer_must_be_logged",
						Label:        `Require Customer To Be Logged In To Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    15,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "cart",
				Label:     `Shopping Cart`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/cart/delete_quote_after`,
						ID:           "delete_quote_after",
						Label:        `Quote Lifetime (days)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      30,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `checkout/cart/redirect_to_cart`,
						ID:           "redirect_to_cart",
						Label:        `After Adding a Product Redirect to Shopping Cart`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "cart_link",
				Label:     `My Cart Link`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/cart_link/use_qty`,
						ID:           "use_qty",
						Label:        `Display Cart Summary`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Checkout\Model\Config\Source\Cart\Summary
					},
				},
			},

			&config.Group{
				ID:        "sidebar",
				Label:     `Shopping Cart Sidebar`,
				Comment:   ``,
				SortOrder: 4,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/sidebar/display`,
						ID:           "display",
						Label:        `Display Shopping Cart Sidebar`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `checkout/sidebar/count`,
						ID:           "count",
						Label:        `Maximum Display Recently Added Item(s)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      5,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "payment_failed",
				Label:     `Payment Failed Emails`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/payment_failed/identity`,
						ID:           "identity",
						Label:        `Payment Failed Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `checkout/payment_failed/receiver`,
						ID:           "receiver",
						Label:        `Payment Failed Email Receiver`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `checkout/payment_failed/template`,
						ID:           "template",
						Label:        `Payment Failed Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `checkout_payment_failed_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `checkout/payment_failed/copy_to`,
						ID:           "copy_to",
						Label:        `Send Payment Failed Email Copy To`,
						Comment:      `Separate by ",".`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `checkout/payment_failed/copy_method`,
						ID:           "copy_method",
						Label:        `Send Payment Failed Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},
		},
	},
)
