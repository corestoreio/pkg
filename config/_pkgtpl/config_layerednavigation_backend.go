// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CatalogLayeredNavigationDisplayProductCount => Display Product Count.
	// Path: catalog/layered_navigation/display_product_count
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationDisplayProductCount model.Bool

	// CatalogLayeredNavigationPriceRangeCalculation => Price Navigation Step Calculation.
	// Path: catalog/layered_navigation/price_range_calculation
	// SourceModel: Magento\Catalog\Model\Config\Source\Price\Step
	CatalogLayeredNavigationPriceRangeCalculation model.Str

	// CatalogLayeredNavigationPriceRangeStep => Default Price Navigation Step.
	// Path: catalog/layered_navigation/price_range_step
	CatalogLayeredNavigationPriceRangeStep model.Str

	// CatalogLayeredNavigationPriceRangeMaxIntervals => Maximum Number of Price Intervals.
	// Maximum number of price intervals is 100
	// Path: catalog/layered_navigation/price_range_max_intervals
	CatalogLayeredNavigationPriceRangeMaxIntervals model.Str

	// CatalogLayeredNavigationOnePriceInterval => Display Price Interval as One Price.
	// This setting will be applied when all prices in the specific price interval
	// are equal.
	// Path: catalog/layered_navigation/one_price_interval
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationOnePriceInterval model.Bool

	// CatalogLayeredNavigationIntervalDivisionLimit => Interval Division Limit.
	// Please specify the number of products, that will not be divided into
	// subintervals.
	// Path: catalog/layered_navigation/interval_division_limit
	CatalogLayeredNavigationIntervalDivisionLimit model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogLayeredNavigationDisplayProductCount = model.NewBool(`catalog/layered_navigation/display_product_count`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeCalculation = model.NewStr(`catalog/layered_navigation/price_range_calculation`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeStep = model.NewStr(`catalog/layered_navigation/price_range_step`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationPriceRangeMaxIntervals = model.NewStr(`catalog/layered_navigation/price_range_max_intervals`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationOnePriceInterval = model.NewBool(`catalog/layered_navigation/one_price_interval`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogLayeredNavigationIntervalDivisionLimit = model.NewStr(`catalog/layered_navigation/interval_division_limit`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
