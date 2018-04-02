// +build ignore

package fedex

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID: "carriers",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "fedex",
					Label:     `FedEx`,
					SortOrder: 120,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: carriers/fedex/active
							ID:        "active",
							Label:     `Enabled for Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/active_rma
							ID:        "active_rma",
							Label:     `Enabled for RMA`,
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Federal Express`,
						},

						element.Field{
							// Path: carriers/fedex/account
							ID:        "account",
							Label:     `Account ID`,
							Comment:   text.Long(`Please make sure to use only digits here. No dashes are allowed.`),
							Type:      element.TypeObscure,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/fedex/meter_number
							ID:        "meter_number",
							Label:     `Meter Number`,
							Type:      element.TypeObscure,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/fedex/key
							ID:        "key",
							Label:     `Key`,
							Type:      element.TypeObscure,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/fedex/password
							ID:        "password",
							Label:     `Password`,
							Type:      element.TypeObscure,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   nil,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted @todo Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: carriers/fedex/sandbox_mode
							ID:        "sandbox_mode",
							Label:     `Sandbox Mode`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/production_webservices_url
							ID:        "production_webservices_url",
							Label:     `Web-Services URL (Production)`,
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `https://ws.fedex.com:443/web-services/`,
						},

						element.Field{
							// Path: carriers/fedex/sandbox_webservices_url
							ID:        "sandbox_webservices_url",
							Label:     `Web-Services URL (Sandbox)`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `https://wsbeta.fedex.com:443/web-services/`,
						},

						element.Field{
							// Path: carriers/fedex/shipment_requesttype
							ID:        "shipment_requesttype",
							Label:     `Packages Request Type`,
							Type:      element.TypeSelect,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Online\Requesttype
						},

						element.Field{
							// Path: carriers/fedex/packaging
							ID:        "packaging",
							Label:     `Packaging`,
							Type:      element.TypeSelect,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `YOUR_PACKAGING`,
							// SourceModel: Magento\Fedex\Model\Source\Packaging
						},

						element.Field{
							// Path: carriers/fedex/dropoff
							ID:        "dropoff",
							Label:     `Dropoff`,
							Type:      element.TypeSelect,
							SortOrder: 130,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `REGULAR_PICKUP`,
							// SourceModel: Magento\Fedex\Model\Source\Dropoff
						},

						element.Field{
							// Path: carriers/fedex/unit_of_measure
							ID:        "unit_of_measure",
							Label:     `Weight Unit`,
							Type:      element.TypeSelect,
							SortOrder: 135,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `LB`,
							// SourceModel: Magento\Fedex\Model\Source\Unitofmeasure
						},

						element.Field{
							// Path: carriers/fedex/max_package_weight
							ID:        "max_package_weight",
							Label:     `Maximum Package Weight (Please consult your shipping carrier for maximum supported shipping weight)`,
							Type:      element.TypeText,
							SortOrder: 140,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   150,
						},

						element.Field{
							// Path: carriers/fedex/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `F`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingType
						},

						element.Field{
							// Path: carriers/fedex/handling_action
							ID:        "handling_action",
							Label:     `Handling Applied`,
							Type:      element.TypeSelect,
							SortOrder: 160,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `O`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingAction
						},

						element.Field{
							// Path: carriers/fedex/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 170,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/fedex/residence_delivery
							ID:        "residence_delivery",
							Label:     `Residential Delivery`,
							Type:      element.TypeSelect,
							SortOrder: 180,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/allowed_methods
							ID:         "allowed_methods",
							Label:      `Allowed Methods`,
							Type:       element.TypeMultiselect,
							SortOrder:  190,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							Default:    `EUROPE_FIRST_INTERNATIONAL_PRIORITY,FEDEX_1_DAY_FREIGHT,FEDEX_2_DAY_FREIGHT,FEDEX_2_DAY,FEDEX_2_DAY_AM,FEDEX_3_DAY_FREIGHT,FEDEX_EXPRESS_SAVER,FEDEX_GROUND,FIRST_OVERNIGHT,GROUND_HOME_DELIVERY,INTERNATIONAL_ECONOMY,INTERNATIONAL_ECONOMY_FREIGHT,INTERNATIONAL_FIRST,INTERNATIONAL_GROUND,INTERNATIONAL_PRIORITY,INTERNATIONAL_PRIORITY_FREIGHT,PRIORITY_OVERNIGHT,SMART_POST,STANDARD_OVERNIGHT,FEDEX_FREIGHT,FEDEX_NATIONAL_FREIGHT`,
							// SourceModel: Magento\Fedex\Model\Source\Method
						},

						element.Field{
							// Path: carriers/fedex/smartpost_hubid
							ID:        "smartpost_hubid",
							Label:     `Hub ID`,
							Comment:   text.Long(`The field is applicable if the Smart Post method is selected.`),
							Type:      element.TypeText,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/fedex/free_method
							ID:        "free_method",
							Label:     `Free Method`,
							Type:      element.TypeSelect,
							SortOrder: 210,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `FEDEX_GROUND`,
							// SourceModel: Magento\Fedex\Model\Source\Freemethod
						},

						element.Field{
							// Path: carriers/fedex/free_shipping_enable
							ID:        "free_shipping_enable",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeSelect,
							SortOrder: 220,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						element.Field{
							// Path: carriers/fedex/free_shipping_subtotal
							ID:        "free_shipping_subtotal",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeText,
							SortOrder: 230,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: carriers/fedex/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 240,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						},

						element.Field{
							// Path: carriers/fedex/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 250,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
						},

						element.Field{
							// Path: carriers/fedex/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  260,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: carriers/fedex/debug
							ID:        "debug",
							Label:     `Debug`,
							Type:      element.TypeSelect,
							SortOrder: 270,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 280,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: carriers/fedex/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 290,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "carriers",
			Groups: element.MakeGroups(
				element.Group{
					ID: "fedex",
					Fields: element.MakeFields(
						element.Field{
							// Path: carriers/fedex/cutoff_cost
							ID:      `cutoff_cost`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: carriers/fedex/handling
							ID:      `handling`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: carriers/fedex/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Fedex\Model\Carrier`,
						},

						element.Field{
							// Path: carriers/fedex/is_online
							ID:      `is_online`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
