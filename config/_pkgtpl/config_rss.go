// +build ignore

package rss

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
			ID:        "rss",
			Label:     `RSS Feeds`,
			SortOrder: 80,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Rss::rss
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "config",
					Label:     `Rss Config`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: rss/config/active
							ID:        "active",
							Label:     `Enable RSS`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Rss\Model\System\Config\Backend\Links
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Path = NewPath(ConfigStructure)
}
