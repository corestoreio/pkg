// +build ignore

package multishipping

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
			ID:        "multishipping",
			Label:     `Multishipping Settings`,
			SortOrder: 311,
			Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
			Resource:  0, // Otnegam_Multishipping::config_multishipping
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "options",
					Label:     `Options`,
					SortOrder: 2,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: multishipping/options/checkout_multiple
							ID:        "checkout_multiple",
							Label:     `Allow Shipping to Multiple Addresses`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: multishipping/options/checkout_multiple_maximum_qty
							ID:        "checkout_multiple_maximum_qty",
							Label:     `Maximum Qty Allowed for Shipping to Multiple Addresses`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   100,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
