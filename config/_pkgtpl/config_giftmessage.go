// +build ignore

package giftmessage

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "sales",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "gift_options",
				Label:     `Gift Options`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/gift_options/allow_order`,
						ID:           "allow_order",
						Label:        `Allow Gift Messages on Order Level`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales/gift_options/allow_items`,
						ID:           "allow_items",
						Label:        `Allow Gift Messages for Order Items`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "sales",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "gift_messages",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/gift_messages/allow_items`,
						ID:      "allow_items",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `sales/gift_messages/allow_order`,
						ID:      "allow_order",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},
				},
			},
		},
	},
)
