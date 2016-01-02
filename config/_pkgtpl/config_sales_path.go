// +build ignore

package sales

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSalesGeneralHideCustomerIp => Hide Customer IP.
// Choose whether a customer IP is shown in orders, invoices, shipments, and
// credit memos.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesGeneralHideCustomerIp = model.NewBool(`sales/general/hide_customer_ip`, model.WithPkgCfg(PackageConfiguration))

// PathSalesTotalsSortDiscount => Discount.
var PathSalesTotalsSortDiscount = model.NewStr(`sales/totals_sort/discount`, model.WithPkgCfg(PackageConfiguration))

// PathSalesTotalsSortGrandTotal => Grand Total.
var PathSalesTotalsSortGrandTotal = model.NewStr(`sales/totals_sort/grand_total`, model.WithPkgCfg(PackageConfiguration))

// PathSalesTotalsSortShipping => Shipping.
var PathSalesTotalsSortShipping = model.NewStr(`sales/totals_sort/shipping`, model.WithPkgCfg(PackageConfiguration))

// PathSalesTotalsSortSubtotal => Subtotal.
var PathSalesTotalsSortSubtotal = model.NewStr(`sales/totals_sort/subtotal`, model.WithPkgCfg(PackageConfiguration))

// PathSalesTotalsSortTax => Tax.
var PathSalesTotalsSortTax = model.NewStr(`sales/totals_sort/tax`, model.WithPkgCfg(PackageConfiguration))

// PathSalesReorderAllow => Allow Reorder.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesReorderAllow = model.NewBool(`sales/reorder/allow`, model.WithPkgCfg(PackageConfiguration))

// PathSalesIdentityLogo => Logo for PDF Print-outs (200x50).
// Your default logo will be used in PDF and HTML documents.(jpeg, tiff, png)
// If your pdf image is distorted, try to use larger file-size image.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Pdf
var PathSalesIdentityLogo = model.NewStr(`sales/identity/logo`, model.WithPkgCfg(PackageConfiguration))

// PathSalesIdentityLogoHtml => Logo for HTML Print View.
// Logo for HTML documents only. If empty, default will be used.(jpeg, gif,
// png)
// BackendModel: Otnegam\Config\Model\Config\Backend\Image
var PathSalesIdentityLogoHtml = model.NewStr(`sales/identity/logo_html`, model.WithPkgCfg(PackageConfiguration))

