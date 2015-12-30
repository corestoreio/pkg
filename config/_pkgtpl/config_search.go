// +build ignore

package search

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
				ID: "search",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/search/engine
						ID:        "engine",
						Label:     `Search Engine`,
						Type:      config.TypeSelect,
						SortOrder: 19,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Search\Model\Adminhtml\System\Config\Source\Engine
					},

					&config.Field{
						// Path: catalog/search/search_type
						ID:      "search_type",
						Type:    config.Type,
						Visible: config.VisibleYes,
					},
				),
			},
		),
	},
)
