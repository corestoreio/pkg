// +build ignore

package dhl

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
				ID:        "dhl",
				Label:     `DHL`,
				SortOrder: 140,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/dhl/active
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
						// Path: carriers/dhl/active_rma
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
						// Path: carriers/dhl/gateway_url
						ID:        "gateway_url",
						Label:     `Gateway URL`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `https://xmlpi-ea.dhl.com/XMLShippingServlet`,
					},

					&config.Field{
						// Path: carriers/dhl/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `DHL`,
					},

					&config.Field{
						// Path: carriers/dhl/id
						ID:        "id",
						Label:     `Access ID`,
						Type:      config.TypeObscure,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/dhl/password
						ID:        "password",
						Label:     `Password`,
						Type:      config.TypeObscure,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   nil,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: carriers/dhl/account
						ID:        "account",
						Label:     `Account Number`,
						Type:      config.TypeText,
						SortOrder: 70,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/dhl/content_type
						ID:        "content_type",
						Label:     `Content Type`,
						Type:      config.TypeSelect,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `N`,
						// SourceModel: Otnegam\Dhl\Model\Source\Contenttype
					},

					&config.Field{
						// Path: carriers/dhl/handling_type
						ID:        "handling_type",
						Label:     `Calculate Handling Fee`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `F`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: carriers/dhl/handling_action
						ID:        "handling_action",
						Label:     `Handling Applied`,
						Comment:   element.LongText(`"Per Order" allows a single handling fee for the entire order. "Per Package" allows an individual handling fee for each package.`),
						Type:      config.TypeSelect,
						SortOrder: 110,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `O`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
					},

					&config.Field{
						// Path: carriers/dhl/handling_fee
						ID:        "handling_fee",
						Label:     `Handling Fee`,
						Type:      config.TypeText,
						SortOrder: 120,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/dhl/divide_order_weight
						ID:        "divide_order_weight",
						Label:     `Divide Order Weight`,
						Comment:   element.LongText(`Select this to allow DHL to optimize shipping charges by splitting the order if it exceeds 70 kg.`),
						Type:      config.TypeSelect,
						SortOrder: 130,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/dhl/unit_of_measure
						ID:        "unit_of_measure",
						Label:     `Weight Unit`,
						Type:      config.TypeSelect,
						SortOrder: 140,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `K`,
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Unitofmeasure
					},

					&config.Field{
						// Path: carriers/dhl/size
						ID:        "size",
						Label:     `Size`,
						Type:      config.TypeSelect,
						SortOrder: 150,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `R`,
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Size
					},

					&config.Field{
						// Path: carriers/dhl/height
						ID:        "height",
						Label:     `Height`,
						Type:      config.TypeText,
						SortOrder: 151,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: carriers/dhl/depth
						ID:        "depth",
						Label:     `Depth`,
						Type:      config.TypeText,
						SortOrder: 152,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: carriers/dhl/width
						ID:        "width",
						Label:     `Width`,
						Type:      config.TypeText,
						SortOrder: 153,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: carriers/dhl/doc_methods
						ID:        "doc_methods",
						Label:     `Allowed Methods`,
						Type:      config.TypeMultiselect,
						SortOrder: 170,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `2,5,6,7,9,B,C,D,U,K,L,G,W,I,N,O,R,S,T,X`,
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Doc
					},

					&config.Field{
						// Path: carriers/dhl/nondoc_methods
						ID:        "nondoc_methods",
						Label:     `Allowed Methods`,
						Type:      config.TypeMultiselect,
						SortOrder: 170,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `1,3,4,8,P,Q,E,F,H,J,M,V,Y`,
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Nondoc
					},

					&config.Field{
						// Path: carriers/dhl/ready_time
						ID:        "ready_time",
						Label:     `Ready time`,
						Comment:   element.LongText(`Package ready time after order submission (in hours)`),
						Type:      config.TypeText,
						SortOrder: 180,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/dhl/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 800,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
					},

					&config.Field{
						// Path: carriers/dhl/free_method_doc
						ID:        "free_method_doc",
						Label:     `Free Method`,
						Type:      config.TypeSelect,
						SortOrder: 1200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Freedoc
					},

					&config.Field{
						// Path: carriers/dhl/free_method_nondoc
						ID:        "free_method_nondoc",
						Label:     `Free Method`,
						Type:      config.TypeSelect,
						SortOrder: 1200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Dhl\Model\Source\Method\Freenondoc
					},

					&config.Field{
						// Path: carriers/dhl/free_shipping_enable
						ID:        "free_shipping_enable",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeSelect,
						SortOrder: 1210,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: carriers/dhl/free_shipping_subtotal
						ID:        "free_shipping_subtotal",
						Label:     `Free Shipping Amount Threshold`,
						Type:      config.TypeText,
						SortOrder: 1220,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/dhl/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 1900,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/dhl/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  1910,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/dhl/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 1940,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/dhl/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 2000,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/dhl/debug
						ID:        "debug",
						Label:     `Debug`,
						Type:      config.TypeSelect,
						SortOrder: 1950,
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
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/allowed_resources
						ID:      `allowed_resources`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"dhl_folder":"dhl"}`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "carriers",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "dhl",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/dhl/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Dhl\Model\Carrier`,
					},

					&config.Field{
						// Path: carriers/dhl/account
						ID:      `account`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},

					&config.Field{
						// Path: carriers/dhl/free_method
						ID:      `free_method`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `G`,
					},

					&config.Field{
						// Path: carriers/dhl/shipment_days
						ID:      `shipment_days`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Mon,Tue,Wed,Thu,Fri`,
					},

					&config.Field{
						// Path: carriers/dhl/is_online
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
