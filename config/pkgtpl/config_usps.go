// +build ignore

package usps

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "carriers",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "usps",
				Label:     `USPS`,
				Comment:   ``,
				SortOrder: 110,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/usps/active`,
						ID:           "active",
						Label:        `Enabled for Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/usps/gateway_url`,
						ID:           "gateway_url",
						Label:        `Gateway URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `http://production.shippingapis.com/ShippingAPI.dll`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/gateway_secure_url`,
						ID:           "gateway_secure_url",
						Label:        `Secure Gateway URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `https://secure.shippingapis.com/ShippingAPI.dll`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `United States Postal Service`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/userid`,
						ID:           "userid",
						Label:        `User ID`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/password`,
						ID:           "password",
						Label:        `Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    53,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/mode`,
						ID:           "mode",
						Label:        `Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    54,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Online\Mode
					},

					&config.Field{
						// Path: `carriers/usps/shipment_requesttype`,
						ID:           "shipment_requesttype",
						Label:        `Packages Request Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    55,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Online\Requesttype
					},

					&config.Field{
						// Path: `carriers/usps/container`,
						ID:           "container",
						Label:        `Container`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `VARIABLE`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Usps\Model\Source\Container
					},

					&config.Field{
						// Path: `carriers/usps/size`,
						ID:           "size",
						Label:        `Size`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `REGULAR`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Usps\Model\Source\Size
					},

					&config.Field{
						// Path: `carriers/usps/width`,
						ID:           "width",
						Label:        `Width`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    73,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/length`,
						ID:           "length",
						Label:        `Length`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    72,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/height`,
						ID:           "height",
						Label:        `Height`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    74,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/girth`,
						ID:           "girth",
						Label:        `Girth`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    76,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/machinable`,
						ID:           "machinable",
						Label:        `Machinable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `true`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Usps\Model\Source\Machinable
					},

					&config.Field{
						// Path: `carriers/usps/max_package_weight`,
						ID:           "max_package_weight",
						Label:        `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      70,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/handling_type`,
						ID:           "handling_type",
						Label:        `Calculate Handling Fee`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `F`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: `carriers/usps/handling_action`,
						ID:           "handling_action",
						Label:        `Handling Applied`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `O`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: `carriers/usps/handling_fee`,
						ID:           "handling_fee",
						Label:        `Handling Fee`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/allowed_methods`,
						ID:           "allowed_methods",
						Label:        `Allowed Methods`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    130,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `0_FCLE,0_FCL,0_FCP,1,2,3,4,6,7,13,16,17,22,23,25,27,28,33,34,35,36,37,42,43,53,55,56,57,61,INT_1,INT_2,INT_4,INT_6,INT_7,INT_8,INT_9,INT_10,INT_11,INT_12,INT_13,INT_14,INT_15,INT_16,INT_20,INT_26`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Usps\Model\Source\Method
					},

					&config.Field{
						// Path: `carriers/usps/free_method`,
						ID:           "free_method",
						Label:        `Free Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    140,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Usps\Model\Source\Freemethod
					},

					&config.Field{
						// Path: `carriers/usps/free_shipping_enable`,
						ID:           "free_shipping_enable",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1500,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `carriers/usps/free_shipping_subtotal`,
						ID:           "free_shipping_subtotal",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    160,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/usps/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    180,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/usps/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    190,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/usps/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/usps/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    210,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/usps/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    220,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "carriers",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "usps",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/usps/cutoff_cost`,
						ID:      "cutoff_cost",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/usps/free_method`,
						ID:      "free_method",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/usps/handling`,
						ID:      "handling",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/usps/methods`,
						ID:      "methods",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/usps/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `Magento\Usps\Model\Carrier`,
					},

					&config.Field{
						// Path: `carriers/usps/isproduction`,
						ID:      "isproduction",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/usps/active_rma`,
						ID:      "active_rma",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/usps/is_online`,
						ID:      "is_online",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: true,
					},
				},
			},
		},
	},
)
