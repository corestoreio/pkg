// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package.
// Used in frontend and backend. See init() for details.
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(
		&element.Section{
			ID:        "carriers",
			SortOrder: 320,
			Scope:     scope.PermAll,
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "flatrate",
					Label:     `Flat Rate`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/flatrate/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/flatrate/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Fixed`,
						},

						&element.Field{
							// Path: carriers/flatrate/price
							ID:        "price",
							Label:     `Price`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   5.00,
						},

						&element.Field{
							// Path: carriers/flatrate/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `F`,
							// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
						},

						&element.Field{
							// Path: carriers/flatrate/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/flatrate/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/flatrate/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Flat Rate`,
						},

						&element.Field{
							// Path: carriers/flatrate/type
							ID:        "type",
							Label:     `Type`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `I`,
							// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Flatrate
						},

						&element.Field{
							// Path: carriers/flatrate/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/flatrate/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/flatrate/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/flatrate/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
						},
					),
				},

				&element.Group{
					ID:        "tablerate",
					Label:     `Table Rates`,
					SortOrder: 3,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/tablerate/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `F`,
							// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
						},

						&element.Field{
							// Path: carriers/tablerate/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/tablerate/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/condition_name
							ID:        "condition_name",
							Label:     `Condition`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `package_weight`,
							// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Tablerate
						},

						&element.Field{
							// Path: carriers/tablerate/include_virtual_price
							ID:        "include_virtual_price",
							Label:     `Include Virtual Products in Price Calculation`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/export
							ID:        "export",
							Label:     `Export`,
							Type:      element.TypeCustom, // @todo: Otnegam\OfflineShipping\Block\Adminhtml\Form\Field\Export,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/tablerate/import
							ID:        "import",
							Label:     `Import`,
							Type:      element.TypeCustom, // @todo: Otnegam\OfflineShipping\Block\Adminhtml\Form\Field\Import,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.WebsiteID),
							// BackendModel: Otnegam\OfflineShipping\Model\Config\Backend\Tablerate
						},

						&element.Field{
							// Path: carriers/tablerate/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Table Rate`,
						},

						&element.Field{
							// Path: carriers/tablerate/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/tablerate/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Best Way`,
						},

						&element.Field{
							// Path: carriers/tablerate/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/tablerate/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/tablerate/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
						},
					),
				},

				&element.Group{
					ID:        "freeshipping",
					Label:     `Free Shipping`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/freeshipping/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/freeshipping/free_shipping_subtotal
							ID:        "free_shipping_subtotal",
							Label:     `Minimum Order Amount`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/freeshipping/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Free`,
						},

						&element.Field{
							// Path: carriers/freeshipping/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: carriers/freeshipping/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Free Shipping`,
						},

						&element.Field{
							// Path: carriers/freeshipping/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/freeshipping/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/freeshipping/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/freeshipping/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "carriers",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "flatrate",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/flatrate/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\OfflineShipping\Model\Carrier\Flatrate`,
						},
					),
				},

				&element.Group{
					ID: "tablerate",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/tablerate/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\OfflineShipping\Model\Carrier\Tablerate`,
						},
					),
				},

				&element.Group{
					ID: "freeshipping",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/freeshipping/cutoff_cost
							ID:      `cutoff_cost`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 50,
						},

						&element.Field{
							// Path: carriers/freeshipping/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\OfflineShipping\Model\Carrier\Freeshipping`,
						},
					),
				},
			),
		},
	)
}
