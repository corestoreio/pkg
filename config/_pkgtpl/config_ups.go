// +build ignore

package ups

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "carriers",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "ups",
				Label:     `UPS`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/ups/access_license_number
						ID:        "access_license_number",
						Label:     `Access License Number`,
						Type:      config.TypeObscure,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/ups/active
						ID:        "active",
						Label:     `Enabled for Checkout`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/ups/active_rma
						ID:        "active_rma",
						Label:     `Enabled for RMA`,
						Type:      config.TypeSelect,
						SortOrder: 15,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/ups/allowed_methods
						ID:         "allowed_methods",
						Label:      `Allowed Methods`,
						Type:       config.TypeMultiselect,
						SortOrder:  170,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						Default:    `1DM,1DML,1DA,1DAL,1DAPI,1DP,1DPL,2DM,2DML,2DA,2DAL,3DS,GND,GNDCOM,GNDRES,STD,XPR,WXS,XPRL,XDM,XDML,XPD,01,02,03,07,08,11,12,14,54,59,65`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Method
					},

					&config.Field{
						// Path: carriers/ups/shipment_requesttype
						ID:        "shipment_requesttype",
						Label:     `Packages Request Type`,
						Type:      config.TypeSelect,
						SortOrder: 47,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
					},

					&config.Field{
						// Path: carriers/ups/container
						ID:        "container",
						Label:     `Container`,
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `CP`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Container
					},

					&config.Field{
						// Path: carriers/ups/free_shipping_enable
						ID:        "free_shipping_enable",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeSelect,
						SortOrder: 210,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: carriers/ups/free_shipping_subtotal
						ID:        "free_shipping_subtotal",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeText,
						SortOrder: 220,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/ups/dest_type
						ID:        "dest_type",
						Label:     `Destination Type`,
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `RES`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\DestType
					},

					&config.Field{
						// Path: carriers/ups/free_method
						ID:        "free_method",
						Label:     `Free Method`,
						Type:      config.TypeSelect,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `GND`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Freemethod
					},

					&config.Field{
						// Path: carriers/ups/gateway_url
						ID:        "gateway_url",
						Label:     `Gateway URL`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `http://www.ups.com/using/services/rave/qcostcgi.cgi`,
					},

					&config.Field{
						// Path: carriers/ups/gateway_xml_url
						ID:        "gateway_xml_url",
						Label:     `Gateway XML URL`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `https://onlinetools.ups.com/ups.app/xml/Rate`,
					},

					&config.Field{
						// Path: carriers/ups/handling_type
						ID:        "handling_type",
						Label:     `Calculate Handling Fee`,
						Type:      config.TypeSelect,
						SortOrder: 110,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `F`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: carriers/ups/handling_action
						ID:        "handling_action",
						Label:     `Handling Applied`,
						Type:      config.TypeSelect,
						SortOrder: 120,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `O`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: carriers/ups/handling_fee
						ID:        "handling_fee",
						Label:     `Handling Fee`,
						Type:      config.TypeText,
						SortOrder: 130,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/ups/max_package_weight
						ID:        "max_package_weight",
						Label:     `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
						Type:      config.TypeText,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   150,
					},

					&config.Field{
						// Path: carriers/ups/min_package_weight
						ID:        "min_package_weight",
						Label:     `Minimum Package Weight (Please consult your shipping carrier for minimum supported shipping weight)`,
						Type:      config.TypeText,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   0.1,
					},

					&config.Field{
						// Path: carriers/ups/origin_shipment
						ID:        "origin_shipment",
						Label:     `Origin of the Shipment`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `Shipments Originating in United States`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\OriginShipment
					},

					&config.Field{
						// Path: carriers/ups/password
						ID:        "password",
						Label:     `Password`,
						Type:      config.TypeObscure,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/ups/pickup
						ID:        "pickup",
						Label:     `Pickup Method`,
						Type:      config.TypeSelect,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `CC`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Pickup
					},

					&config.Field{
						// Path: carriers/ups/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 1000,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/ups/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `United Parcel Service`,
					},

					&config.Field{
						// Path: carriers/ups/tracking_xml_url
						ID:        "tracking_xml_url",
						Label:     `Tracking XML URL`,
						Type:      config.TypeText,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `https://www.ups.com/ups.app/xml/Track`,
					},

					&config.Field{
						// Path: carriers/ups/type
						ID:        "type",
						Label:     `UPS Type`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `UPS`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Type
					},

					&config.Field{
						// Path: carriers/ups/is_account_live
						ID:        "is_account_live",
						Label:     `Live Account`,
						Type:      config.TypeSelect,
						SortOrder: 25,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/ups/unit_of_measure
						ID:        "unit_of_measure",
						Label:     `Weight Unit`,
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `LBS`,
						// SourceModel: Otnegam\Ups\Model\Config\Source\Unitofmeasure
					},

					&config.Field{
						// Path: carriers/ups/username
						ID:        "username",
						Label:     `User ID`,
						Type:      config.TypeObscure,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/ups/negotiated_active
						ID:        "negotiated_active",
						Label:     `Enable Negotiated Rates`,
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/ups/shipper_number
						ID:        "shipper_number",
						Label:     `Shipper Number`,
						Comment:   element.LongText(`Required for negotiated rates; 6-character UPS`),
						Type:      config.TypeText,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/ups/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 900,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/ups/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  910,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/ups/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 920,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/ups/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 800,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
					},

					&config.Field{
						// Path: carriers/ups/mode_xml
						ID:        "mode_xml",
						Label:     `Mode`,
						Comment:   element.LongText(`This enables or disables SSL verification of the Otnegam server by UPS.`),
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Mode
					},

					&config.Field{
						// Path: carriers/ups/debug
						ID:        "debug",
						Label:     `Debug`,
						Type:      config.TypeSelect,
						SortOrder: 920,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "carriers",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "ups",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/ups/cutoff_cost
						ID:      `cutoff_cost`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},

					&config.Field{
						// Path: carriers/ups/handling
						ID:      `handling`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: carriers/ups/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Ups\Model\Carrier`,
					},

					&config.Field{
						// Path: carriers/ups/is_online
						ID:      `is_online`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},
		),
	},
)
