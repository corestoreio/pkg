// +build ignore

package fedex

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
				ID:        "fedex",
				Label:     `FedEx`,
				SortOrder: 120,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/fedex/active
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
						// Path: carriers/fedex/active_rma
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
						// Path: carriers/fedex/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Federal Express`,
					},

					&config.Field{
						// Path: carriers/fedex/account
						ID:        "account",
						Label:     `Account ID`,
						Comment:   element.LongText(`Please make sure to use only digits here. No dashes are allowed.`),
						Type:      config.TypeObscure,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/fedex/meter_number
						ID:        "meter_number",
						Label:     `Meter Number`,
						Type:      config.TypeObscure,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/fedex/key
						ID:        "key",
						Label:     `Key`,
						Type:      config.TypeObscure,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/fedex/password
						ID:        "password",
						Label:     `Password`,
						Type:      config.TypeObscure,
						SortOrder: 70,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/fedex/sandbox_mode
						ID:        "sandbox_mode",
						Label:     `Sandbox Mode`,
						Type:      config.TypeSelect,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/fedex/production_webservices_url
						ID:        "production_webservices_url",
						Label:     `Web-Services URL (Production)`,
						Type:      config.TypeText,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `https://ws.fedex.com:443/web-services/`,
					},

					&config.Field{
						// Path: carriers/fedex/sandbox_webservices_url
						ID:        "sandbox_webservices_url",
						Label:     `Web-Services URL (Sandbox)`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `https://wsbeta.fedex.com:443/web-services/`,
					},

					&config.Field{
						// Path: carriers/fedex/shipment_requesttype
						ID:        "shipment_requesttype",
						Label:     `Packages Request Type`,
						Type:      config.TypeSelect,
						SortOrder: 110,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Online\Requesttype
					},

					&config.Field{
						// Path: carriers/fedex/packaging
						ID:        "packaging",
						Label:     `Packaging`,
						Type:      config.TypeSelect,
						SortOrder: 120,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `YOUR_PACKAGING`,
						// SourceModel: Otnegam\Fedex\Model\Source\Packaging
					},

					&config.Field{
						// Path: carriers/fedex/dropoff
						ID:        "dropoff",
						Label:     `Dropoff`,
						Type:      config.TypeSelect,
						SortOrder: 130,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `REGULAR_PICKUP`,
						// SourceModel: Otnegam\Fedex\Model\Source\Dropoff
					},

					&config.Field{
						// Path: carriers/fedex/unit_of_measure
						ID:        "unit_of_measure",
						Label:     `Weight Unit`,
						Type:      config.TypeSelect,
						SortOrder: 135,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `LB`,
						// SourceModel: Otnegam\Fedex\Model\Source\Unitofmeasure
					},

					&config.Field{
						// Path: carriers/fedex/max_package_weight
						ID:        "max_package_weight",
						Label:     `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
						Type:      config.TypeText,
						SortOrder: 140,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   150,
					},

					&config.Field{
						// Path: carriers/fedex/handling_type
						ID:        "handling_type",
						Label:     `Calculate Handling Fee`,
						Type:      config.TypeSelect,
						SortOrder: 150,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `F`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: carriers/fedex/handling_action
						ID:        "handling_action",
						Label:     `Handling Applied`,
						Type:      config.TypeSelect,
						SortOrder: 160,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `O`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: carriers/fedex/handling_fee
						ID:        "handling_fee",
						Label:     `Handling Fee`,
						Type:      config.TypeText,
						SortOrder: 170,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/fedex/residence_delivery
						ID:        "residence_delivery",
						Label:     `Residential Delivery`,
						Type:      config.TypeSelect,
						SortOrder: 180,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/fedex/allowed_methods
						ID:         "allowed_methods",
						Label:      `Allowed Methods`,
						Type:       config.TypeMultiselect,
						SortOrder:  190,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						Default:    `EUROPE_FIRST_INTERNATIONAL_PRIORITY,FEDEX_1_DAY_FREIGHT,FEDEX_2_DAY_FREIGHT,FEDEX_2_DAY,FEDEX_2_DAY_AM,FEDEX_3_DAY_FREIGHT,FEDEX_EXPRESS_SAVER,FEDEX_GROUND,FIRST_OVERNIGHT,GROUND_HOME_DELIVERY,INTERNATIONAL_ECONOMY,INTERNATIONAL_ECONOMY_FREIGHT,INTERNATIONAL_FIRST,INTERNATIONAL_GROUND,INTERNATIONAL_PRIORITY,INTERNATIONAL_PRIORITY_FREIGHT,PRIORITY_OVERNIGHT,SMART_POST,STANDARD_OVERNIGHT,FEDEX_FREIGHT,FEDEX_NATIONAL_FREIGHT`,
						// SourceModel: Otnegam\Fedex\Model\Source\Method
					},

					&config.Field{
						// Path: carriers/fedex/smartpost_hubid
						ID:        "smartpost_hubid",
						Label:     `Hub ID`,
						Comment:   element.LongText(`The field is applicable if the Smart Post method is selected.`),
						Type:      config.TypeText,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/fedex/free_method
						ID:        "free_method",
						Label:     `Free Method`,
						Type:      config.TypeSelect,
						SortOrder: 210,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `FEDEX_GROUND`,
						// SourceModel: Otnegam\Fedex\Model\Source\Freemethod
					},

					&config.Field{
						// Path: carriers/fedex/free_shipping_enable
						ID:        "free_shipping_enable",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeSelect,
						SortOrder: 220,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: carriers/fedex/free_shipping_subtotal
						ID:        "free_shipping_subtotal",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeText,
						SortOrder: 230,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/fedex/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 240,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
					},

					&config.Field{
						// Path: carriers/fedex/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 250,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/fedex/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  260,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/fedex/debug
						ID:        "debug",
						Label:     `Debug`,
						Type:      config.TypeSelect,
						SortOrder: 270,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/fedex/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 280,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/fedex/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 290,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
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
				ID: "fedex",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/fedex/cutoff_cost
						ID:      `cutoff_cost`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},

					&config.Field{
						// Path: carriers/fedex/handling
						ID:      `handling`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: carriers/fedex/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Fedex\Model\Carrier`,
					},

					&config.Field{
						// Path: carriers/fedex/is_online
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
