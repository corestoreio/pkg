// +build ignore

package offlineshipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "carriers",
		Label:     "",
		SortOrder: 320,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "flatrate",
				Label:     `Flat Rate`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/flatrate/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/flatrate/name`,
						ID:           "name",
						Label:        `Method Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Fixed`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/flatrate/price`,
						ID:           "price",
						Label:        `Price`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      5.00,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/flatrate/handling_type`,
						ID:           "handling_type",
						Label:        `Calculate Handling Fee`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `F`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: `carriers/flatrate/handling_fee`,
						ID:           "handling_fee",
						Label:        `Handling Fee`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    8,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/flatrate/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/flatrate/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Flat Rate`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/flatrate/type`,
						ID:           "type",
						Label:        `Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `I`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\OfflineShipping\Model\Config\Source\Flatrate
					},

					&config.Field{
						// Path: `carriers/flatrate/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/flatrate/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    91,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/flatrate/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    92,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/flatrate/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `This shipping method is not available. To use this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "tablerate",
				Label:     `Table Rates`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/tablerate/handling_type`,
						ID:           "handling_type",
						Label:        `Calculate Handling Fee`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `F`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Source\HandlingType
					},

					&config.Field{
						// Path: `carriers/tablerate/handling_fee`,
						ID:           "handling_fee",
						Label:        `Handling Fee`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    8,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/tablerate/condition_name`,
						ID:           "condition_name",
						Label:        `Condition`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `package_weight`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\OfflineShipping\Model\Config\Source\Tablerate
					},

					&config.Field{
						// Path: `carriers/tablerate/include_virtual_price`,
						ID:           "include_virtual_price",
						Label:        `Include Virtual Products in Price Calculation`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/tablerate/export`,
						ID:           "export",
						Label:        `Export`,
						Comment:      ``,
						Type:         config.TypeCustom, // @todo: Magento\OfflineShipping\Block\Adminhtml\Form\Field\Export,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/import`,
						ID:           "import",
						Label:        `Import`,
						Comment:      ``,
						Type:         config.TypeCustom, // @todo: Magento\OfflineShipping\Block\Adminhtml\Form\Field\Import,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\OfflineShipping\Model\Config\Backend\Tablerate
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/name`,
						ID:           "name",
						Label:        `Method Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Table Rate`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Best Way`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/tablerate/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/tablerate/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    91,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/tablerate/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    92,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/tablerate/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `This shipping method is not available. To use this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "freeshipping",
				Label:     `Free Shipping`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/freeshipping/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/freeshipping/free_shipping_subtotal`,
						ID:           "free_shipping_subtotal",
						Label:        `Minimum Order Amount`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/freeshipping/name`,
						ID:           "name",
						Label:        `Method Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Free`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/freeshipping/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/freeshipping/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Free Shipping`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `carriers/freeshipping/sallowspecific`,
						ID:           "sallowspecific",
						Label:        `Ship to Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Shipping\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `carriers/freeshipping/specificcountry`,
						ID:           "specificcountry",
						Label:        `Ship to Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    91,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `carriers/freeshipping/showmethod`,
						ID:           "showmethod",
						Label:        `Show Method if Not Applicable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    92,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `carriers/freeshipping/specificerrmsg`,
						ID:           "specificerrmsg",
						Label:        `Displayed Error Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `This shipping method is not available. To use this shipping method, please contact us.`,
						BackendModel: nil,
						SourceModel:  nil,
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
				ID: "flatrate",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/flatrate/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\OfflineShipping\Model\Carrier\Flatrate`,
					},
				},
			},

			&config.Group{
				ID: "tablerate",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/tablerate/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\OfflineShipping\Model\Carrier\Tablerate`,
					},
				},
			},

			&config.Group{
				ID: "freeshipping",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `carriers/freeshipping/cutoff_cost`,
						ID:      "cutoff_cost",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: 50,
					},

					&config.Field{
						// Path: `carriers/freeshipping/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\OfflineShipping\Model\Carrier\Freeshipping`,
					},
				},
			},
		},
	},
)
