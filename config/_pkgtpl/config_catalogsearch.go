// +build ignore

package catalogsearch

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "seo",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/seo/search_terms
						ID:        "search_terms",
						Label:     `Popular Search Terms`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},

			&config.Group{
				ID:        "search",
				Label:     `Catalog Search`,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/search/engine
						ID:      "engine",
						Type:    config.Type,
						Visible: config.VisibleYes,
						Default: `mysql`,
						// BackendModel: Otnegam\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
					},

					&config.Field{
						// Path: catalog/search/min_query_length
						ID:        "min_query_length",
						Label:     `Minimal Query Length`,
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   1,
					},

					&config.Field{
						// Path: catalog/search/max_query_length
						ID:        "max_query_length",
						Label:     `Maximum Query Length`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   128,
					},
				),
			},
		),
	},
)