// PathSalesIdentityAddress => Address.
var PathSalesIdentityAddress = model.NewStr(`sales/identity/address`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderActive => Enable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderActive = model.NewBool(`sales/minimum_order/active`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderAmount => Minimum Amount.
// Subtotal after discount
var PathSalesMinimumOrderAmount = model.NewStr(`sales/minimum_order/amount`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderTaxIncluding => Include Tax to Amount.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderTaxIncluding = model.NewBool(`sales/minimum_order/tax_including`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderDescription => Description Message.
// This message will be shown in the shopping cart when the subtotal (after
// discount) is lower than the minimum allowed amount.
var PathSalesMinimumOrderDescription = model.NewStr(`sales/minimum_order/description`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderErrorMessage => Error to Show in Shopping Cart.
var PathSalesMinimumOrderErrorMessage = model.NewStr(`sales/minimum_order/error_message`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderMultiAddress => Validate Each Address Separately in Multi-address Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesMinimumOrderMultiAddress = model.NewBool(`sales/minimum_order/multi_address`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderMultiAddressDescription => Multi-address Description Message.
// We'll use the default description above if you leave this empty.
var PathSalesMinimumOrderMultiAddressDescription = model.NewStr(`sales/minimum_order/multi_address_description`, model.WithPkgCfg(PackageConfiguration))

// PathSalesMinimumOrderMultiAddressErrorMessage => Multi-address Error to Show in Shopping Cart.
// We'll use the default error above if you leave this empty.
var PathSalesMinimumOrderMultiAddressErrorMessage = model.NewStr(`sales/minimum_order/multi_address_error_message`, model.WithPkgCfg(PackageConfiguration))

// PathSalesDashboardUseAggregatedData => Use Aggregated Data (beta).
// Improves dashboard performance but provides non-realtime data.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesDashboardUseAggregatedData = model.NewBool(`sales/dashboard/use_aggregated_data`, model.WithPkgCfg(PackageConfiguration))

// PathSalesOrdersDeletePendingAfter => Pending Payment Order Lifetime (minutes).
var PathSalesOrdersDeletePendingAfter = model.NewStr(`sales/orders/delete_pending_after`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailGeneralAsyncSending => Asynchronous sending.
// BackendModel: Otnegam\Sales\Model\Config\Backend\Email\AsyncSending
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathSalesEmailGeneralAsyncSending = model.NewBool(`sales_email/general/async_sending`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailOrderEnabled = model.NewBool(`sales_email/order/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderIdentity => New Order Confirmation Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailOrderIdentity = model.NewStr(`sales_email/order/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderTemplate => New Order Confirmation Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderTemplate = model.NewStr(`sales_email/order/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderGuestTemplate => New Order Confirmation Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderGuestTemplate = model.NewStr(`sales_email/order/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCopyTo => Send Order Email Copy To.
// Comma-separated
var PathSalesEmailOrderCopyTo = model.NewStr(`sales_email/order/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCopyMethod => Send Order Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailOrderCopyMethod = model.NewStr(`sales_email/order/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailOrderCommentEnabled = model.NewBool(`sales_email/order_comment/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentIdentity => Order Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailOrderCommentIdentity = model.NewStr(`sales_email/order_comment/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentTemplate => Order Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderCommentTemplate = model.NewStr(`sales_email/order_comment/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentGuestTemplate => Order Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailOrderCommentGuestTemplate = model.NewStr(`sales_email/order_comment/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentCopyTo => Send Order Comment Email Copy To.
// Comma-separated
var PathSalesEmailOrderCommentCopyTo = model.NewStr(`sales_email/order_comment/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailOrderCommentCopyMethod => Send Order Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailOrderCommentCopyMethod = model.NewStr(`sales_email/order_comment/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailInvoiceEnabled = model.NewBool(`sales_email/invoice/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceIdentity => Invoice Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailInvoiceIdentity = model.NewStr(`sales_email/invoice/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceTemplate => Invoice Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceTemplate = model.NewStr(`sales_email/invoice/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceGuestTemplate => Invoice Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceGuestTemplate = model.NewStr(`sales_email/invoice/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCopyTo => Send Invoice Email Copy To.
// Comma-separated
var PathSalesEmailInvoiceCopyTo = model.NewStr(`sales_email/invoice/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCopyMethod => Send Invoice Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailInvoiceCopyMethod = model.NewStr(`sales_email/invoice/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailInvoiceCommentEnabled = model.NewBool(`sales_email/invoice_comment/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentIdentity => Invoice Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailInvoiceCommentIdentity = model.NewStr(`sales_email/invoice_comment/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentTemplate => Invoice Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceCommentTemplate = model.NewStr(`sales_email/invoice_comment/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentGuestTemplate => Invoice Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailInvoiceCommentGuestTemplate = model.NewStr(`sales_email/invoice_comment/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentCopyTo => Send Invoice Comment Email Copy To.
// Comma-separated
var PathSalesEmailInvoiceCommentCopyTo = model.NewStr(`sales_email/invoice_comment/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailInvoiceCommentCopyMethod => Send Invoice Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailInvoiceCommentCopyMethod = model.NewStr(`sales_email/invoice_comment/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailShipmentEnabled = model.NewBool(`sales_email/shipment/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentIdentity => Shipment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailShipmentIdentity = model.NewStr(`sales_email/shipment/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentTemplate => Shipment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentTemplate = model.NewStr(`sales_email/shipment/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentGuestTemplate => Shipment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentGuestTemplate = model.NewStr(`sales_email/shipment/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCopyTo => Send Shipment Email Copy To.
// Comma-separated
var PathSalesEmailShipmentCopyTo = model.NewStr(`sales_email/shipment/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCopyMethod => Send Shipment Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailShipmentCopyMethod = model.NewStr(`sales_email/shipment/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailShipmentCommentEnabled = model.NewBool(`sales_email/shipment_comment/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentIdentity => Shipment Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailShipmentCommentIdentity = model.NewStr(`sales_email/shipment_comment/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentTemplate => Shipment Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentCommentTemplate = model.NewStr(`sales_email/shipment_comment/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentGuestTemplate => Shipment Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailShipmentCommentGuestTemplate = model.NewStr(`sales_email/shipment_comment/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentCopyTo => Send Shipment Comment Email Copy To.
// Comma-separated
var PathSalesEmailShipmentCommentCopyTo = model.NewStr(`sales_email/shipment_comment/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailShipmentCommentCopyMethod => Send Shipment Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailShipmentCommentCopyMethod = model.NewStr(`sales_email/shipment_comment/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailCreditmemoEnabled = model.NewBool(`sales_email/creditmemo/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoIdentity => Credit Memo Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailCreditmemoIdentity = model.NewStr(`sales_email/creditmemo/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoTemplate => Credit Memo Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoTemplate = model.NewStr(`sales_email/creditmemo/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoGuestTemplate => Credit Memo Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoGuestTemplate = model.NewStr(`sales_email/creditmemo/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCopyTo => Send Credit Memo Email Copy To.
// Comma-separated
var PathSalesEmailCreditmemoCopyTo = model.NewStr(`sales_email/creditmemo/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCopyMethod => Send Credit Memo Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailCreditmemoCopyMethod = model.NewStr(`sales_email/creditmemo/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesEmailCreditmemoCommentEnabled = model.NewBool(`sales_email/creditmemo_comment/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentIdentity => Credit Memo Comment Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathSalesEmailCreditmemoCommentIdentity = model.NewStr(`sales_email/creditmemo_comment/identity`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentTemplate => Credit Memo Comment Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoCommentTemplate = model.NewStr(`sales_email/creditmemo_comment/template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentGuestTemplate => Credit Memo Comment Email Template for Guest.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathSalesEmailCreditmemoCommentGuestTemplate = model.NewStr(`sales_email/creditmemo_comment/guest_template`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentCopyTo => Send Credit Memo Comment Email Copy To.
// Comma-separated
var PathSalesEmailCreditmemoCommentCopyTo = model.NewStr(`sales_email/creditmemo_comment/copy_to`, model.WithPkgCfg(PackageConfiguration))

// PathSalesEmailCreditmemoCommentCopyMethod => Send Credit Memo Comments Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathSalesEmailCreditmemoCommentCopyMethod = model.NewStr(`sales_email/creditmemo_comment/copy_method`, model.WithPkgCfg(PackageConfiguration))

// PathSalesPdfInvoicePutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfInvoicePutOrderId = model.NewBool(`sales_pdf/invoice/put_order_id`, model.WithPkgCfg(PackageConfiguration))

// PathSalesPdfShipmentPutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfShipmentPutOrderId = model.NewBool(`sales_pdf/shipment/put_order_id`, model.WithPkgCfg(PackageConfiguration))

// PathSalesPdfCreditmemoPutOrderId => Display Order ID in Header.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSalesPdfCreditmemoPutOrderId = model.NewBool(`sales_pdf/creditmemo/put_order_id`, model.WithPkgCfg(PackageConfiguration))

// PathRssOrderStatus => Customer Order Status Notification.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssOrderStatus = model.NewBool(`rss/order/status`, model.WithPkgCfg(PackageConfiguration))

// PathDevGridAsyncIndexing => Asynchronous indexing.
// BackendModel: Otnegam\Sales\Model\Config\Backend\Grid\AsyncIndexing
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathDevGridAsyncIndexing = model.NewBool(`dev/grid/async_indexing`, model.WithPkgCfg(PackageConfiguration))
