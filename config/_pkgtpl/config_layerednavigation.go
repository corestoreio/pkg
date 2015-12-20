// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 40,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "layered_navigation",
				Label:     `Layered Navigation`,
				Comment:   ``,
				SortOrder: 490,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/layered_navigation/display_product_count`,
						ID:           "display_product_count",
						Label:        `Display Product Count`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/layered_navigation/price_range_calculation`,
						ID:           "price_range_calculation",
						Label:        `Price Navigation Step Calculation`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `auto`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Price\Step
					},

					&config.Field{
						// Path: `catalog/layered_navigation/price_range_step`,
						ID:           "price_range_step",
						Label:        `Default Price Navigation Step`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    15,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      100,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/layered_navigation/price_range_max_intervals`,
						ID:           "price_range_max_intervals",
						Label:        `Maximum Number of Price Intervals`,
						Comment:      `Maximum number of price intervals is 100`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      10,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/layered_navigation/one_price_interval`,
						ID:           "one_price_interval",
						Label:        `Display Price Interval as One Price`,
						Comment:      `This setting will be applied when all prices in the specific price interval are equal.`,
						Type:         config.TypeSelect,
						SortOrder:    15,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/layered_navigation/interval_division_limit`,
						ID:           "interval_division_limit",
						Label:        `Interval Division Limit`,
						Comment:      `Please specify the number of products, that will not be divided into subintervals.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      9,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
