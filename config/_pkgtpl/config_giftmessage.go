// +build ignore

package giftmessage

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "sales",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "gift_options",
				Label:     `Gift Options`,
				SortOrder: 100,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sales/gift_options/allow_order
						ID:        "allow_order",
						Label:     `Allow Gift Messages on Order Level`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: sales/gift_options/allow_items
						ID:        "allow_items",
						Label:     `Allow Gift Messages for Order Items`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "sales",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "gift_messages",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sales/gift_messages/allow_items
						ID:      `allow_items`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: sales/gift_messages/allow_order
						ID:      `allow_order`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},
				),
			},
		),
	},
)
