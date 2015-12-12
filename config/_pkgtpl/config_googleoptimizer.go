// +build ignore

package googleoptimizer

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID: "google",
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "analytics",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/analytics/experiments`,
						ID:           "experiments",
						Label:        `Enable Content Experiments`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "google",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "optimizer",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/optimizer/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},
				},
			},
		},
	},
)
