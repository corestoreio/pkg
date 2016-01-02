// +build ignore

package downloadable

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
	// CatalogDownloadableOrderItemStatus => Order Item Status to Enable Downloads.
	// Path: catalog/downloadable/order_item_status
	// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Orderitemstatus
	CatalogDownloadableOrderItemStatus model.Str

	// CatalogDownloadableDownloadsNumber => Default Maximum Number of Downloads.
	// Path: catalog/downloadable/downloads_number
	CatalogDownloadableDownloadsNumber model.Str

	// CatalogDownloadableShareable => Shareable.
	// Path: catalog/downloadable/shareable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogDownloadableShareable model.Bool

	// CatalogDownloadableSamplesTitle => Default Sample Title.
	// Path: catalog/downloadable/samples_title
	CatalogDownloadableSamplesTitle model.Str

	// CatalogDownloadableLinksTitle => Default Link Title.
	// Path: catalog/downloadable/links_title
	CatalogDownloadableLinksTitle model.Str

	// CatalogDownloadableLinksTargetNewWindow => Open Links in New Window.
	// Path: catalog/downloadable/links_target_new_window
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogDownloadableLinksTargetNewWindow model.Bool

	// CatalogDownloadableContentDisposition => Use Content-Disposition.
	// Path: catalog/downloadable/content_disposition
	// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Contentdisposition
	CatalogDownloadableContentDisposition model.Str

	// CatalogDownloadableDisableGuestCheckout => Disable Guest Checkout if Cart Contains Downloadable Items.
	// Guest checkout will only work with shareable.
	// Path: catalog/downloadable/disable_guest_checkout
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogDownloadableDisableGuestCheckout model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogDownloadableOrderItemStatus = model.NewStr(`catalog/downloadable/order_item_status`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableDownloadsNumber = model.NewStr(`catalog/downloadable/downloads_number`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableShareable = model.NewBool(`catalog/downloadable/shareable`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableSamplesTitle = model.NewStr(`catalog/downloadable/samples_title`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableLinksTitle = model.NewStr(`catalog/downloadable/links_title`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableLinksTargetNewWindow = model.NewBool(`catalog/downloadable/links_target_new_window`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableContentDisposition = model.NewStr(`catalog/downloadable/content_disposition`, model.WithPkgCfg(pkgCfg))
	pp.CatalogDownloadableDisableGuestCheckout = model.NewBool(`catalog/downloadable/disable_guest_checkout`, model.WithPkgCfg(pkgCfg))

	return pp
}
