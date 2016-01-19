// +build ignore

package dhl

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "carriers",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "dhl",
					Label:     `DHL`,
					SortOrder: 140,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/dhl/active
							ID:        "active",
							Label:     `Enabled for Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/dhl/active_rma
							ID:        "active_rma",
							Label:     `Enabled for RMA`,
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/dhl/gateway_url
							ID:        "gateway_url",
							Label:     `Gateway URL`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `https://xmlpi-ea.dhl.com/XMLShippingServlet`,
						},

						&element.Field{
							// Path: carriers/dhl/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `DHL`,
						},

						&element.Field{
							// Path: carriers/dhl/id
							ID:        "id",
							Label:     `Access ID`,
							Type:      element.TypeObscure,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   nil,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: carriers/dhl/password
							ID:        "password",
							Label:     `Password`,
							Type:      element.TypeObscure,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   nil,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: carriers/dhl/account
							ID:        "account",
							Label:     `Account Number`,
							Type:      element.TypeText,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/dhl/content_type
							ID:        "content_type",
							Label:     `Content Type`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `N`,
							// SourceModel: Otnegam\Dhl\Model\Source\Contenttype
						},

						&element.Field{
							// Path: carriers/dhl/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `F`,
							// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
						},

						&element.Field{
							// Path: carriers/dhl/handling_action
							ID:        "handling_action",
							Label:     `Handling Applied`,
							Comment:   text.Long(`"Per Order" allows a single handling fee for the entire order. "Per Package" allows an individual handling fee for each package.`),
							Type:      element.TypeSelect,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `O`,
							// SourceModel: Otnegam\Shipping\Model\Source\HandlingAction
						},

						&element.Field{
							// Path: carriers/dhl/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/dhl/divide_order_weight
							ID:        "divide_order_weight",
							Label:     `Divide Order Weight`,
							Comment:   text.Long(`Select this to allow DHL to optimize shipping charges by splitting the order if it exceeds 70 kg.`),
							Type:      element.TypeSelect,
							SortOrder: 130,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/dhl/unit_of_measure
							ID:        "unit_of_measure",
							Label:     `Weight Unit`,
							Type:      element.TypeSelect,
							SortOrder: 140,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `K`,
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Unitofmeasure
						},

						&element.Field{
							// Path: carriers/dhl/size
							ID:        "size",
							Label:     `Size`,
							Type:      element.TypeSelect,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `R`,
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Size
						},

						&element.Field{
							// Path: carriers/dhl/height
							ID:        "height",
							Label:     `Height`,
							Type:      element.TypeText,
							SortOrder: 151,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: carriers/dhl/depth
							ID:        "depth",
							Label:     `Depth`,
							Type:      element.TypeText,
							SortOrder: 152,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: carriers/dhl/width
							ID:        "width",
							Label:     `Width`,
							Type:      element.TypeText,
							SortOrder: 153,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: carriers/dhl/doc_methods
							ID:        "doc_methods",
							Label:     `Allowed Methods`,
							Type:      element.TypeMultiselect,
							SortOrder: 170,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `2,5,6,7,9,B,C,D,U,K,L,G,W,I,N,O,R,S,T,X`,
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Doc
						},

						&element.Field{
							// Path: carriers/dhl/nondoc_methods
							ID:        "nondoc_methods",
							Label:     `Allowed Methods`,
							Type:      element.TypeMultiselect,
							SortOrder: 170,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `1,3,4,8,P,Q,E,F,H,J,M,V,Y`,
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Nondoc
						},

						&element.Field{
							// Path: carriers/dhl/ready_time
							ID:        "ready_time",
							Label:     `Ready time`,
							Comment:   text.Long(`Package ready time after order submission (in hours)`),
							Type:      element.TypeText,
							SortOrder: 180,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/dhl/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 800,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `This shipping method is currently unavailable. If you would like to ship using this shipping method, please contact us.`,
						},

						&element.Field{
							// Path: carriers/dhl/free_method_doc
							ID:        "free_method_doc",
							Label:     `Free Method`,
							Type:      element.TypeSelect,
							SortOrder: 1200,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Freedoc
						},

						&element.Field{
							// Path: carriers/dhl/free_method_nondoc
							ID:        "free_method_nondoc",
							Label:     `Free Method`,
							Type:      element.TypeSelect,
							SortOrder: 1200,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Dhl\Model\Source\Method\Freenondoc
						},

						&element.Field{
							// Path: carriers/dhl/free_shipping_enable
							ID:        "free_shipping_enable",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeSelect,
							SortOrder: 1210,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},

						&element.Field{
							// Path: carriers/dhl/free_shipping_subtotal
							ID:        "free_shipping_subtotal",
							Label:     `Free Shipping Amount Threshold`,
							Type:      element.TypeText,
							SortOrder: 1220,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/dhl/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 1900,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/dhl/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  1910,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/dhl/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 1940,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/dhl/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 2000,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/dhl/debug
							ID:        "debug",
							Label:     `Debug`,
							Type:      element.TypeSelect,
							SortOrder: 1950,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "media_storage_configuration",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"dhl_folder":"dhl"}`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "carriers",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "dhl",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/dhl/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\Dhl\Model\Carrier`,
						},

						&element.Field{
							// Path: carriers/dhl/account
							ID:      `account`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						&element.Field{
							// Path: carriers/dhl/free_method
							ID:      `free_method`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `G`,
						},

						&element.Field{
							// Path: carriers/dhl/shipment_days
							ID:      `shipment_days`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Mon,Tue,Wed,Thu,Fri`,
						},

						&element.Field{
							// Path: carriers/dhl/is_online
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
