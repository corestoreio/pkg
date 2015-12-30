// +build ignore

package salesrule

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "promo",
		Label:     `Promotions`,
		SortOrder: 400,
		Scope:     scope.NewPerm(scope.DefaultID),
		Resource:  0, // Otnegam_SalesRule::config_promo
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "auto_generated_coupon_codes",
				Label:     `Auto Generated Specific Coupon Codes`,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: promo/auto_generated_coupon_codes/length
						ID:        "length",
						Label:     `Code Length`,
						Comment:   element.LongText(`Excluding prefix, suffix and separators.`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   12,
					},

					&config.Field{
						// Path: promo/auto_generated_coupon_codes/format
						ID:        "format",
						Label:     `Code Format`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\SalesRule\Model\System\Config\Source\Coupon\Format
					},

					&config.Field{
						// Path: promo/auto_generated_coupon_codes/prefix
						ID:        "prefix",
						Label:     `Code Prefix`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: promo/auto_generated_coupon_codes/suffix
						ID:        "suffix",
						Label:     `Code Suffix`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: promo/auto_generated_coupon_codes/dash
						ID:        "dash",
						Label:     `Dash Every X Characters`,
						Comment:   element.LongText(`If empty no separation.`),
						Type:      config.TypeText,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},
		),
	},
	&config.Section{
		ID: "rss",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "catalog",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: rss/catalog/discounts
						ID:        "discounts",
						Label:     `Coupons/Discounts`,
						Type:      config.TypeSelect,
						SortOrder: 12,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},
		),
	},
)
