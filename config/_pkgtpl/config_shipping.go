// +build ignore

package shipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "shipping",
		Label:     `Shipping Settings`,
		SortOrder: 310,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Shipping::config_shipping
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "origin",
				Label:     `Origin`,
				SortOrder: 1,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: shipping/origin/country_id
						ID:        "country_id",
						Label:     `Country`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `US`,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: shipping/origin/region_id
						ID:        "region_id",
						Label:     `Region/State`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   12,
					},

					&config.Field{
						// Path: shipping/origin/postcode
						ID:        "postcode",
						Label:     `ZIP/Postal Code`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   90034,
					},

					&config.Field{
						// Path: shipping/origin/city
						ID:        "city",
						Label:     `City`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: shipping/origin/street_line1
						ID:        "street_line1",
						Label:     `Street Address`,
						Type:      config.TypeText,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: shipping/origin/street_line2
						ID:        "street_line2",
						Label:     `Street Address Line 2`,
						Type:      config.TypeText,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},
				),
			},

			&config.Group{
				ID:        "shipping_policy",
				Label:     `Shipping Policy Parameters`,
				SortOrder: 120,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: shipping/shipping_policy/enable_shipping_policy
						ID:        "enable_shipping_policy",
						Label:     `Apply custom Shipping Policy`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: shipping/shipping_policy/shipping_policy_content
						ID:        "shipping_policy_content",
						Label:     `Shipping Policy`,
						Type:      config.TypeTextarea,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "carriers",
		Label:     `Shipping Methods`,
		SortOrder: 320,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Shipping::carriers
		Groups:    config.NewGroupSlice(),
	},
)
