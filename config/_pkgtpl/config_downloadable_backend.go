// +build ignore

package downloadable

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
	// CatalogDownloadableOrderItemStatus => Order Item Status to Enable Downloads.
	// Path: catalog/downloadable/order_item_status
	// SourceModel: Magento\Downloadable\Model\System\Config\Source\Orderitemstatus
	CatalogDownloadableOrderItemStatus model.Str

	// CatalogDownloadableDownloadsNumber => Default Maximum Number of Downloads.
	// Path: catalog/downloadable/downloads_number
	CatalogDownloadableDownloadsNumber model.Str

	// CatalogDownloadableShareable => Shareable.
	// Path: catalog/downloadable/shareable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableShareable model.Bool

	// CatalogDownloadableSamplesTitle => Default Sample Title.
	// Path: catalog/downloadable/samples_title
	CatalogDownloadableSamplesTitle model.Str

	// CatalogDownloadableLinksTitle => Default Link Title.
	// Path: catalog/downloadable/links_title
	CatalogDownloadableLinksTitle model.Str

	// CatalogDownloadableLinksTargetNewWindow => Open Links in New Window.
	// Path: catalog/downloadable/links_target_new_window
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableLinksTargetNewWindow model.Bool

	// CatalogDownloadableContentDisposition => Use Content-Disposition.
	// Path: catalog/downloadable/content_disposition
	// SourceModel: Magento\Downloadable\Model\System\Config\Source\Contentdisposition
	CatalogDownloadableContentDisposition model.Str

	// CatalogDownloadableDisableGuestCheckout => Disable Guest Checkout if Cart Contains Downloadable Items.
	// Guest checkout will only work with shareable.
	// Path: catalog/downloadable/disable_guest_checkout
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogDownloadableDisableGuestCheckout model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogDownloadableOrderItemStatus = model.NewStr(`catalog/downloadable/order_item_status`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableDownloadsNumber = model.NewStr(`catalog/downloadable/downloads_number`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableShareable = model.NewBool(`catalog/downloadable/shareable`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableSamplesTitle = model.NewStr(`catalog/downloadable/samples_title`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableLinksTitle = model.NewStr(`catalog/downloadable/links_title`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableLinksTargetNewWindow = model.NewBool(`catalog/downloadable/links_target_new_window`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableContentDisposition = model.NewStr(`catalog/downloadable/content_disposition`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogDownloadableDisableGuestCheckout = model.NewBool(`catalog/downloadable/disable_guest_checkout`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
