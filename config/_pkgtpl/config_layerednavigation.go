// +build ignore

package layerednavigation

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
			ID:        "catalog",
			SortOrder: 40,
			Scope:     scope.PermAll,
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "layered_navigation",
					Label:     `Layered Navigation`,
					SortOrder: 490,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/layered_navigation/display_product_count
							ID:        "display_product_count",
							Label:     `Display Product Count`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/layered_navigation/price_range_calculation
							ID:        "price_range_calculation",
							Label:     `Price Navigation Step Calculation`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `auto`,
							// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Step
						},

						&element.Field{
							// Path: catalog/layered_navigation/price_range_step
							ID:        "price_range_step",
							Label:     `Default Price Navigation Step`,
							Type:      element.TypeText,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   100,
						},

						&element.Field{
							// Path: catalog/layered_navigation/price_range_max_intervals
							ID:        "price_range_max_intervals",
							Label:     `Maximum Number of Price Intervals`,
							Comment:   element.LongText(`Maximum number of price intervals is 100`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   10,
						},

						&element.Field{
							// Path: catalog/layered_navigation/one_price_interval
							ID:        "one_price_interval",
							Label:     `Display Price Interval as One Price`,
							Comment:   element.LongText(`This setting will be applied when all prices in the specific price interval are equal.`),
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/layered_navigation/interval_division_limit
							ID:        "interval_division_limit",
							Label:     `Interval Division Limit`,
							Comment:   element.LongText(`Please specify the number of products, that will not be divided into subintervals.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   9,
						},
					),
				},
			),
		},
	)
	Path = NewPath(PackageConfiguration)
}
