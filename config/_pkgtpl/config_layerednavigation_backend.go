// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CatalogLayeredNavigationDisplayProductCount => Display Product Count.
	// Path: catalog/layered_navigation/display_product_count
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationDisplayProductCount cfgmodel.Bool

	// CatalogLayeredNavigationPriceRangeCalculation => Price Navigation Step Calculation.
	// Path: catalog/layered_navigation/price_range_calculation
	// SourceModel: Magento\Catalog\Model\Config\Source\Price\Step
	CatalogLayeredNavigationPriceRangeCalculation cfgmodel.Str

	// CatalogLayeredNavigationPriceRangeStep => Default Price Navigation Step.
	// Path: catalog/layered_navigation/price_range_step
	CatalogLayeredNavigationPriceRangeStep cfgmodel.Str

	// CatalogLayeredNavigationPriceRangeMaxIntervals => Maximum Number of Price Intervals.
	// Maximum number of price intervals is 100
	// Path: catalog/layered_navigation/price_range_max_intervals
	CatalogLayeredNavigationPriceRangeMaxIntervals cfgmodel.Str

	// CatalogLayeredNavigationOnePriceInterval => Display Price Interval as One Price.
	// This setting will be applied when all prices in the specific price interval
	// are equal.
	// Path: catalog/layered_navigation/one_price_interval
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationOnePriceInterval cfgmodel.Bool

	// CatalogLayeredNavigationIntervalDivisionLimit => Interval Division Limit.
	// Please specify the number of products, that will not be divided into
	// subintervals.
	// Path: catalog/layered_navigation/interval_division_limit
	CatalogLayeredNavigationIntervalDivisionLimit cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogLayeredNavigationDisplayProductCount = cfgmodel.NewBool(`catalog/layered_navigation/display_product_count`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeCalculation = cfgmodel.NewStr(`catalog/layered_navigation/price_range_calculation`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeStep = cfgmodel.NewStr(`catalog/layered_navigation/price_range_step`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeMaxIntervals = cfgmodel.NewStr(`catalog/layered_navigation/price_range_max_intervals`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationOnePriceInterval = cfgmodel.NewBool(`catalog/layered_navigation/one_price_interval`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationIntervalDivisionLimit = cfgmodel.NewStr(`catalog/layered_navigation/interval_division_limit`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
