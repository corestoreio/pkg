// +build ignore

package weee

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "tax",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "weee",
				Label:     `Fixed Product Taxes`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: tax/weee/enable
						ID:        "enable",
						Label:     `Enable FPT`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: tax/weee/display_list
						ID:        "display_list",
						Label:     `Display Prices In Product Lists`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: tax/weee/display
						ID:        "display",
						Label:     `Display Prices On Product View Page`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: tax/weee/display_sales
						ID:        "display_sales",
						Label:     `Display Prices In Sales Modules`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: tax/weee/display_email
						ID:        "display_email",
						Label:     `Display Prices In Emails`,
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Weee\Model\Config\Source\Display
					},

					&config.Field{
						// Path: tax/weee/apply_vat
						ID:        "apply_vat",
						Label:     `Apply Tax To FPT`,
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: tax/weee/include_in_subtotal
						ID:        "include_in_subtotal",
						Label:     `Include FPT In Subtotal`,
						Type:      config.TypeSelect,
						SortOrder: 70,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID: "sales",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "totals_sort",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sales/totals_sort/weee
						ID:        "weee",
						Label:     `Fixed Product Tax`,
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   35,
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "sales",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "totals_sort",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sales/totals_sort/weee_tax
						ID:      `weee_tax`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 35,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "general",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "validator_data",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/validator_data/input_types
						ID:      `input_types`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"weee":"weee"}`,
					},
				),
			},
		),
	},
)
