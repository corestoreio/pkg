// +build ignore

package swatches

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
			ID:        "catalog",
			SortOrder: 40,
			Scope:     scope.PermAll,
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "frontend",
					SortOrder: 100,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/frontend/swatches_per_product
							ID:        "swatches_per_product",
							Label:     `Swatches per Product`,
							Type:      element.TypeText,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   16,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
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
							Default: `{"swatch_visual":"swatch_visual","swatch_text":"swatch_text"}`,
						},
					),
				},
			),
		},
	)
}
