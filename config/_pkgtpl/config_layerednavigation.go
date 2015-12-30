// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		SortOrder: 40,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "layered_navigation",
				Label:     `Layered Navigation`,
				SortOrder: 490,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/layered_navigation/display_product_count
						ID:        "display_product_count",
						Label:     `Display Product Count`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/layered_navigation/price_range_calculation
						ID:        "price_range_calculation",
						Label:     `Price Navigation Step Calculation`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `auto`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Step
					},

					&config.Field{
						// Path: catalog/layered_navigation/price_range_step
						ID:        "price_range_step",
						Label:     `Default Price Navigation Step`,
						Type:      config.TypeText,
						SortOrder: 15,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   100,
					},

					&config.Field{
						// Path: catalog/layered_navigation/price_range_max_intervals
						ID:        "price_range_max_intervals",
						Label:     `Maximum Number of Price Intervals`,
						Comment:   element.LongText(`Maximum number of price intervals is 100`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   10,
					},

					&config.Field{
						// Path: catalog/layered_navigation/one_price_interval
						ID:        "one_price_interval",
						Label:     `Display Price Interval as One Price`,
						Comment:   element.LongText(`This setting will be applied when all prices in the specific price interval are equal.`),
						Type:      config.TypeSelect,
						SortOrder: 15,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/layered_navigation/interval_division_limit
						ID:        "interval_division_limit",
						Label:     `Interval Division Limit`,
						Comment:   element.LongText(`Please specify the number of products, that will not be divided into subintervals.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   9,
					},
				),
			},
		),
	},
)
