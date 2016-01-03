// +build ignore

package catalogsearch

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
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "seo",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/seo/search_terms
							ID:        "search_terms",
							Label:     `Popular Search Terms`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},
					),
				},

				&element.Group{
					ID:        "search",
					Label:     `Catalog Search`,
					SortOrder: 500,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/search/engine
							ID:      "engine",
							Type:    element.Type,
							Visible: element.VisibleYes,
							Default: `mysql`,
							// BackendModel: Otnegam\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
						},

						&element.Field{
							// Path: catalog/search/min_query_length
							ID:        "min_query_length",
							Label:     `Minimal Query Length`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   1,
						},

						&element.Field{
							// Path: catalog/search/max_query_length
							ID:        "max_query_length",
							Label:     `Maximum Query Length`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   128,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
