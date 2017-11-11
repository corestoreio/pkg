// +build ignore

package downloadable

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
	// CatalogDownloadableOrderItemStatus => Order Item Status to Enable Downloads.
	// Path: catalog/downloadable/order_item_status
	// SourceModel: Magento\Downloadable\Model\System\Config\Source\Orderitemstatus
	CatalogDownloadableOrderItemStatus cfgmodel.Str

	// CatalogDownloadableDownloadsNumber => Default Maximum Number of Downloads.
	// Path: catalog/downloadable/downloads_number
	CatalogDownloadableDownloadsNumber cfgmodel.Str

	// CatalogDownloadableShareable => Shareable.
	// Path: catalog/downloadable/shareable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableShareable cfgmodel.Bool

	// CatalogDownloadableSamplesTitle => Default Sample Title.
	// Path: catalog/downloadable/samples_title
	CatalogDownloadableSamplesTitle cfgmodel.Str

	// CatalogDownloadableLinksTitle => Default Link Title.
	// Path: catalog/downloadable/links_title
	CatalogDownloadableLinksTitle cfgmodel.Str

	// CatalogDownloadableLinksTargetNewWindow => Open Links in New Window.
	// Path: catalog/downloadable/links_target_new_window
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableLinksTargetNewWindow cfgmodel.Bool

	// CatalogDownloadableContentDisposition => Use Content-Disposition.
	// Path: catalog/downloadable/content_disposition
	// SourceModel: Magento\Downloadable\Model\System\Config\Source\Contentdisposition
	CatalogDownloadableContentDisposition cfgmodel.Str

	// CatalogDownloadableDisableGuestCheckout => Disable Guest Checkout if Cart Contains Downloadable Items.
	// Guest checkout will only work with shareable.
	// Path: catalog/downloadable/disable_guest_checkout
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableDisableGuestCheckout cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogDownloadableOrderItemStatus = cfgmodel.NewStr(`catalog/downloadable/order_item_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableDownloadsNumber = cfgmodel.NewStr(`catalog/downloadable/downloads_number`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableShareable = cfgmodel.NewBool(`catalog/downloadable/shareable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableSamplesTitle = cfgmodel.NewStr(`catalog/downloadable/samples_title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableLinksTitle = cfgmodel.NewStr(`catalog/downloadable/links_title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableLinksTargetNewWindow = cfgmodel.NewBool(`catalog/downloadable/links_target_new_window`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableContentDisposition = cfgmodel.NewStr(`catalog/downloadable/content_disposition`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableDisableGuestCheckout = cfgmodel.NewBool(`catalog/downloadable/disable_guest_checkout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
