// +build ignore

package catalogsearch

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "seo",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/seo/search_terms`,
						ID:           "search_terms",
						Label:        `Popular Search Terms`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},

			&config.Group{
				ID:        "search",
				Label:     `Catalog Search`,
				Comment:   ``,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/search/engine`,
						ID:           "engine",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        nil,
						Default:      `mysql`,
						BackendModel: nil, // Magento\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/search/min_query_length`,
						ID:           "min_query_length",
						Label:        `Minimal Query Length`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      1,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/search/max_query_length`,
						ID:           "max_query_length",
						Label:        `Maximum Query Length`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      128,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/search/use_layered_navigation_count`,
						ID:           "use_layered_navigation_count",
						Label:        `Apply Layered Navigation if Search Results are Less Than`,
						Comment:      `Enter "0" to enable layered navigation for any number of results.`,
						Type:         config.TypeText,
						SortOrder:    25,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      0,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
