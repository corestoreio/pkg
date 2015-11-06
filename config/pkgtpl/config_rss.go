// +build ignore

package rss

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "rss",
		Label:     "RSS Feeds",
		SortOrder: 80,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "config",
				Label:     `Rss Config`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/config/active`,
						ID:           "active",
						Label:        `Enable RSS`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Rss\Model\System\Config\Backend\Links
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},
)
