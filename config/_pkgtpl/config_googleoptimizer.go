// +build ignore

package googleoptimizer

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID: "google",
			Groups: element.MakeGroups(
				element.Group{
					ID: "analytics",
					Fields: element.MakeFields(
						element.Field{
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
		element.Section{
			ID: "google",
			Groups: element.MakeGroups(
				element.Group{
					ID: "optimizer",
					Fields: element.MakeFields(
						element.Field{
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
