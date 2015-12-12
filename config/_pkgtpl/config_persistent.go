// +build ignore

package persistent

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "persistent",
		Label:     "Persistent Shopping Cart",
		SortOrder: 500,
		Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     `General Options`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `persistent/options/enabled`,
						ID:           "enabled",
						Label:        `Enable Persistence`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `persistent/options/lifetime`,
						ID:           "lifetime",
						Label:        `Persistence Lifetime (seconds)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      31536000,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `persistent/options/remember_enabled`,
						ID:           "remember_enabled",
						Label:        `Enable "Remember Me"`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `persistent/options/remember_default`,
						ID:           "remember_default",
						Label:        `"Remember Me" Default Value`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `persistent/options/logout_clear`,
						ID:           "logout_clear",
						Label:        `Clear Persistence on Sign Out`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `persistent/options/shopping_cart`,
						ID:           "shopping_cart",
						Label:        `Persist Shopping Cart`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
