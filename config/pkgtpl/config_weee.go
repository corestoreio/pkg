// +build ignore

package weee

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "tax",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "weee",
				Label:     `Fixed Product Taxes`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `tax/weee/enable`,
						ID:           "enable",
						Label:        `Enable FPT`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/weee/display_list`,
						ID:           "display_list",
						Label:        `Display Prices In Product Lists`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: `tax/weee/display`,
						ID:           "display",
						Label:        `Display Prices On Product View Page`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: `tax/weee/display_sales`,
						ID:           "display_sales",
						Label:        `Display Prices In Sales Modules`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: `tax/weee/display_email`,
						ID:           "display_email",
						Label:        `Display Prices In Emails`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: `tax/weee/apply_vat`,
						ID:           "apply_vat",
						Label:        `Apply Tax To FPT`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `tax/weee/include_in_subtotal`,
						ID:           "include_in_subtotal",
						Label:        `Include FPT In Subtotal`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "sales",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "totals_sort",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/totals_sort/weee`,
						ID:           "weee",
						Label:        `Fixed Product Tax`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      35,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "sales",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "totals_sort",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/totals_sort/weee_tax`,
						ID:      "weee_tax",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: 35,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "general",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "validator_data",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/validator_data/input_types`,
						ID:      "input_types",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"weee":"weee"}`,
					},
				},
			},
		},
	},
)
