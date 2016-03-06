// +build ignore

package sales

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
	// SalesGeneralHideCustomerIp => Hide Customer IP.
	// Choose whether a customer IP is shown in orders, invoices, shipments, and
	// credit memos.
	// Path: sales/general/hide_customer_ip
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesGeneralHideCustomerIp model.Bool

	// SalesTotalsSortDiscount => Discount.
	// Path: sales/totals_sort/discount
	SalesTotalsSortDiscount model.Str

	// SalesTotalsSortGrandTotal => Grand Total.
	// Path: sales/totals_sort/grand_total
	SalesTotalsSortGrandTotal model.Str

	// SalesTotalsSortShipping => Shipping.
	// Path: sales/totals_sort/shipping
	SalesTotalsSortShipping model.Str

	// SalesTotalsSortSubtotal => Subtotal.
	// Path: sales/totals_sort/subtotal
	SalesTotalsSortSubtotal model.Str

	// SalesTotalsSortTax => Tax.
	// Path: sales/totals_sort/tax
	SalesTotalsSortTax model.Str

	// SalesReorderAllow => Allow Reorder.
	// Path: sales/reorder/allow
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesReorderAllow model.Bool

	// SalesIdentityLogo => Logo for PDF Print-outs (200x50).
	// Your default logo will be used in PDF and HTML documents.(jpeg, tiff, png)
	// If your pdf image is distorted, try to use larger file-size image.
	// Path: sales/identity/logo
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Pdf
	SalesIdentityLogo model.Str

	// SalesIdentityLogoHtml => Logo for HTML Print View.
	// Logo for HTML documents only. If empty, default will be used.(jpeg, gif,
	// png)
	// Path: sales/identity/logo_html
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	SalesIdentityLogoHtml model.Str

	// SalesIdentityAddress => Address.
	// Path: sales/identity/address
	SalesIdentityAddress model.Str

	// SalesMinimumOrderActive => Enable.
	// Path: sales/minimum_order/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderActive model.Bool

	// SalesMinimumOrderAmount => Minimum Amount.
	// Subtotal after discount
	// Path: sales/minimum_order/amount
	SalesMinimumOrderAmount model.Str

	// SalesMinimumOrderTaxIncluding => Include Tax to Amount.
	// Path: sales/minimum_order/tax_including
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderTaxIncluding model.Bool

	// SalesMinimumOrderDescription => Description Message.
	// This message will be shown in the shopping cart when the subtotal (after
	// discount) is lower than the minimum allowed amount.
	// Path: sales/minimum_order/description
	SalesMinimumOrderDescription model.Str

	// SalesMinimumOrderErrorMessage => Error to Show in Shopping Cart.
	// Path: sales/minimum_order/error_message
	SalesMinimumOrderErrorMessage model.Str

	// SalesMinimumOrderMultiAddress => Validate Each Address Separately in Multi-address Checkout.
	// Path: sales/minimum_order/multi_address
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderMultiAddress model.Bool

	// SalesMinimumOrderMultiAddressDescription => Multi-address Description Message.
	// We'll use the default description above if you leave this empty.
	// Path: sales/minimum_order/multi_address_description
	SalesMinimumOrderMultiAddressDescription model.Str

	// SalesMinimumOrderMultiAddressErrorMessage => Multi-address Error to Show in Shopping Cart.
	// We'll use the default error above if you leave this empty.
	// Path: sales/minimum_order/multi_address_error_message
	SalesMinimumOrderMultiAddressErrorMessage model.Str

	// SalesDashboardUseAggregatedData => Use Aggregated Data (beta).
	// Improves dashboard performance but provides non-realtime data.
	// Path: sales/dashboard/use_aggregated_data
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesDashboardUseAggregatedData model.Bool

	// SalesOrdersDeletePendingAfter => Pending Payment Order Lifetime (minutes).
	// Path: sales/orders/delete_pending_after
	SalesOrdersDeletePendingAfter model.Str

	// SalesEmailGeneralAsyncSending => Asynchronous sending.
	// Path: sales_email/general/async_sending
	// BackendModel: Magento\Sales\Model\Config\Backend\Email\AsyncSending
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	SalesEmailGeneralAsyncSending model.Bool

	// SalesEmailOrderEnabled => Enabled.
	// Path: sales_email/order/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailOrderEnabled model.Bool

	// SalesEmailOrderIdentity => New Order Confirmation Email Sender.
	// Path: sales_email/order/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderIdentity model.Str

	// SalesEmailOrderTemplate => New Order Confirmation Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderTemplate model.Str

	// SalesEmailOrderGuestTemplate => New Order Confirmation Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderGuestTemplate model.Str

	// SalesEmailOrderCopyTo => Send Order Email Copy To.
	// Comma-separated
	// Path: sales_email/order/copy_to
	SalesEmailOrderCopyTo model.Str

	// SalesEmailOrderCopyMethod => Send Order Email Copy Method.
	// Path: sales_email/order/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCopyMethod model.Str

	// SalesEmailOrderCommentEnabled => Enabled.
	// Path: sales_email/order_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailOrderCommentEnabled model.Bool

	// SalesEmailOrderCommentIdentity => Order Comment Email Sender.
	// Path: sales_email/order_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderCommentIdentity model.Str

	// SalesEmailOrderCommentTemplate => Order Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentTemplate model.Str

	// SalesEmailOrderCommentGuestTemplate => Order Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentGuestTemplate model.Str

	// SalesEmailOrderCommentCopyTo => Send Order Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/order_comment/copy_to
	SalesEmailOrderCommentCopyTo model.Str

	// SalesEmailOrderCommentCopyMethod => Send Order Comments Email Copy Method.
	// Path: sales_email/order_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCommentCopyMethod model.Str

	// SalesEmailInvoiceEnabled => Enabled.
	// Path: sales_email/invoice/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceEnabled model.Bool

	// SalesEmailInvoiceIdentity => Invoice Email Sender.
	// Path: sales_email/invoice/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceIdentity model.Str

	// SalesEmailInvoiceTemplate => Invoice Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceTemplate model.Str

	// SalesEmailInvoiceGuestTemplate => Invoice Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceGuestTemplate model.Str

	// SalesEmailInvoiceCopyTo => Send Invoice Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice/copy_to
	SalesEmailInvoiceCopyTo model.Str

	// SalesEmailInvoiceCopyMethod => Send Invoice Email Copy Method.
	// Path: sales_email/invoice/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCopyMethod model.Str

	// SalesEmailInvoiceCommentEnabled => Enabled.
	// Path: sales_email/invoice_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceCommentEnabled model.Bool

	// SalesEmailInvoiceCommentIdentity => Invoice Comment Email Sender.
	// Path: sales_email/invoice_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceCommentIdentity model.Str

	// SalesEmailInvoiceCommentTemplate => Invoice Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentTemplate model.Str

	// SalesEmailInvoiceCommentGuestTemplate => Invoice Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentGuestTemplate model.Str

	// SalesEmailInvoiceCommentCopyTo => Send Invoice Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice_comment/copy_to
	SalesEmailInvoiceCommentCopyTo model.Str

	// SalesEmailInvoiceCommentCopyMethod => Send Invoice Comments Email Copy Method.
	// Path: sales_email/invoice_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCommentCopyMethod model.Str

	// SalesEmailShipmentEnabled => Enabled.
	// Path: sales_email/shipment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailShipmentEnabled model.Bool

	// SalesEmailShipmentIdentity => Shipment Email Sender.
	// Path: sales_email/shipment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentIdentity model.Str

	// SalesEmailShipmentTemplate => Shipment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentTemplate model.Str

	// SalesEmailShipmentGuestTemplate => Shipment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentGuestTemplate model.Str

	// SalesEmailShipmentCopyTo => Send Shipment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment/copy_to
	SalesEmailShipmentCopyTo model.Str

	// SalesEmailShipmentCopyMethod => Send Shipment Email Copy Method.
	// Path: sales_email/shipment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCopyMethod model.Str

	// SalesEmailShipmentCommentEnabled => Enabled.
	// Path: sales_email/shipment_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailShipmentCommentEnabled model.Bool

	// SalesEmailShipmentCommentIdentity => Shipment Comment Email Sender.
	// Path: sales_email/shipment_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentCommentIdentity model.Str

	// SalesEmailShipmentCommentTemplate => Shipment Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentTemplate model.Str

	// SalesEmailShipmentCommentGuestTemplate => Shipment Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentGuestTemplate model.Str

	// SalesEmailShipmentCommentCopyTo => Send Shipment Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment_comment/copy_to
	SalesEmailShipmentCommentCopyTo model.Str

	// SalesEmailShipmentCommentCopyMethod => Send Shipment Comments Email Copy Method.
	// Path: sales_email/shipment_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCommentCopyMethod model.Str

	// SalesEmailCreditmemoEnabled => Enabled.
	// Path: sales_email/creditmemo/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoEnabled model.Bool

	// SalesEmailCreditmemoIdentity => Credit Memo Email Sender.
	// Path: sales_email/creditmemo/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoIdentity model.Str

	// SalesEmailCreditmemoTemplate => Credit Memo Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoTemplate model.Str

	// SalesEmailCreditmemoGuestTemplate => Credit Memo Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoGuestTemplate model.Str

	// SalesEmailCreditmemoCopyTo => Send Credit Memo Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo/copy_to
	SalesEmailCreditmemoCopyTo model.Str

	// SalesEmailCreditmemoCopyMethod => Send Credit Memo Email Copy Method.
	// Path: sales_email/creditmemo/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCopyMethod model.Str

	// SalesEmailCreditmemoCommentEnabled => Enabled.
	// Path: sales_email/creditmemo_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoCommentEnabled model.Bool

	// SalesEmailCreditmemoCommentIdentity => Credit Memo Comment Email Sender.
	// Path: sales_email/creditmemo_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoCommentIdentity model.Str

	// SalesEmailCreditmemoCommentTemplate => Credit Memo Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentTemplate model.Str

	// SalesEmailCreditmemoCommentGuestTemplate => Credit Memo Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentGuestTemplate model.Str

	// SalesEmailCreditmemoCommentCopyTo => Send Credit Memo Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo_comment/copy_to
	SalesEmailCreditmemoCommentCopyTo model.Str

	// SalesEmailCreditmemoCommentCopyMethod => Send Credit Memo Comments Email Copy Method.
	// Path: sales_email/creditmemo_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCommentCopyMethod model.Str

	// SalesPdfInvoicePutOrderId => Display Order ID in Header.
	// Path: sales_pdf/invoice/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfInvoicePutOrderId model.Bool

	// SalesPdfShipmentPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/shipment/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfShipmentPutOrderId model.Bool

	// SalesPdfCreditmemoPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/creditmemo/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfCreditmemoPutOrderId model.Bool

	// RssOrderStatus => Customer Order Status Notification.
	// Path: rss/order/status
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssOrderStatus model.Bool

	// DevGridAsyncIndexing => Asynchronous indexing.
	// Path: dev/grid/async_indexing
	// BackendModel: Magento\Sales\Model\Config\Backend\Grid\AsyncIndexing
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	DevGridAsyncIndexing model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesGeneralHideCustomerIp = model.NewBool(`sales/general/hide_customer_ip`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortDiscount = model.NewStr(`sales/totals_sort/discount`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortGrandTotal = model.NewStr(`sales/totals_sort/grand_total`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortShipping = model.NewStr(`sales/totals_sort/shipping`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortSubtotal = model.NewStr(`sales/totals_sort/subtotal`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortTax = model.NewStr(`sales/totals_sort/tax`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesReorderAllow = model.NewBool(`sales/reorder/allow`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityLogo = model.NewStr(`sales/identity/logo`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityLogoHtml = model.NewStr(`sales/identity/logo_html`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityAddress = model.NewStr(`sales/identity/address`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderActive = model.NewBool(`sales/minimum_order/active`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderAmount = model.NewStr(`sales/minimum_order/amount`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderTaxIncluding = model.NewBool(`sales/minimum_order/tax_including`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderDescription = model.NewStr(`sales/minimum_order/description`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderErrorMessage = model.NewStr(`sales/minimum_order/error_message`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddress = model.NewBool(`sales/minimum_order/multi_address`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddressDescription = model.NewStr(`sales/minimum_order/multi_address_description`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddressErrorMessage = model.NewStr(`sales/minimum_order/multi_address_error_message`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesDashboardUseAggregatedData = model.NewBool(`sales/dashboard/use_aggregated_data`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesOrdersDeletePendingAfter = model.NewStr(`sales/orders/delete_pending_after`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailGeneralAsyncSending = model.NewBool(`sales_email/general/async_sending`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderEnabled = model.NewBool(`sales_email/order/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderIdentity = model.NewStr(`sales_email/order/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderTemplate = model.NewStr(`sales_email/order/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderGuestTemplate = model.NewStr(`sales_email/order/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCopyTo = model.NewStr(`sales_email/order/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCopyMethod = model.NewStr(`sales_email/order/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentEnabled = model.NewBool(`sales_email/order_comment/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentIdentity = model.NewStr(`sales_email/order_comment/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentTemplate = model.NewStr(`sales_email/order_comment/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentGuestTemplate = model.NewStr(`sales_email/order_comment/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentCopyTo = model.NewStr(`sales_email/order_comment/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentCopyMethod = model.NewStr(`sales_email/order_comment/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceEnabled = model.NewBool(`sales_email/invoice/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceIdentity = model.NewStr(`sales_email/invoice/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceTemplate = model.NewStr(`sales_email/invoice/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceGuestTemplate = model.NewStr(`sales_email/invoice/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCopyTo = model.NewStr(`sales_email/invoice/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCopyMethod = model.NewStr(`sales_email/invoice/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentEnabled = model.NewBool(`sales_email/invoice_comment/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentIdentity = model.NewStr(`sales_email/invoice_comment/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentTemplate = model.NewStr(`sales_email/invoice_comment/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentGuestTemplate = model.NewStr(`sales_email/invoice_comment/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentCopyTo = model.NewStr(`sales_email/invoice_comment/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentCopyMethod = model.NewStr(`sales_email/invoice_comment/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentEnabled = model.NewBool(`sales_email/shipment/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentIdentity = model.NewStr(`sales_email/shipment/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentTemplate = model.NewStr(`sales_email/shipment/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentGuestTemplate = model.NewStr(`sales_email/shipment/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCopyTo = model.NewStr(`sales_email/shipment/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCopyMethod = model.NewStr(`sales_email/shipment/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentEnabled = model.NewBool(`sales_email/shipment_comment/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentIdentity = model.NewStr(`sales_email/shipment_comment/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentTemplate = model.NewStr(`sales_email/shipment_comment/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentGuestTemplate = model.NewStr(`sales_email/shipment_comment/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentCopyTo = model.NewStr(`sales_email/shipment_comment/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentCopyMethod = model.NewStr(`sales_email/shipment_comment/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoEnabled = model.NewBool(`sales_email/creditmemo/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoIdentity = model.NewStr(`sales_email/creditmemo/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoTemplate = model.NewStr(`sales_email/creditmemo/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoGuestTemplate = model.NewStr(`sales_email/creditmemo/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCopyTo = model.NewStr(`sales_email/creditmemo/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCopyMethod = model.NewStr(`sales_email/creditmemo/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentEnabled = model.NewBool(`sales_email/creditmemo_comment/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentIdentity = model.NewStr(`sales_email/creditmemo_comment/identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentTemplate = model.NewStr(`sales_email/creditmemo_comment/template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentGuestTemplate = model.NewStr(`sales_email/creditmemo_comment/guest_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentCopyTo = model.NewStr(`sales_email/creditmemo_comment/copy_to`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentCopyMethod = model.NewStr(`sales_email/creditmemo_comment/copy_method`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfInvoicePutOrderId = model.NewBool(`sales_pdf/invoice/put_order_id`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfShipmentPutOrderId = model.NewBool(`sales_pdf/shipment/put_order_id`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfCreditmemoPutOrderId = model.NewBool(`sales_pdf/creditmemo/put_order_id`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.RssOrderStatus = model.NewBool(`rss/order/status`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DevGridAsyncIndexing = model.NewBool(`dev/grid/async_indexing`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
