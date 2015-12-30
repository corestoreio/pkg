// +build ignore

package sales

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSalesGeneralHideCustomerIp => Hide Customer IP.
// Choose whether a customer IP is shown in orders, invoices, shipments, and
// credit memos.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesGeneralHideCustomerIp = model.NewBool(`sales/general/hide_customer_ip`)

// PathSalesTotalsSortDiscount => Discount.
var PathSalesTotalsSortDiscount = model.NewStr(`sales/totals_sort/discount`)

// PathSalesTotalsSortGrandTotal => Grand Total.
var PathSalesTotalsSortGrandTotal = model.NewStr(`sales/totals_sort/grand_total`)

// PathSalesTotalsSortShipping => Shipping.
var PathSalesTotalsSortShipping = model.NewStr(`sales/totals_sort/shipping`)

// PathSalesTotalsSortSubtotal => Subtotal.
var PathSalesTotalsSortSubtotal = model.NewStr(`sales/totals_sort/subtotal`)

// PathSalesTotalsSortTax => Tax.
var PathSalesTotalsSortTax = model.NewStr(`sales/totals_sort/tax`)

// PathSalesReorderAllow => Allow Reorder.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesReorderAllow = model.NewBool(`sales/reorder/allow`)

// PathSalesIdentityLogo => Logo for PDF Print-outs (200x50).
// Your default logo will be used in PDF and HTML documents.(jpeg, tiff, png)
// If your pdf image is distorted, try to use larger file-size image.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Pdf
var PathSalesIdentityLogo = model.NewStr(`sales/identity/logo`)

// PathSalesIdentityLogoHtml => Logo for HTML Print View.
// Logo for HTML documents only. If empty, default will be used.(jpeg, gif,
// png)
// BackendModel: Otnegam\Config\Model\Config\Backend\Image
var PathSalesIdentityLogoHtml = model.NewStr(`sales/identity/logo_html`)

// PathSalesIdentityAddress => Address.
var PathSalesIdentityAddress = model.NewStr(`sales/identity/address`)

// PathSalesMinimumOrderActive => Enable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderActive = model.NewBool(`sales/minimum_order/active`)

// PathSalesMinimumOrderAmount => Minimum Amount.
// Subtotal after discount
var PathSalesMinimumOrderAmount = model.NewStr(`sales/minimum_order/amount`)

// PathSalesMinimumOrderTaxIncluding => Include Tax to Amount.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderTaxIncluding = model.NewBool(`sales/minimum_order/tax_including`)

// PathSalesMinimumOrderDescription => Description Message.
// This message will be shown in the shopping cart when the subtotal (after
// discount) is lower than the minimum allowed amount.
var PathSalesMinimumOrderDescription = model.NewStr(`sales/minimum_order/description`)

// PathSalesMinimumOrderErrorMessage => Error to Show in Shopping Cart.
var PathSalesMinimumOrderErrorMessage = model.NewStr(`sales/minimum_order/error_message`)

// PathSalesMinimumOrderMultiAddress => Validate Each Address Separately in Multi-address Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderMultiAddress = model.NewBool(`sales/minimum_order/multi_address`)

// PathSalesMinimumOrderMultiAddressDescription => Multi-address Description Message.
// We'll use the default description above if you leave this empty.
var PathSalesMinimumOrderMultiAddressDescription = model.NewStr(`sales/minimum_order/multi_address_description`)

// PathSalesMinimumOrderMultiAddressErrorMessage => Multi-address Error to Show in Shopping Cart.
// We'll use the default error above if you leave this empty.
var PathSalesMinimumOrderMultiAddressErrorMessage = model.NewStr(`sales/minimum_order/multi_address_error_message`)

// PathSalesDashboardUseAggregatedData => Use Aggregated Data (beta).
// Improves dashboard performance but provides non-realtime data.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesDashboardUseAggregatedData = model.NewBool(`sales/dashboard/use_aggregated_data`)

// PathSalesOrdersDeletePendingAfter => Pending Payment Order Lifetime (minutes).
var PathSalesOrdersDeletePendingAfter = model.NewStr(`sales/orders/delete_pending_after`)

// PathSalesEmailGeneralAsyncSending => Asynchronous sending.
// BackendModel: Otnegam\Sales\Model\Config\Backend\Email\AsyncSending
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathSalesEmailGeneralAsyncSending = model.NewBool(`sales_email/general/async_sending`)

// PathSalesEmailOrderEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailOrderEnabled = model.NewBool(`sales_email/order/enabled`)

// PathSalesEmailOrderIdentity => New Order Confirmation Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailOrderIdentity = model.NewStr(`sales_email/order/identity`)

// PathSalesEmailOrderTemplate => New Order Confirmation Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderTemplate = model.NewStr(`sales_email/order/template`)

// PathSalesEmailOrderGuestTemplate => New Order Confirmation Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderGuestTemplate = model.NewStr(`sales_email/order/guest_template`)

// PathSalesEmailOrderCopyTo => Send Order Email Copy To.
// Comma-separated
var PathSalesEmailOrderCopyTo = model.NewStr(`sales_email/order/copy_to`)

// PathSalesEmailOrderCopyMethod => Send Order Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailOrderCopyMethod = model.NewStr(`sales_email/order/copy_method`)

// PathSalesEmailOrderCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailOrderCommentEnabled = model.NewBool(`sales_email/order_comment/enabled`)

// PathSalesEmailOrderCommentIdentity => Order Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailOrderCommentIdentity = model.NewStr(`sales_email/order_comment/identity`)

// PathSalesEmailOrderCommentTemplate => Order Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderCommentTemplate = model.NewStr(`sales_email/order_comment/template`)

// PathSalesEmailOrderCommentGuestTemplate => Order Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderCommentGuestTemplate = model.NewStr(`sales_email/order_comment/guest_template`)

// PathSalesEmailOrderCommentCopyTo => Send Order Comment Email Copy To.
// Comma-separated
var PathSalesEmailOrderCommentCopyTo = model.NewStr(`sales_email/order_comment/copy_to`)

// PathSalesEmailOrderCommentCopyMethod => Send Order Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailOrderCommentCopyMethod = model.NewStr(`sales_email/order_comment/copy_method`)

// PathSalesEmailInvoiceEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailInvoiceEnabled = model.NewBool(`sales_email/invoice/enabled`)

// PathSalesEmailInvoiceIdentity => Invoice Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailInvoiceIdentity = model.NewStr(`sales_email/invoice/identity`)

// PathSalesEmailInvoiceTemplate => Invoice Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceTemplate = model.NewStr(`sales_email/invoice/template`)

// PathSalesEmailInvoiceGuestTemplate => Invoice Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceGuestTemplate = model.NewStr(`sales_email/invoice/guest_template`)

// PathSalesEmailInvoiceCopyTo => Send Invoice Email Copy To.
// Comma-separated
var PathSalesEmailInvoiceCopyTo = model.NewStr(`sales_email/invoice/copy_to`)

// PathSalesEmailInvoiceCopyMethod => Send Invoice Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailInvoiceCopyMethod = model.NewStr(`sales_email/invoice/copy_method`)

// PathSalesEmailInvoiceCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailInvoiceCommentEnabled = model.NewBool(`sales_email/invoice_comment/enabled`)

// PathSalesEmailInvoiceCommentIdentity => Invoice Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailInvoiceCommentIdentity = model.NewStr(`sales_email/invoice_comment/identity`)

// PathSalesEmailInvoiceCommentTemplate => Invoice Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceCommentTemplate = model.NewStr(`sales_email/invoice_comment/template`)

// PathSalesEmailInvoiceCommentGuestTemplate => Invoice Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceCommentGuestTemplate = model.NewStr(`sales_email/invoice_comment/guest_template`)

// PathSalesEmailInvoiceCommentCopyTo => Send Invoice Comment Email Copy To.
// Comma-separated
var PathSalesEmailInvoiceCommentCopyTo = model.NewStr(`sales_email/invoice_comment/copy_to`)

// PathSalesEmailInvoiceCommentCopyMethod => Send Invoice Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailInvoiceCommentCopyMethod = model.NewStr(`sales_email/invoice_comment/copy_method`)

// PathSalesEmailShipmentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailShipmentEnabled = model.NewBool(`sales_email/shipment/enabled`)

// PathSalesEmailShipmentIdentity => Shipment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailShipmentIdentity = model.NewStr(`sales_email/shipment/identity`)

// PathSalesEmailShipmentTemplate => Shipment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentTemplate = model.NewStr(`sales_email/shipment/template`)

// PathSalesEmailShipmentGuestTemplate => Shipment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentGuestTemplate = model.NewStr(`sales_email/shipment/guest_template`)

