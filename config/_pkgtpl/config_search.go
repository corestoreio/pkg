// +build ignore

package search

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "catalog",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "search",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: catalog/search/engine
						ID:        "engine",
						Label:     `Search Engine`,
						Type:      element.TypeSelect,
						SortOrder: 19,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Search\Model\Adminhtml\System\Config\Source\Engine
					},

					&element.Field{
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
