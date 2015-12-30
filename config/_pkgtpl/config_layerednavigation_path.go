// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogLayeredNavigationDisplayProductCount => Display Product Count.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogLayeredNavigationDisplayProductCount = model.NewBool(`catalog/layered_navigation/display_product_count`)

// PathCatalogLayeredNavigationPriceRangeCalculation => Price Navigation Step Calculation.
// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Step
var PathCatalogLayeredNavigationPriceRangeCalculation = model.NewStr(`catalog/layered_navigation/price_range_calculation`)

// PathCatalogLayeredNavigationPriceRangeStep => Default Price Navigation Step.
var PathCatalogLayeredNavigationPriceRangeStep = model.NewStr(`catalog/layered_navigation/price_range_step`)

// PathCatalogLayeredNavigationPriceRangeMaxIntervals => Maximum Number of Price Intervals.
// Maximum number of price intervals is 100
var PathCatalogLayeredNavigationPriceRangeMaxIntervals = model.NewStr(`catalog/layered_navigation/price_range_max_intervals`)

// PathCatalogLayeredNavigationOnePriceInterval => Display Price Interval as One Price.
// This setting will be applied when all prices in the specific price interval
// are equal.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogLayeredNavigationOnePriceInterval = model.NewBool(`catalog/layered_navigation/one_price_interval`)

// PathCatalogLayeredNavigationIntervalDivisionLimit => Interval Division Limit.
// Please specify the number of products, that will not be divided into
// subintervals.
var PathCatalogLayeredNavigationIntervalDivisionLimit = model.NewStr(`catalog/layered_navigation/interval_division_limit`)
