package fedex

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "carriers",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "fedex",
				Label:     `FedEx`,
				Comment:   ``,
				SortOrder: 120,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/fedex/active`,
						ID:           "active",
						Label:        `Enabled for Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/fedex/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `Federal Express`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/account`,
						ID:           "account",
						Label:        `Account ID`,
						Comment:      `Please make sure to use only digits here. No dashes are allowed.`,
						Type:         config.TypeObscure,
						SortOrder:    40,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/meter_number`,
						ID:           "meter_number",
						Label:        `Meter Number`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/key`,
						ID:           "key",
						Label:        `Key`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    60,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/password`,
						ID:           "password",
						Label:        `Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    70,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/sandbox_mode`,
						ID:           "sandbox_mode",
						Label:        `Sandbox Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/fedex/production_webservices_url`,
						ID:           "production_webservices_url",
						Label:        `Web-Services URL (Production)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    90,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `https://ws.fedex.com:443/web-services/`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/sandbox_webservices_url`,
						ID:           "sandbox_webservices_url",
						Label:        `Web-Services URL (Sandbox)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `https://wsbeta.fedex.com:443/web-services/`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/shipment_requesttype`,
						ID:           "shipment_requesttype",
						Label:        `Packages Request Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Online\Requesttype
					},

					&config.Field{
						// Path: `carriers/fedex/packaging`,
						ID:           "packaging",
						Label:        `Packaging`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `YOUR_PACKAGING`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Fedex\Model\Source\Packaging
					},

					&config.Field{
						// Path: `carriers/fedex/dropoff`,
						ID:           "dropoff",
						Label:        `Dropoff`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    130,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `REGULAR_PICKUP`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Fedex\Model\Source\Dropoff
					},

					&config.Field{
						// Path: `carriers/fedex/unit_of_measure`,
						ID:           "unit_of_measure",
						Label:        `Weight Unit`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    135,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `LB`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Fedex\Model\Source\Unitofmeasure
					},

					&config.Field{
						// Path: `carriers/fedex/max_package_weight`,
						ID:           "max_package_weight",
						Label:        `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    140,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      150,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/handling_type`,
						ID:           "handling_type",
						Label:        `Calculate Handling Fee`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    150,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `F`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: `carriers/fedex/handling_action`,
						ID:           "handling_action",
						Label:        `Handling Applied`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    160,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `O`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: `carriers/fedex/handling_fee`,
						ID:           "handling_fee",
						Label:        `Handling Fee`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    170,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/residence_delivery`,
						ID:           "residence_delivery",
						Label:        `Residential Delivery`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    180,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/fedex/allowed_methods`,
						ID:           "allowed_methods",
						Label:        `Allowed Methods`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    190,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `EUROPE_FIRST_INTERNATIONAL_PRIORITY,FEDEX_1_DAY_FREIGHT,FEDEX_2_DAY_FREIGHT,FEDEX_2_DAY,FEDEX_2_DAY_AM,FEDEX_3_DAY_FREIGHT,FEDEX_EXPRESS_SAVER,FEDEX_GROUND,FIRST_OVERNIGHT,GROUND_HOME_DELIVERY,INTERNATIONAL_ECONOMY,INTERNATIONAL_ECONOMY_FREIGHT,INTERNATIONAL_FIRST,INTERNATIONAL_GROUND,INTERNATIONAL_PRIORITY,INTERNATIONAL_PRIORITY_FREIGHT,PRIORITY_OVERNIGHT,SMART_POST,STANDARD_OVERNIGHT,FEDEX_FREIGHT,FEDEX_NATIONAL_FREIGHT`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Fedex\Model\Source\Method
					},

					&config.Field{
						// Path: `carriers/fedex/smartpost_hubid`,
						ID:           "smartpost_hubid",
						Label:        `Hub ID`,
						Comment:      `The field is applicable if the Smart Post method is selected.`,
						Type:         config.TypeText,
						SortOrder:    200,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/free_method`,
						ID:           "free_method",
						Label:        `Free Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    210,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `FEDEX_GROUND`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Fedex\Model\Source\Freemethod
					},

					&config.Field{
						// Path: `carriers/fedex/free_shipping_enable`,
						ID:           "free_shipping_enable",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    220,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `carriers/fedex/free_shipping_subtotal`,
						ID:           "free_shipping_subtotal",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    230,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    240,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/fedex/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    250,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/fedex/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    260,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/fedex/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    270,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/fedex/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    280,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/fedex/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    290,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
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
		ID: "carriers",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "fedex",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/fedex/cutoff_cost`,
						ID:      "cutoff_cost",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/fedex/handling`,
						ID:      "handling",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/fedex/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `Magento\Fedex\Model\Carrier`,
					},

					&config.Field{
						// Path: `carriers/fedex/active_rma`,
						ID:      "active_rma",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/fedex/is_online`,
						ID:      "is_online",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: true,
					},
				},
			},
		},
	},
)
