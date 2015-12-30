// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "carriers",
		SortOrder: 320,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "flatrate",
				Label:     `Flat Rate`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/flatrate/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/flatrate/name
						ID:        "name",
						Label:     `Method Name`,
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Fixed`,
					},

					&config.Field{
						// Path: carriers/flatrate/price
						ID:        "price",
						Label:     `Price`,
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   5.00,
					},

					&config.Field{
						// Path: carriers/flatrate/handling_type
						ID:        "handling_type",
						Label:     `Calculate Handling Fee`,
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `F`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: carriers/flatrate/handling_fee
						ID:        "handling_fee",
						Label:     `Handling Fee`,
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/flatrate/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/flatrate/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Flat Rate`,
					},

					&config.Field{
						// Path: carriers/flatrate/type
						ID:        "type",
						Label:     `Type`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `I`,
						// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Flatrate
					},

					&config.Field{
						// Path: carriers/flatrate/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/flatrate/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  91,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/flatrate/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 92,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/flatrate/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
					},
				),
			},

			&config.Group{
				ID:        "tablerate",
				Label:     `Table Rates`,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/tablerate/handling_type
						ID:        "handling_type",
						Label:     `Calculate Handling Fee`,
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `F`,
						// SourceModel: Otnegam\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: carriers/tablerate/handling_fee
						ID:        "handling_fee",
						Label:     `Handling Fee`,
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/tablerate/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/tablerate/condition_name
						ID:        "condition_name",
						Label:     `Condition`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `package_weight`,
						// SourceModel: Otnegam\OfflineShipping\Model\Config\Source\Tablerate
					},

					&config.Field{
						// Path: carriers/tablerate/include_virtual_price
						ID:        "include_virtual_price",
						Label:     `Include Virtual Products in Price Calculation`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/tablerate/export
						ID:        "export",
						Label:     `Export`,
						Type:      config.TypeCustom, // @todo: Otnegam\OfflineShipping\Block\Adminhtml\Form\Field\Export,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/tablerate/import
						ID:        "import",
						Label:     `Import`,
						Type:      config.TypeCustom, // @todo: Otnegam\OfflineShipping\Block\Adminhtml\Form\Field\Import,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.WebsiteID),
						// BackendModel: Otnegam\OfflineShipping\Model\Config\Backend\Tablerate
					},

					&config.Field{
						// Path: carriers/tablerate/name
						ID:        "name",
						Label:     `Method Name`,
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Table Rate`,
					},

					&config.Field{
						// Path: carriers/tablerate/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/tablerate/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Best Way`,
					},

					&config.Field{
						// Path: carriers/tablerate/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/tablerate/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  91,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/tablerate/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 92,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/tablerate/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
					},
				),
			},

			&config.Group{
				ID:        "freeshipping",
				Label:     `Free Shipping`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/freeshipping/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/freeshipping/free_shipping_subtotal
						ID:        "free_shipping_subtotal",
						Label:     `Minimum Order Amount`,
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/freeshipping/name
						ID:        "name",
						Label:     `Method Name`,
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Free`,
					},

					&config.Field{
						// Path: carriers/freeshipping/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: carriers/freeshipping/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Free Shipping`,
					},

					&config.Field{
						// Path: carriers/freeshipping/sallowspecific
						ID:        "sallowspecific",
						Label:     `Ship to Applicable Countries`,
						Type:      config.TypeSelect,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: carriers/freeshipping/specificcountry
						ID:         "specificcountry",
						Label:      `Ship to Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  91,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: carriers/freeshipping/showmethod
						ID:        "showmethod",
						Label:     `Show Method if Not Applicable`,
						Type:      config.TypeSelect,
						SortOrder: 92,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: carriers/freeshipping/specificerrmsg
						ID:        "specificerrmsg",
						Label:     `Displayed Error Message`,
						Type:      config.TypeTextarea,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `This shipping method is not available. To use this shipping method, please contact us.`,
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
				ID: "flatrate",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/flatrate/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\OfflineShipping\Model\Carrier\Flatrate`,
					},
				),
			},

			&config.Group{
				ID: "tablerate",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/tablerate/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\OfflineShipping\Model\Carrier\Tablerate`,
					},
				),
			},

			&config.Group{
				ID: "freeshipping",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: carriers/freeshipping/cutoff_cost
						ID:      `cutoff_cost`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 50,
					},

					&config.Field{
						// Path: carriers/freeshipping/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\OfflineShipping\Model\Carrier\Freeshipping`,
					},
				),
			},
		),
	},
)
