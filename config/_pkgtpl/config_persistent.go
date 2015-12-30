// +build ignore

package persistent

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "persistent",
		Label:     `Persistent Shopping Cart`,
		SortOrder: 500,
		Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
		Resource:  0, // Otnegam_Persistent::persistent
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "options",
				Label:     `General Options`,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: persistent/options/enabled
						ID:        "enabled",
						Label:     `Enable Persistence`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: persistent/options/lifetime
						ID:        "lifetime",
						Label:     `Persistence Lifetime (seconds)`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   31536000,
					},

					&config.Field{
						// Path: persistent/options/remember_enabled
						ID:        "remember_enabled",
						Label:     `Enable "Remember Me"`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: persistent/options/remember_default
						ID:        "remember_default",
						Label:     `"Remember Me" Default Value`,
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: persistent/options/logout_clear
						ID:        "logout_clear",
						Label:     `Clear Persistence on Sign Out`,
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: persistent/options/shopping_cart
						ID:        "shopping_cart",
						Label:     `Persist Shopping Cart`,
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
