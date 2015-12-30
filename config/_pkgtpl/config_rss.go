// +build ignore

package rss

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "rss",
		Label:     `RSS Feeds`,
		SortOrder: 80,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Rss::rss
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "config",
				Label:     `Rss Config`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: rss/config/active
						ID:        "active",
						Label:     `Enable RSS`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Rss\Model\System\Config\Backend\Links
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},
		),
	},
)
