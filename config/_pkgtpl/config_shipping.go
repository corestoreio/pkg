// +build ignore

package shipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "shipping",
		Label:     "Shipping Settings",
		SortOrder: 310,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "origin",
				Label:     `Origin`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `shipping/origin/country_id`,
						ID:           "country_id",
						Label:        `Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `US`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `shipping/origin/region_id`,
						ID:           "region_id",
						Label:        `Region/State`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      12,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/postcode`,
						ID:           "postcode",
						Label:        `ZIP/Postal Code`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      90034,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/city`,
						ID:           "city",
						Label:        `City`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/street_line1`,
						ID:           "street_line1",
						Label:        `Street Address`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `shipping/origin/street_line2`,
						ID:           "street_line2",
						Label:        `Street Address Line 2`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "shipping_policy",
				Label:     `Shipping Policy Parameters`,
				Comment:   ``,
				SortOrder: 120,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `shipping/shipping_policy/enable_shipping_policy`,
						ID:           "enable_shipping_policy",
						Label:        `Apply custom Shipping Policy`,
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
						// Path: `shipping/shipping_policy/shipping_policy_content`,
						ID:           "shipping_policy_content",
						Label:        `Shipping Policy`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "carriers",
		Label:     "Shipping Methods",
		SortOrder: 320,
		Scope:     scope.PermAll,
		Groups:    config.GroupSlice{},
	},
)
