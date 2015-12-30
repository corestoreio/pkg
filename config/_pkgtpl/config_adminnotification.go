// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "adminnotification",
				Label:     `Notifications`,
				SortOrder: 250,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/adminnotification/use_https
						ID:        "use_https",
						Label:     `Use HTTPS to Get Feed`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: system/adminnotification/frequency
						ID:        "frequency",
						Label:     `Update Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\AdminNotification\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: system/adminnotification/last_update
						ID:        "last_update",
						Label:     `Last Update`,
						Type:      config.TypeLabel,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "adminnotification",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/adminnotification/feed_url
						ID:      `feed_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `notifications.magentocommerce.com/magento2/community/notifications.rss`,
					},

					&config.Field{
						// Path: system/adminnotification/popup_url
						ID:      `popup_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `widgets.magentocommerce.com/notificationPopup`,
					},

					&config.Field{
						// Path: system/adminnotification/severity_icons_url
						ID:      `severity_icons_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `widgets.magentocommerce.com/%s/%s.gif`,
					},
				),
			},
		),
	},
)
