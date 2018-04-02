// +build ignore

package checkoutagreements

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
			ID: "checkout",
			Groups: element.MakeGroups(
				element.Group{
					ID: "options",
					Fields: element.MakeFields(
						element.Field{
							// Path: checkout/options/enable_agreements
							ID:        "enable_agreements",
							Label:     `Enable Terms and Conditions`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
