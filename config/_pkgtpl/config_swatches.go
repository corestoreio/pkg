// +build ignore

package swatches

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "catalog",
			SortOrder: 40,
			Scopes:    scope.PermStore,
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "frontend",
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: catalog/frontend/swatches_per_product
							ID:        "swatches_per_product",
							Label:     `Swatches per Product`,
							Type:      element.TypeText,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   16,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "general",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "validator_data",
					Fields: element.NewFieldSlice(
						element.Field{
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
	Backend = NewBackend(ConfigStructure)
}
