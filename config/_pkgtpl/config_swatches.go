// +build ignore

package swatches

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 40,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "frontend",
				Label:     ``,
				Comment:   ``,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/frontend/swatches_per_product`,
						ID:           "swatches_per_product",
						Label:        `Swatches per Product`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    300,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      16,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
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
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"swatch_visual":"swatch_visual","swatch_text":"swatch_text"}`,
					},
				},
			},
		},
	},
)
