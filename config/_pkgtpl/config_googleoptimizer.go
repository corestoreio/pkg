// +build ignore

package googleoptimizer

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
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
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
	Backend = NewBackend(ConfigStructure)
}
