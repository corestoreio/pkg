// +build ignore

package multishipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "multishipping",
		Label:     `Multishipping Settings`,
		SortOrder: 311,
		Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
		Resource:  0, // Otnegam_Multishipping::config_multishipping
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "options",
				Label:     `Options`,
				SortOrder: 2,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: multishipping/options/checkout_multiple
						ID:        "checkout_multiple",
						Label:     `Allow Shipping to Multiple Addresses`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: multishipping/options/checkout_multiple_maximum_qty
						ID:        "checkout_multiple_maximum_qty",
						Label:     `Maximum Qty Allowed for Shipping to Multiple Addresses`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   100,
					},
				),
			},
		),
	},
)
