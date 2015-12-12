// +build ignore

package search

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "search",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/search/engine`,
						ID:           "engine",
						Label:        `Search Engine`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    19,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Search\Model\Adminhtml\System\Config\Source\Engine
					},

					&config.Field{
						// Path: `catalog/search/search_type`,
						ID:           "search_type",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        nil,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/search/use_layered_navigation_count`,
						ID:           "use_layered_navigation_count",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        nil,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