// PathSalesEmailShipmentCopyTo => Send Shipment Email Copy To.
// Comma-separated
var PathSalesEmailShipmentCopyTo = model.NewStr(`sales_email/shipment/copy_to`)

// PathSalesEmailShipmentCopyMethod => Send Shipment Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailShipmentCopyMethod = model.NewStr(`sales_email/shipment/copy_method`)

// PathSalesEmailShipmentCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailShipmentCommentEnabled = model.NewBool(`sales_email/shipment_comment/enabled`)

// PathSalesEmailShipmentCommentIdentity => Shipment Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailShipmentCommentIdentity = model.NewStr(`sales_email/shipment_comment/identity`)

// PathSalesEmailShipmentCommentTemplate => Shipment Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentCommentTemplate = model.NewStr(`sales_email/shipment_comment/template`)

// PathSalesEmailShipmentCommentGuestTemplate => Shipment Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentCommentGuestTemplate = model.NewStr(`sales_email/shipment_comment/guest_template`)

// PathSalesEmailShipmentCommentCopyTo => Send Shipment Comment Email Copy To.
// Comma-separated
var PathSalesEmailShipmentCommentCopyTo = model.NewStr(`sales_email/shipment_comment/copy_to`)

// PathSalesEmailShipmentCommentCopyMethod => Send Shipment Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailShipmentCommentCopyMethod = model.NewStr(`sales_email/shipment_comment/copy_method`)

// PathSalesEmailCreditmemoEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailCreditmemoEnabled = model.NewBool(`sales_email/creditmemo/enabled`)

// PathSalesEmailCreditmemoIdentity => Credit Memo Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailCreditmemoIdentity = model.NewStr(`sales_email/creditmemo/identity`)

// PathSalesEmailCreditmemoTemplate => Credit Memo Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoTemplate = model.NewStr(`sales_email/creditmemo/template`)

// PathSalesEmailCreditmemoGuestTemplate => Credit Memo Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoGuestTemplate = model.NewStr(`sales_email/creditmemo/guest_template`)

// PathSalesEmailCreditmemoCopyTo => Send Credit Memo Email Copy To.
// Comma-separated
var PathSalesEmailCreditmemoCopyTo = model.NewStr(`sales_email/creditmemo/copy_to`)

// PathSalesEmailCreditmemoCopyMethod => Send Credit Memo Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailCreditmemoCopyMethod = model.NewStr(`sales_email/creditmemo/copy_method`)

// PathSalesEmailCreditmemoCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailCreditmemoCommentEnabled = model.NewBool(`sales_email/creditmemo_comment/enabled`)

// PathSalesEmailCreditmemoCommentIdentity => Credit Memo Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailCreditmemoCommentIdentity = model.NewStr(`sales_email/creditmemo_comment/identity`)

// PathSalesEmailCreditmemoCommentTemplate => Credit Memo Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoCommentTemplate = model.NewStr(`sales_email/creditmemo_comment/template`)

// PathSalesEmailCreditmemoCommentGuestTemplate => Credit Memo Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoCommentGuestTemplate = model.NewStr(`sales_email/creditmemo_comment/guest_template`)

// PathSalesEmailCreditmemoCommentCopyTo => Send Credit Memo Comment Email Copy To.
// Comma-separated
var PathSalesEmailCreditmemoCommentCopyTo = model.NewStr(`sales_email/creditmemo_comment/copy_to`)

// PathSalesEmailCreditmemoCommentCopyMethod => Send Credit Memo Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailCreditmemoCommentCopyMethod = model.NewStr(`sales_email/creditmemo_comment/copy_method`)

// PathSalesPdfInvoicePutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfInvoicePutOrderId = model.NewBool(`sales_pdf/invoice/put_order_id`)

// PathSalesPdfShipmentPutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfShipmentPutOrderId = model.NewBool(`sales_pdf/shipment/put_order_id`)

// PathSalesPdfCreditmemoPutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfCreditmemoPutOrderId = model.NewBool(`sales_pdf/creditmemo/put_order_id`)

// PathRssOrderStatus => Customer Order Status Notification.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssOrderStatus = model.NewBool(`rss/order/status`)

// PathDevGridAsyncIndexing => Asynchronous indexing.
// BackendModel: Otnegam\Sales\Model\Config\Backend\Grid\AsyncIndexing
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathDevGridAsyncIndexing = model.NewBool(`dev/grid/async_indexing`)
