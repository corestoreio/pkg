// +build ignore

package reports

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
	// CatalogRecentlyProductsScope => Show for Current.
	// Path: catalog/recently_products/scope
	// SourceModel: Otnegam\Config\Model\Config\Source\Reports\Scope
	CatalogRecentlyProductsScope model.Str

	// CatalogRecentlyProductsViewedCount => Default Recently Viewed Products Count.
	// Path: catalog/recently_products/viewed_count
	CatalogRecentlyProductsViewedCount model.Str

	// CatalogRecentlyProductsComparedCount => Default Recently Compared Products Count.
	// Path: catalog/recently_products/compared_count
	CatalogRecentlyProductsComparedCount model.Str

	// ReportsDashboardYtdStart => Year-To-Date Starts.
	// Path: reports/dashboard/ytd_start
	ReportsDashboardYtdStart model.Str

	// ReportsDashboardMtdStart => Current Month Starts.
	// Select day of the month.
	// Path: reports/dashboard/mtd_start
	ReportsDashboardMtdStart model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogRecentlyProductsScope = model.NewStr(`catalog/recently_products/scope`, model.WithPkgCfg(pkgCfg))
	pp.CatalogRecentlyProductsViewedCount = model.NewStr(`catalog/recently_products/viewed_count`, model.WithPkgCfg(pkgCfg))
	pp.CatalogRecentlyProductsComparedCount = model.NewStr(`catalog/recently_products/compared_count`, model.WithPkgCfg(pkgCfg))
	pp.ReportsDashboardYtdStart = model.NewStr(`reports/dashboard/ytd_start`, model.WithPkgCfg(pkgCfg))
	pp.ReportsDashboardMtdStart = model.NewStr(`reports/dashboard/mtd_start`, model.WithPkgCfg(pkgCfg))

	return pp
}
