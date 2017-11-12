// +build ignore

package rss

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "rss",
			Label:     `RSS Feeds`,
			SortOrder: 80,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Rss::rss
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "config",
					Label:     `Rss Config`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: rss/config/active
							ID:        "active",
							Label:     `Enable RSS`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Rss\Model\System\Config\Backend\Links
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
