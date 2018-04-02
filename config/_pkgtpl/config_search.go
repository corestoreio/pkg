// +build ignore

package search

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
			ID: "catalog",
			Groups: element.MakeGroups(
				element.Group{
					ID: "search",
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/search/engine
							ID:        "engine",
							Label:     `Search Engine`,
							Type:      element.TypeSelect,
							SortOrder: 19,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Search\Model\Adminhtml\System\Config\Source\Engine
						},

						element.Field{
							// Path: catalog/search/search_type
							ID:      "search_type",
							Type:    element.Type,
							Visible: element.VisibleYes,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
