// +build ignore

package reports

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogRecentlyProductsScope => Show for Current.
// SourceModel: Otnegam\Config\Model\Config\Source\Reports\Scope
var PathCatalogRecentlyProductsScope = model.NewStr(`catalog/recently_products/scope`)

// PathCatalogRecentlyProductsViewedCount => Default Recently Viewed Products Count.
var PathCatalogRecentlyProductsViewedCount = model.NewStr(`catalog/recently_products/viewed_count`)

// PathCatalogRecentlyProductsComparedCount => Default Recently Compared Products Count.
var PathCatalogRecentlyProductsComparedCount = model.NewStr(`catalog/recently_products/compared_count`)

// PathReportsDashboardYtdStart => Year-To-Date Starts.
var PathReportsDashboardYtdStart = model.NewStr(`reports/dashboard/ytd_start`)

// PathReportsDashboardMtdStart => Current Month Starts.
// Select day of the month.
var PathReportsDashboardMtdStart = model.NewStr(`reports/dashboard/mtd_start`)
