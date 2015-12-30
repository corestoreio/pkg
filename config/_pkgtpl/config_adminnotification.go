// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "system",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "adminnotification",
				Label:     `Notifications`,
				SortOrder: 250,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: system/adminnotification/use_https
						ID:        "use_https",
						Label:     `Use HTTPS to Get Feed`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: system/adminnotification/frequency
						ID:        "frequency",
						Label:     `Update Frequency`,
						Type:      element.TypeSelect,
						SortOrder: 2,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\AdminNotification\Model\Config\Source\Frequency
					},

					&element.Field{
						// Path: system/adminnotification/last_update
						ID:        "last_update",
						Label:     `Last Update`,
						Type:      element.TypeLabel,
						SortOrder: 3,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&element.Section{
		ID: "system",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "adminnotification",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: system/adminnotification/feed_url
						ID:      `feed_url`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: `notifications.magentocommerce.com/magento2/community/notifications.rss`,
					},

					&element.Field{
						// Path: system/adminnotification/popup_url
						ID:      `popup_url`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: `widgets.magentocommerce.com/notificationPopup`,
					},

					&element.Field{
						// Path: system/adminnotification/severity_icons_url
						ID:      `severity_icons_url`,
						Type:    element.TypeHidden,
						Visible: element.VisibleNo,
						Default: `widgets.magentocommerce.com/%s/%s.gif`,
					},
				),
			},
		),
	},
)
