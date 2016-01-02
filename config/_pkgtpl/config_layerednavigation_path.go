// +build ignore

package layerednavigation

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CatalogLayeredNavigationDisplayProductCount => Display Product Count.
	// Path: catalog/layered_navigation/display_product_count
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationDisplayProductCount model.Bool

	// CatalogLayeredNavigationPriceRangeCalculation => Price Navigation Step Calculation.
	// Path: catalog/layered_navigation/price_range_calculation
	// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Step
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
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogLayeredNavigationOnePriceInterval model.Bool

	// CatalogLayeredNavigationIntervalDivisionLimit => Interval Division Limit.
	// Please specify the number of products, that will not be divided into
	// subintervals.
	// Path: catalog/layered_navigation/interval_division_limit
	CatalogLayeredNavigationIntervalDivisionLimit model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogLayeredNavigationDisplayProductCount = model.NewBool(`catalog/layered_navigation/display_product_count`, model.WithPkgCfg(pkgCfg))
	pp.CatalogLayeredNavigationPriceRangeCalculation = model.NewStr(`catalog/layered_navigation/price_range_calculation`, model.WithPkgCfg(pkgCfg))
	pp.CatalogLayeredNavigationPriceRangeStep = model.NewStr(`catalog/layered_navigation/price_range_step`, model.WithPkgCfg(pkgCfg))
	pp.CatalogLayeredNavigationPriceRangeMaxIntervals = model.NewStr(`catalog/layered_navigation/price_range_max_intervals`, model.WithPkgCfg(pkgCfg))
	pp.CatalogLayeredNavigationOnePriceInterval = model.NewBool(`catalog/layered_navigation/one_price_interval`, model.WithPkgCfg(pkgCfg))
	pp.CatalogLayeredNavigationIntervalDivisionLimit = model.NewStr(`catalog/layered_navigation/interval_division_limit`, model.WithPkgCfg(pkgCfg))

	return pp
}
