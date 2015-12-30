// +build ignore

package salesrule

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID:        "promo",
		Label:     `Promotions`,
		SortOrder: 400,
		Scope:     scope.NewPerm(scope.DefaultID),
		Resource:  0, // Otnegam_SalesRule::config_promo
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "auto_generated_coupon_codes",
				Label:     `Auto Generated Specific Coupon Codes`,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: promo/auto_generated_coupon_codes/length
						ID:        "length",
						Label:     `Code Length`,
						Comment:   element.LongText(`Excluding prefix, suffix and separators.`),
						Type:      element.TypeText,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   12,
					},

					&element.Field{
						// Path: promo/auto_generated_coupon_codes/format
						ID:        "format",
						Label:     `Code Format`,
						Type:      element.TypeSelect,
						SortOrder: 20,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\SalesRule\Model\System\Config\Source\Coupon\Format
					},

					&element.Field{
						// Path: promo/auto_generated_coupon_codes/prefix
						ID:        "prefix",
						Label:     `Code Prefix`,
						Type:      element.TypeText,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&element.Field{
						// Path: promo/auto_generated_coupon_codes/suffix
						ID:        "suffix",
						Label:     `Code Suffix`,
						Type:      element.TypeText,
						SortOrder: 40,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&element.Field{
						// Path: promo/auto_generated_coupon_codes/dash
						ID:        "dash",
						Label:     `Dash Every X Characters`,
						Comment:   element.LongText(`If empty no separation.`),
						Type:      element.TypeText,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},
		),
	},
	&element.Section{
		ID: "rss",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "catalog",
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: rss/catalog/discounts
						ID:        "discounts",
						Label:     `Coupons/Discounts`,
						Type:      element.TypeSelect,
						SortOrder: 12,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},
		),
	},
)
