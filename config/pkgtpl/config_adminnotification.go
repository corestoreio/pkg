// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "system",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "adminnotification",
				Label:     `Notifications`,
				Comment:   ``,
				SortOrder: 250,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/adminnotification/use_https`,
						ID:           "use_https",
						Label:        `Use HTTPS to Get Feed`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `system/adminnotification/frequency`,
						ID:           "frequency",
						Label:        `Update Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\AdminNotification\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `system/adminnotification/last_update`,
						ID:           "last_update",
						Label:        `Last Update`,
						Comment:      ``,
						Type:         config.TypeLabel,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "adminnotification",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/adminnotification/feed_url`,
						ID:      "feed_url",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `notifications.magentocommerce.com/magento2/community/notifications.rss`,
					},

					&config.Field{
						// Path: `system/adminnotification/popup_url`,
						ID:      "popup_url",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `widgets.magentocommerce.com/notificationPopup`,
					},

					&config.Field{
						// Path: `system/adminnotification/severity_icons_url`,
						ID:      "severity_icons_url",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `widgets.magentocommerce.com/%s/%s.gif`,
					},
				},
			},
		},
	},
)
