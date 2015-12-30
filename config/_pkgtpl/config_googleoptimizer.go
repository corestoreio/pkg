// +build ignore

package googleoptimizer

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "google",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "analytics",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: google/analytics/experiments
						ID:        "experiments",
						Label:     `Enable Content Experiments`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "google",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "optimizer",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: google/optimizer/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},
				),
			},
		),
	},
)
