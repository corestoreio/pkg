// +build ignore

package salesrule

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "promo",
		Label:     "Promotions",
		SortOrder: 400,
		Scope:     scope.NewPerm(scope.DefaultID),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "auto_generated_coupon_codes",
				Label:     `Auto Generated Specific Coupon Codes`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/length`,
						ID:           "length",
						Label:        `Code Length`,
						Comment:      `Excluding prefix, suffix and separators.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      12,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/format`,
						ID:           "format",
						Label:        `Code Format`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\SalesRule\Model\System\Config\Source\Coupon\Format
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/prefix`,
						ID:           "prefix",
						Label:        `Code Prefix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/suffix`,
						ID:           "suffix",
						Label:        `Code Suffix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/dash`,
						ID:           "dash",
						Label:        `Dash Every X Characters`,
						Comment:      `If empty no separation.`,
						Type:         config.TypeText,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "rss",
		Label:     "",
		SortOrder: 0,
		Scope:     scope.NewPerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "catalog",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     scope.NewPerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/catalog/discounts`,
						ID:           "discounts",
						Label:        `Coupons/Discounts`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    12,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},
)
