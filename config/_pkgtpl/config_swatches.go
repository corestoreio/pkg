// +build ignore

package swatches

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		SortOrder: 40,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "frontend",
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/frontend/swatches_per_product
						ID:        "swatches_per_product",
						Label:     `Swatches per Product`,
						Type:      config.TypeText,
						SortOrder: 300,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   16,
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
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
						Default: `{"swatch_visual":"swatch_visual","swatch_text":"swatch_text"}`,
					},
				),
			},
		),
	},
)
