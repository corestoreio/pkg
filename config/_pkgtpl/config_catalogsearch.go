// +build ignore

package catalogsearch

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
					ID: "seo",
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/seo/search_terms
							ID:        "search_terms",
							Label:     `Popular Search Terms`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},

				element.Group{
					ID:        "search",
					Label:     `Catalog Search`,
					SortOrder: 500,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/search/engine
							ID:      "engine",
							Type:    element.Type,
							Visible: element.VisibleYes,
							Default: `mysql`,
							// BackendModel: Magento\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
						},

						element.Field{
							// Path: catalog/search/min_query_length
							ID:        "min_query_length",
							Label:     `Minimal Query Length`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   1,
						},

						element.Field{
							// Path: catalog/search/max_query_length
							ID:        "max_query_length",
							Label:     `Maximum Query Length`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   128,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
