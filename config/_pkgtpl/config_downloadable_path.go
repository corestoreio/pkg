// +build ignore

package downloadable

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogDownloadableOrderItemStatus => Order Item Status to Enable Downloads.
// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Orderitemstatus
var PathCatalogDownloadableOrderItemStatus = model.NewStr(`catalog/downloadable/order_item_status`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableDownloadsNumber => Default Maximum Number of Downloads.
var PathCatalogDownloadableDownloadsNumber = model.NewStr(`catalog/downloadable/downloads_number`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableShareable => Shareable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogDownloadableShareable = model.NewBool(`catalog/downloadable/shareable`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableSamplesTitle => Default Sample Title.
var PathCatalogDownloadableSamplesTitle = model.NewStr(`catalog/downloadable/samples_title`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableLinksTitle => Default Link Title.
var PathCatalogDownloadableLinksTitle = model.NewStr(`catalog/downloadable/links_title`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableLinksTargetNewWindow => Open Links in New Window.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogDownloadableLinksTargetNewWindow = model.NewBool(`catalog/downloadable/links_target_new_window`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableContentDisposition => Use Content-Disposition.
// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Contentdisposition
var PathCatalogDownloadableContentDisposition = model.NewStr(`catalog/downloadable/content_disposition`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogDownloadableDisableGuestCheckout => Disable Guest Checkout if Cart Contains Downloadable Items.
// Guest checkout will only work with shareable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogDownloadableDisableGuestCheckout = model.NewBool(`catalog/downloadable/disable_guest_checkout`, model.WithPkgCfg(PackageConfiguration))
