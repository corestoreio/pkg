// +build ignore

package multishipping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "multishipping",
		Label:     "Multishipping Settings",
		SortOrder: 311,
		Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     `Options`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `multishipping/options/checkout_multiple`,
						ID:           "checkout_multiple",
						Label:        `Allow Shipping to Multiple Addresses`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `multishipping/options/checkout_multiple_maximum_qty`,
						ID:           "checkout_multiple_maximum_qty",
						Label:        `Maximum Qty Allowed for Shipping to Multiple Addresses`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      100,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
