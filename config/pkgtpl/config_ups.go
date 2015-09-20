// +build ignore

package ups

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "carriers",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "ups",
				Label:     `UPS`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/ups/access_license_number`,
						ID:           "access_license_number",
						Label:        `Access License Number`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/active`,
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
						// Path: `carriers/ups/allowed_methods`,
						ID:           "allowed_methods",
						Label:        `Allowed Methods`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `1DM,1DML,1DA,1DAL,1DAPI,1DP,1DPL,2DM,2DML,2DA,2DAL,3DS,GND,GNDCOM,GNDRES,STD,XPR,WXS,XPRL,XDM,XDML,XPD,01,02,03,07,08,11,12,14,54,59,65`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Method
					},

					&config.Field{
						// Path: `carriers/ups/shipment_requesttype`,
						ID:           "shipment_requesttype",
						Label:        `Packages Request Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    47,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Online\Requesttype
					},

					&config.Field{
						// Path: `carriers/ups/container`,
						ID:           "container",
						Label:        `Container`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `CP`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Container
					},

					&config.Field{
						// Path: `carriers/ups/free_shipping_enable`,
						ID:           "free_shipping_enable",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    210,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `carriers/ups/free_shipping_subtotal`,
						ID:           "free_shipping_subtotal",
						Label:        `Free Shipping Amount Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    220,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/dest_type`,
						ID:           "dest_type",
						Label:        `Destination Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `RES`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\DestType
					},

					&config.Field{
						// Path: `carriers/ups/free_method`,
						ID:           "free_method",
						Label:        `Free Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `GND`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Freemethod
					},

					&config.Field{
						// Path: `carriers/ups/gateway_url`,
						ID:           "gateway_url",
						Label:        `Gateway URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `http://www.ups.com/using/services/rave/qcostcgi.cgi`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/gateway_xml_url`,
						ID:           "gateway_xml_url",
						Label:        `Gateway XML URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `https://onlinetools.ups.com/ups.app/xml/Rate`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/handling_type`,
						ID:           "handling_type",
						Label:        `Calculate Handling Fee`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `F`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: `carriers/ups/handling_action`,
						ID:           "handling_action",
						Label:        `Handling Applied`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `O`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: `carriers/ups/handling_fee`,
						ID:           "handling_fee",
						Label:        `Handling Fee`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    130,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/max_package_weight`,
						ID:           "max_package_weight",
						Label:        `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      150,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/min_package_weight`,
						ID:           "min_package_weight",
						Label:        `Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      0.1,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/origin_shipment`,
						ID:           "origin_shipment",
						Label:        `Origin of the Shipment`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `Shipments Originating in United States`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\OriginShipment
					},

					&config.Field{
						// Path: `carriers/ups/password`,
						ID:           "password",
						Label:        `Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/pickup`,
						ID:           "pickup",
						Label:        `Pickup Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `CC`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Pickup
					},

					&config.Field{
						// Path: `carriers/ups/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1000,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `United Parcel Service`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/tracking_xml_url`,
						ID:           "tracking_xml_url",
						Label:        `Tracking XML URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `https://www.ups.com/ups.app/xml/Track`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/type`,
						ID:           "type",
						Label:        `UPS Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `UPS`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Type
					},

					&config.Field{
						// Path: `carriers/ups/is_account_live`,
						ID:           "is_account_live",
						Label:        `Live account`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    25,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/ups/unit_of_measure`,
						ID:           "unit_of_measure",
						Label:        `Weight Unit`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      `LBS`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Ups\Model\Config\Source\Unitofmeasure
					},

					&config.Field{
						// Path: `carriers/ups/username`,
						ID:           "username",
						Label:        `User ID`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/negotiated_active`,
						ID:           "negotiated_active",
						Label:        `Enable Negotiated Rates`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/ups/shipper_number`,
						ID:           "shipper_number",
						Label:        `Shipper Number`,
						Comment:      `Required for negotiated rates; 6-character UPS`,
						Type:         config.TypeText,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/ups/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    900,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/ups/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    910,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/ups/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    920,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/ups/specificerrmsg`,
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
						// Path: `carriers/ups/mode_xml`,
						ID:           "mode_xml",
						Label:        `Mode`,
						Comment:      `This enables or disables SSL verification of the Magento server by UPS.`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Online\Mode
					},

					&config.Field{
						// Path: `carriers/ups/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    920,
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
		ID: "carriers",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "ups",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/ups/cutoff_cost`,
						ID:      "cutoff_cost",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `carriers/ups/handling`,
						ID:      "handling",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/ups/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: `Magento\Ups\Model\Carrier`,
					},

					&config.Field{
						// Path: `carriers/ups/active_rma`,
						ID:      "active_rma",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `carriers/ups/is_online`,
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
