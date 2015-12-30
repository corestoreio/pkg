// +build ignore

package googleoptimizer

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "google",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "analytics",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: google/analytics/experiments
						ID:        "experiments",
						Label:     `Enable Content Experiments`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&element.Section{
		ID: "google",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "optimizer",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: google/optimizer/active
						ID:      `active`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: false,
					},
				),
			},
		),
	},
)
