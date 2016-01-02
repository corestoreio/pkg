// +build ignore

package shipping

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package.
// Used in frontend and backend. See init() for details.
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(
		&element.Section{
			ID:        "shipping",
			Label:     `Shipping Settings`,
			SortOrder: 310,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Shipping::config_shipping
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "origin",
					Label:     `Origin`,
					SortOrder: 1,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: shipping/origin/country_id
							ID:        "country_id",
							Label:     `Country`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `US`,
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: shipping/origin/region_id
							ID:        "region_id",
							Label:     `Region/State`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   12,
						},

						&element.Field{
							// Path: shipping/origin/postcode
							ID:        "postcode",
							Label:     `ZIP/Postal Code`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   90034,
						},

						&element.Field{
							// Path: shipping/origin/city
							ID:        "city",
							Label:     `City`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: shipping/origin/street_line1
							ID:        "street_line1",
							Label:     `Street Address`,
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: shipping/origin/street_line2
							ID:        "street_line2",
							Label:     `Street Address Line 2`,
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},
					),
				},

				&element.Group{
					ID:        "shipping_policy",
					Label:     `Shipping Policy Parameters`,
					SortOrder: 120,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: shipping/shipping_policy/enable_shipping_policy
							ID:        "enable_shipping_policy",
							Label:     `Apply custom Shipping Policy`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: shipping/shipping_policy/shipping_policy_content
							ID:        "shipping_policy_content",
							Label:     `Shipping Policy`,
							Type:      element.TypeTextarea,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "carriers",
			Label:     `Shipping Methods`,
			SortOrder: 320,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Shipping::carriers
			Groups:    element.NewGroupSlice(),
		},
	)
	Path = NewPath(PackageConfiguration)
}
