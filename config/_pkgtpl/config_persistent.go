// +build ignore

package persistent

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "persistent",
			Label:     `Persistent Shopping Cart`,
			SortOrder: 500,
			Scopes:    scope.PermWebsite,
			Resource:  0, // Magento_Persistent::persistent
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "options",
					Label:     `General Options`,
					SortOrder: 10,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: persistent/options/enabled
							ID:        "enabled",
							Label:     `Enable Persistence`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: persistent/options/lifetime
							ID:        "lifetime",
							Label:     `Persistence Lifetime (seconds)`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   31536000,
						},

						element.Field{
							// Path: persistent/options/remember_enabled
							ID:        "remember_enabled",
							Label:     `Enable "Remember Me"`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: persistent/options/remember_default
							ID:        "remember_default",
							Label:     `"Remember Me" Default Value`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: persistent/options/logout_clear
							ID:        "logout_clear",
							Label:     `Clear Persistence on Sign Out`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: persistent/options/shopping_cart
							ID:        "shopping_cart",
							Label:     `Persist Shopping Cart`,
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
