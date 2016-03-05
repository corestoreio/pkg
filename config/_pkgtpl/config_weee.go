// +build ignore

package weee

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
			ID: "tax",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "weee",
					Label:     `Fixed Product Taxes`,
					SortOrder: 100,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: tax/weee/enable
							ID:        "enable",
							Label:     `Enable FPT`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: tax/weee/display_list
							ID:        "display_list",
							Label:     `Display Prices In Product Lists`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Weee\Model\Config\Source\Display
						},

						&element.Field{
							// Path: tax/weee/display
							ID:        "display",
							Label:     `Display Prices On Product View Page`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Weee\Model\Config\Source\Display
						},

						&element.Field{
							// Path: tax/weee/display_sales
							ID:        "display_sales",
							Label:     `Display Prices In Sales Modules`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Weee\Model\Config\Source\Display
						},

						&element.Field{
							// Path: tax/weee/display_email
							ID:        "display_email",
							Label:     `Display Prices In Emails`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Weee\Model\Config\Source\Display
						},

						&element.Field{
							// Path: tax/weee/apply_vat
							ID:        "apply_vat",
							Label:     `Apply Tax To FPT`,
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: tax/weee/include_in_subtotal
							ID:        "include_in_subtotal",
							Label:     `Include FPT In Subtotal`,
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID: "sales",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "totals_sort",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/totals_sort/weee
							ID:        "weee",
							Label:     `Fixed Product Tax`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   35,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "sales",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "totals_sort",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/totals_sort/weee_tax
							ID:      `weee_tax`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 35,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "general",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "validator_data",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/validator_data/input_types
							ID:      `input_types`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"weee":"weee"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
