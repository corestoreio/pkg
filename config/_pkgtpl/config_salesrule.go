// +build ignore

package salesrule

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        "promo",
			Label:     `Promotions`,
			SortOrder: 400,
			Scopes:    scope.PermDefault,
			Resource:  0, // Magento_SalesRule::config_promo
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "auto_generated_coupon_codes",
					Label:     `Auto Generated Specific Coupon Codes`,
					SortOrder: 10,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: promo/auto_generated_coupon_codes/length
							ID:        "length",
							Label:     `Code Length`,
							Comment:   text.Long(`Excluding prefix, suffix and separators.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   12,
						},

						&element.Field{
							// Path: promo/auto_generated_coupon_codes/format
							ID:        "format",
							Label:     `Code Format`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   true,
							// SourceModel: Magento\SalesRule\Model\System\Config\Source\Coupon\Format
						},

						&element.Field{
							// Path: promo/auto_generated_coupon_codes/prefix
							ID:        "prefix",
							Label:     `Code Prefix`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						&element.Field{
							// Path: promo/auto_generated_coupon_codes/suffix
							ID:        "suffix",
							Label:     `Code Suffix`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						&element.Field{
							// Path: promo/auto_generated_coupon_codes/dash
							ID:        "dash",
							Label:     `Dash Every X Characters`,
							Comment:   text.Long(`If empty no separation.`),
							Type:      element.TypeText,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
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
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
