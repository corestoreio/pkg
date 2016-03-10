// +build ignore

package offlineshipping

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
			ID:        "carriers",
			SortOrder: 320,
			Scopes:    scope.PermStore,
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "flatrate",
					Label:     `Flat Rate`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/flatrate/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/flatrate/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Fixed`,
						},

						&element.Field{
							// Path: carriers/flatrate/price
							ID:        "price",
							Label:     `Price`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   5.00,
						},

						&element.Field{
							// Path: carriers/flatrate/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `F`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingType
						},

						&element.Field{
							// Path: carriers/flatrate/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/flatrate/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/flatrate/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Flat Rate`,
						},

						&element.Field{
							// Path: carriers/flatrate/type
							ID:        "type",
							Label:     `Type`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `I`,
							// SourceModel: Magento\OfflineShipping\Model\Config\Source\Flatrate
						},

						&element.Field{
							// Path: carriers/flatrate/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/flatrate/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/flatrate/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/flatrate/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
						},
					),
				},

				&element.Group{
					ID:        "tablerate",
					Label:     `Table Rates`,
					SortOrder: 3,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/tablerate/handling_type
							ID:        "handling_type",
							Label:     `Calculate Handling Fee`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `F`,
							// SourceModel: Magento\Shipping\Model\Source\HandlingType
						},

						&element.Field{
							// Path: carriers/tablerate/handling_fee
							ID:        "handling_fee",
							Label:     `Handling Fee`,
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/tablerate/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/condition_name
							ID:        "condition_name",
							Label:     `Condition`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `package_weight`,
							// SourceModel: Magento\OfflineShipping\Model\Config\Source\Tablerate
						},

						&element.Field{
							// Path: carriers/tablerate/include_virtual_price
							ID:        "include_virtual_price",
							Label:     `Include Virtual Products in Price Calculation`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/export
							ID:        "export",
							Label:     `Export`,
							Type:      element.TypeCustom, // @todo: Magento\OfflineShipping\Block\Adminhtml\Form\Field\Export,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/tablerate/import
							ID:        "import",
							Label:     `Import`,
							Type:      element.TypeCustom, // @todo: Magento\OfflineShipping\Block\Adminhtml\Form\Field\Import,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\OfflineShipping\Model\Config\Backend\Tablerate
						},

						&element.Field{
							// Path: carriers/tablerate/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Table Rate`,
						},

						&element.Field{
							// Path: carriers/tablerate/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/tablerate/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Best Way`,
						},

						&element.Field{
							// Path: carriers/tablerate/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/tablerate/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/tablerate/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/tablerate/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
						},
					),
				},

				&element.Group{
					ID:        "freeshipping",
					Label:     `Free Shipping`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/freeshipping/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/freeshipping/free_shipping_subtotal
							ID:        "free_shipping_subtotal",
							Label:     `Minimum Order Amount`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/freeshipping/name
							ID:        "name",
							Label:     `Method Name`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Free`,
						},

						&element.Field{
							// Path: carriers/freeshipping/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						&element.Field{
							// Path: carriers/freeshipping/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Free Shipping`,
						},

						&element.Field{
							// Path: carriers/freeshipping/sallowspecific
							ID:        "sallowspecific",
							Label:     `Ship to Applicable Countries`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Shipping\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: carriers/freeshipping/specificcountry
							ID:         "specificcountry",
							Label:      `Ship to Specific Countries`,
							Type:       element.TypeMultiselect,
							SortOrder:  91,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: carriers/freeshipping/showmethod
							ID:        "showmethod",
							Label:     `Show Method if Not Applicable`,
							Type:      element.TypeSelect,
							SortOrder: 92,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: carriers/freeshipping/specificerrmsg
							ID:        "specificerrmsg",
							Label:     `Displayed Error Message`,
							Type:      element.TypeTextarea,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
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
							// Path: carriers/flatrate/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\OfflineShipping\Model\Carrier\Flatrate`,
						},
					),
				},

				&element.Group{
					ID: "tablerate",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: carriers/tablerate/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\OfflineShipping\Model\Carrier\Tablerate`,
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
							// Path: carriers/freeshipping/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\OfflineShipping\Model\Carrier\Freeshipping`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
