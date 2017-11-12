// +build ignore

package reports

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CatalogRecentlyProductsScope => Show for Current.
	// Path: catalog/recently_products/scope
	// SourceModel: Magento\Config\Model\Config\Source\Reports\Scope
	CatalogRecentlyProductsScope cfgmodel.Str

	// CatalogRecentlyProductsViewedCount => Default Recently Viewed Products Count.
	// Path: catalog/recently_products/viewed_count
	CatalogRecentlyProductsViewedCount cfgmodel.Str

	// CatalogRecentlyProductsComparedCount => Default Recently Compared Products Count.
	// Path: catalog/recently_products/compared_count
	CatalogRecentlyProductsComparedCount cfgmodel.Str

	// ReportsDashboardYtdStart => Year-To-Date Starts.
	// Path: reports/dashboard/ytd_start
	ReportsDashboardYtdStart cfgmodel.Str

	// ReportsDashboardMtdStart => Current Month Starts.
	// Select day of the month.
	// Path: reports/dashboard/mtd_start
	ReportsDashboardMtdStart cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogRecentlyProductsScope = cfgmodel.NewStr(`catalog/recently_products/scope`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogRecentlyProductsViewedCount = cfgmodel.NewStr(`catalog/recently_products/viewed_count`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogRecentlyProductsComparedCount = cfgmodel.NewStr(`catalog/recently_products/compared_count`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ReportsDashboardYtdStart = cfgmodel.NewStr(`reports/dashboard/ytd_start`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ReportsDashboardMtdStart = cfgmodel.NewStr(`reports/dashboard/mtd_start`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
