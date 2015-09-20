// +build ignore

package dhl

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "carriers",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "dhl",
				Label:     `DHL`,
				Comment:   ``,
				SortOrder: 140,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/dhl/active`,
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
						// Path: `carriers/dhl/gateway_url`,
						ID:           "gateway_url",
						Label:        `Gateway URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `https://xmlpi-ea.dhl.com/XMLShippingServlet`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `DHL`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/id`,
						ID:           "id",
						Label:        `Access ID`,
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
						// Path: `carriers/dhl/password`,
						ID:           "password",
						Label:        `Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/account`,
						ID:           "account",
						Label:        `Account Number`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/content_type`,
						ID:           "content_type",
						Label:        `Content Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `N`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Contenttype
					},

					&config.Field{
						// Path: `carriers/dhl/handling_type`,
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
						// Path: `carriers/dhl/handling_action`,
						ID:           "handling_action",
						Label:        `Handling Applied`,
						Comment:      `"Per Order" allows a single handling fee for the entire order. "Per Package" allows an individual handling fee for each package.`,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `O`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: `carriers/dhl/handling_fee`,
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
						// Path: `carriers/dhl/divide_order_weight`,
						ID:           "divide_order_weight",
						Label:        `Divide Order Weight`,
						Comment:      `Select this to allow DHL to optimize shipping charges by splitting the order if it exceeds 70 kg.`,
						Type:         config.TypeSelect,
						SortOrder:    130,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/dhl/unit_of_measure`,
						ID:           "unit_of_measure",
						Label:        `Weight Unit`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    140,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `K`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Unitofmeasure
					},

					&config.Field{
						// Path: `carriers/dhl/size`,
						ID:           "size",
						Label:        `Size`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `R`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Size
					},

					&config.Field{
						// Path: `carriers/dhl/height`,
						ID:           "height",
						Label:        `Height`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    151,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/depth`,
						ID:           "depth",
						Label:        `Depth`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    152,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/width`,
						ID:           "width",
						Label:        `Width`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    153,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/doc_methods`,
						ID:           "doc_methods",
						Label:        `Allowed Methods`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `2,5,6,7,9,B,C,D,U,K,L,G,W,I,N,O,R,S,T,X`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Doc
					},

					&config.Field{
						// Path: `carriers/dhl/nondoc_methods`,
						ID:           "nondoc_methods",
						Label:        `Allowed Methods`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `1,3,4,8,P,Q,E,F,H,J,M,V,Y`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Nondoc
					},

					&config.Field{
						// Path: `carriers/dhl/ready_time`,
						ID:           "ready_time",
						Label:        `Ready time`,
						Comment:      `Package ready time after order submission (in hours)`,
						Type:         config.TypeText,
						SortOrder:    180,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    800,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/free_method_doc`,
						ID:           "free_method_doc",
						Label:        `Free Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1200,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Freedoc
					},

					&config.Field{
						// Path: `carriers/dhl/free_method_nondoc`,
						ID:           "free_method_nondoc",
						Label:        `Free Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1200,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Dhl\Model\Source\Method\Freenondoc
					},

					&config.Field{
						// Path: `carriers/dhl/free_shipping_enable`,
						ID:           "free_shipping_enable",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1210,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `carriers/dhl/free_shipping_subtotal`,
						ID:           "free_shipping_subtotal",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1220,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1900,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/dhl/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    1910,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/dhl/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1940,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/dhl/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2000,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/dhl/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1950,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/media_storage_configuration/allowed_resources`,
						ID:      "allowed_resources",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `{"dhl_folder":"dhl"}`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "carriers",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "dhl",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/dhl/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `Magento\Dhl\Model\Carrier`,
					},

					&config.Field{
						// Path: `carriers/dhl/account`,
						ID:      "account",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/dhl/free_method`,
						ID:      "free_method",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `G`,
					},

					&config.Field{
						// Path: `carriers/dhl/shipment_days`,
						ID:      "shipment_days",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `Mon,Tue,Wed,Thu,Fri`,
					},

					&config.Field{
						// Path: `carriers/dhl/active_rma`,
						ID:      "active_rma",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/dhl/is_online`,
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
