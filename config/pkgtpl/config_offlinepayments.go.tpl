package offlinepayments

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "payment",
		Label:     "",
		SortOrder: 400,
		Scope:     config.IDScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "checkmo",
				Label:     `Check / Money Order`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/checkmo/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/checkmo/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `pending`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: `payment/checkmo/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `Check / Money order`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/checkmo/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    51,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/checkmo/payable_to`,
						ID:           "payable_to",
						Label:        `Make Check Payable to`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    61,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/mailing_address`,
						ID:           "mailing_address",
						Label:        `Send Check to`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    62,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/min_order_total`,
						ID:           "min_order_total",
						Label:        `Minimum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    98,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/max_order_total`,
						ID:           "max_order_total",
						Label:        `Maximum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    99,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/checkmo/model`,
						ID:           "model",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      true,
						Scope:        config.NewScopePerm(),
						Default:      `Magento\OfflinePayments\Model\Checkmo`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "purchaseorder",
				Label:     `Purchase Order`,
				Comment:   ``,
				SortOrder: 32,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/purchaseorder/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/purchaseorder/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `pending`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: `payment/purchaseorder/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/purchaseorder/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `Purchase Order`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/purchaseorder/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/purchaseorder/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    51,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/purchaseorder/min_order_total`,
						ID:           "min_order_total",
						Label:        `Minimum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    98,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/purchaseorder/max_order_total`,
						ID:           "max_order_total",
						Label:        `Maximum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    99,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/purchaseorder/model`,
						ID:           "model",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      true,
						Scope:        config.NewScopePerm(),
						Default:      `Magento\OfflinePayments\Model\Purchaseorder`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "banktransfer",
				Label:     `Bank Transfer Payment`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/banktransfer/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/banktransfer/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `Bank Transfer Payment`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/banktransfer/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `pending`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: `payment/banktransfer/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/banktransfer/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    51,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/banktransfer/instructions`,
						ID:           "instructions",
						Label:        `Instructions`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    62,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/banktransfer/min_order_total`,
						ID:           "min_order_total",
						Label:        `Minimum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    98,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/banktransfer/max_order_total`,
						ID:           "max_order_total",
						Label:        `Maximum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    99,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/banktransfer/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "cashondelivery",
				Label:     `Cash On Delivery Payment`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/cashondelivery/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/cashondelivery/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `Cash On Delivery`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/cashondelivery/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `pending`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: `payment/cashondelivery/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/cashondelivery/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    51,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/cashondelivery/instructions`,
						ID:           "instructions",
						Label:        `Instructions`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    62,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/cashondelivery/min_order_total`,
						ID:           "min_order_total",
						Label:        `Minimum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    98,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/cashondelivery/max_order_total`,
						ID:           "max_order_total",
						Label:        `Maximum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    99,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/cashondelivery/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "free",
				Label:     `Zero Subtotal Checkout`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/free/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/free/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\Newprocessing
					},

					&config.Field{
						// Path: `payment/free/payment_action`,
						ID:           "payment_action",
						Label:        `Automatically Invoice All Items`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Source\Invoice
					},

					&config.Field{
						// Path: `payment/free/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/free/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/free/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/free/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    51,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/free/model`,
						ID:           "model",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      true,
						Scope:        config.NewScopePerm(),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "payment",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "checkmo",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/checkmo/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `offline`,
					},
				},
			},

			&config.Group{
				ID: "purchaseorder",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/purchaseorder/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `offline`,
					},
				},
			},

			&config.Group{
				ID: "banktransfer",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/banktransfer/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `Magento\OfflinePayments\Model\Banktransfer`,
					},

					&config.Field{
						// Path: `payment/banktransfer/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `offline`,
					},
				},
			},

			&config.Group{
				ID: "cashondelivery",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/cashondelivery/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `Magento\OfflinePayments\Model\Cashondelivery`,
					},

					&config.Field{
						// Path: `payment/cashondelivery/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `offline`,
					},
				},
			},

			&config.Group{
				ID: "free",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/free/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `offline`,
					},
				},
			},
		},
	},
)
