// +build ignore

package sales

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
	// SalesGeneralHideCustomerIp => Hide Customer IP.
	// Choose whether a customer IP is shown in orders, invoices, shipments, and
	// credit memos.
	// Path: sales/general/hide_customer_ip
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
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
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesReorderAllow model.Bool

	// SalesIdentityLogo => Logo for PDF Print-outs (200x50).
	// Your default logo will be used in PDF and HTML documents.(jpeg, tiff, png)
	// If your pdf image is distorted, try to use larger file-size image.
	// Path: sales/identity/logo
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Pdf
	SalesIdentityLogo model.Str

	// SalesIdentityLogoHtml => Logo for HTML Print View.
	// Logo for HTML documents only. If empty, default will be used.(jpeg, gif,
	// png)
	// Path: sales/identity/logo_html
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image
	SalesIdentityLogoHtml model.Str

	// SalesIdentityAddress => Address.
	// Path: sales/identity/address
	SalesIdentityAddress model.Str

	// SalesMinimumOrderActive => Enable.
	// Path: sales/minimum_order/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesMinimumOrderActive model.Bool

	// SalesMinimumOrderAmount => Minimum Amount.
	// Subtotal after discount
	// Path: sales/minimum_order/amount
	SalesMinimumOrderAmount model.Str

	// SalesMinimumOrderTaxIncluding => Include Tax to Amount.
	// Path: sales/minimum_order/tax_including
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
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
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
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
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesDashboardUseAggregatedData model.Bool

	// SalesOrdersDeletePendingAfter => Pending Payment Order Lifetime (minutes).
	// Path: sales/orders/delete_pending_after
	SalesOrdersDeletePendingAfter model.Str

	// SalesEmailGeneralAsyncSending => Asynchronous sending.
	// Path: sales_email/general/async_sending
	// BackendModel: Otnegam\Sales\Model\Config\Backend\Email\AsyncSending
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	SalesEmailGeneralAsyncSending model.Bool

	// SalesEmailOrderEnabled => Enabled.
	// Path: sales_email/order/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailOrderEnabled model.Bool

	// SalesEmailOrderIdentity => New Order Confirmation Email Sender.
	// Path: sales_email/order/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderIdentity model.Str

	// SalesEmailOrderTemplate => New Order Confirmation Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailOrderTemplate model.Str

	// SalesEmailOrderGuestTemplate => New Order Confirmation Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailOrderGuestTemplate model.Str

	// SalesEmailOrderCopyTo => Send Order Email Copy To.
	// Comma-separated
	// Path: sales_email/order/copy_to
	SalesEmailOrderCopyTo model.Str

	// SalesEmailOrderCopyMethod => Send Order Email Copy Method.
	// Path: sales_email/order/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCopyMethod model.Str

	// SalesEmailOrderCommentEnabled => Enabled.
	// Path: sales_email/order_comment/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailOrderCommentEnabled model.Bool

	// SalesEmailOrderCommentIdentity => Order Comment Email Sender.
	// Path: sales_email/order_comment/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderCommentIdentity model.Str

	// SalesEmailOrderCommentTemplate => Order Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentTemplate model.Str

	// SalesEmailOrderCommentGuestTemplate => Order Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentGuestTemplate model.Str

	// SalesEmailOrderCommentCopyTo => Send Order Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/order_comment/copy_to
	SalesEmailOrderCommentCopyTo model.Str

	// SalesEmailOrderCommentCopyMethod => Send Order Comments Email Copy Method.
	// Path: sales_email/order_comment/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCommentCopyMethod model.Str

	// SalesEmailInvoiceEnabled => Enabled.
	// Path: sales_email/invoice/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceEnabled model.Bool

	// SalesEmailInvoiceIdentity => Invoice Email Sender.
	// Path: sales_email/invoice/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceIdentity model.Str

	// SalesEmailInvoiceTemplate => Invoice Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceTemplate model.Str

	// SalesEmailInvoiceGuestTemplate => Invoice Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceGuestTemplate model.Str

	// SalesEmailInvoiceCopyTo => Send Invoice Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice/copy_to
	SalesEmailInvoiceCopyTo model.Str

	// SalesEmailInvoiceCopyMethod => Send Invoice Email Copy Method.
	// Path: sales_email/invoice/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCopyMethod model.Str

	// SalesEmailInvoiceCommentEnabled => Enabled.
	// Path: sales_email/invoice_comment/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceCommentEnabled model.Bool

	// SalesEmailInvoiceCommentIdentity => Invoice Comment Email Sender.
	// Path: sales_email/invoice_comment/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceCommentIdentity model.Str

	// SalesEmailInvoiceCommentTemplate => Invoice Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentTemplate model.Str

	// SalesEmailInvoiceCommentGuestTemplate => Invoice Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentGuestTemplate model.Str

	// SalesEmailInvoiceCommentCopyTo => Send Invoice Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice_comment/copy_to
	SalesEmailInvoiceCommentCopyTo model.Str

	// SalesEmailInvoiceCommentCopyMethod => Send Invoice Comments Email Copy Method.
	// Path: sales_email/invoice_comment/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCommentCopyMethod model.Str

	// SalesEmailShipmentEnabled => Enabled.
	// Path: sales_email/shipment/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailShipmentEnabled model.Bool

	// SalesEmailShipmentIdentity => Shipment Email Sender.
	// Path: sales_email/shipment/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentIdentity model.Str

	// SalesEmailShipmentTemplate => Shipment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentTemplate model.Str

	// SalesEmailShipmentGuestTemplate => Shipment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentGuestTemplate model.Str

	// SalesEmailShipmentCopyTo => Send Shipment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment/copy_to
	SalesEmailShipmentCopyTo model.Str

	// SalesEmailShipmentCopyMethod => Send Shipment Email Copy Method.
	// Path: sales_email/shipment/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCopyMethod model.Str

	// SalesEmailShipmentCommentEnabled => Enabled.
	// Path: sales_email/shipment_comment/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailShipmentCommentEnabled model.Bool

	// SalesEmailShipmentCommentIdentity => Shipment Comment Email Sender.
	// Path: sales_email/shipment_comment/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentCommentIdentity model.Str

	// SalesEmailShipmentCommentTemplate => Shipment Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentTemplate model.Str

	// SalesEmailShipmentCommentGuestTemplate => Shipment Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentGuestTemplate model.Str

	// SalesEmailShipmentCommentCopyTo => Send Shipment Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment_comment/copy_to
	SalesEmailShipmentCommentCopyTo model.Str

	// SalesEmailShipmentCommentCopyMethod => Send Shipment Comments Email Copy Method.
	// Path: sales_email/shipment_comment/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCommentCopyMethod model.Str

	// SalesEmailCreditmemoEnabled => Enabled.
	// Path: sales_email/creditmemo/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoEnabled model.Bool

	// SalesEmailCreditmemoIdentity => Credit Memo Email Sender.
	// Path: sales_email/creditmemo/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoIdentity model.Str

	// SalesEmailCreditmemoTemplate => Credit Memo Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoTemplate model.Str

	// SalesEmailCreditmemoGuestTemplate => Credit Memo Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoGuestTemplate model.Str

	// SalesEmailCreditmemoCopyTo => Send Credit Memo Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo/copy_to
	SalesEmailCreditmemoCopyTo model.Str

	// SalesEmailCreditmemoCopyMethod => Send Credit Memo Email Copy Method.
	// Path: sales_email/creditmemo/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCopyMethod model.Str

	// SalesEmailCreditmemoCommentEnabled => Enabled.
	// Path: sales_email/creditmemo_comment/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoCommentEnabled model.Bool

	// SalesEmailCreditmemoCommentIdentity => Credit Memo Comment Email Sender.
	// Path: sales_email/creditmemo_comment/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoCommentIdentity model.Str

	// SalesEmailCreditmemoCommentTemplate => Credit Memo Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentTemplate model.Str

	// SalesEmailCreditmemoCommentGuestTemplate => Credit Memo Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/guest_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentGuestTemplate model.Str

	// SalesEmailCreditmemoCommentCopyTo => Send Credit Memo Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo_comment/copy_to
	SalesEmailCreditmemoCommentCopyTo model.Str

	// SalesEmailCreditmemoCommentCopyMethod => Send Credit Memo Comments Email Copy Method.
	// Path: sales_email/creditmemo_comment/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCommentCopyMethod model.Str

	// SalesPdfInvoicePutOrderId => Display Order ID in Header.
	// Path: sales_pdf/invoice/put_order_id
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesPdfInvoicePutOrderId model.Bool

	// SalesPdfShipmentPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/shipment/put_order_id
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesPdfShipmentPutOrderId model.Bool

	// SalesPdfCreditmemoPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/creditmemo/put_order_id
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesPdfCreditmemoPutOrderId model.Bool

	// RssOrderStatus => Customer Order Status Notification.
	// Path: rss/order/status
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssOrderStatus model.Bool

	// DevGridAsyncIndexing => Asynchronous indexing.
	// Path: dev/grid/async_indexing
	// BackendModel: Otnegam\Sales\Model\Config\Backend\Grid\AsyncIndexing
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	DevGridAsyncIndexing model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesGeneralHideCustomerIp = model.NewBool(`sales/general/hide_customer_ip`, model.WithPkgCfg(pkgCfg))
	pp.SalesTotalsSortDiscount = model.NewStr(`sales/totals_sort/discount`, model.WithPkgCfg(pkgCfg))
	pp.SalesTotalsSortGrandTotal = model.NewStr(`sales/totals_sort/grand_total`, model.WithPkgCfg(pkgCfg))
	pp.SalesTotalsSortShipping = model.NewStr(`sales/totals_sort/shipping`, model.WithPkgCfg(pkgCfg))
	pp.SalesTotalsSortSubtotal = model.NewStr(`sales/totals_sort/subtotal`, model.WithPkgCfg(pkgCfg))
	pp.SalesTotalsSortTax = model.NewStr(`sales/totals_sort/tax`, model.WithPkgCfg(pkgCfg))
	pp.SalesReorderAllow = model.NewBool(`sales/reorder/allow`, model.WithPkgCfg(pkgCfg))
	pp.SalesIdentityLogo = model.NewStr(`sales/identity/logo`, model.WithPkgCfg(pkgCfg))
	pp.SalesIdentityLogoHtml = model.NewStr(`sales/identity/logo_html`, model.WithPkgCfg(pkgCfg))
	pp.SalesIdentityAddress = model.NewStr(`sales/identity/address`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderActive = model.NewBool(`sales/minimum_order/active`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderAmount = model.NewStr(`sales/minimum_order/amount`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderTaxIncluding = model.NewBool(`sales/minimum_order/tax_including`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderDescription = model.NewStr(`sales/minimum_order/description`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderErrorMessage = model.NewStr(`sales/minimum_order/error_message`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderMultiAddress = model.NewBool(`sales/minimum_order/multi_address`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderMultiAddressDescription = model.NewStr(`sales/minimum_order/multi_address_description`, model.WithPkgCfg(pkgCfg))
	pp.SalesMinimumOrderMultiAddressErrorMessage = model.NewStr(`sales/minimum_order/multi_address_error_message`, model.WithPkgCfg(pkgCfg))
	pp.SalesDashboardUseAggregatedData = model.NewBool(`sales/dashboard/use_aggregated_data`, model.WithPkgCfg(pkgCfg))
	pp.SalesOrdersDeletePendingAfter = model.NewStr(`sales/orders/delete_pending_after`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailGeneralAsyncSending = model.NewBool(`sales_email/general/async_sending`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderEnabled = model.NewBool(`sales_email/order/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderIdentity = model.NewStr(`sales_email/order/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderTemplate = model.NewStr(`sales_email/order/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderGuestTemplate = model.NewStr(`sales_email/order/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCopyTo = model.NewStr(`sales_email/order/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCopyMethod = model.NewStr(`sales_email/order/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentEnabled = model.NewBool(`sales_email/order_comment/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentIdentity = model.NewStr(`sales_email/order_comment/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentTemplate = model.NewStr(`sales_email/order_comment/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentGuestTemplate = model.NewStr(`sales_email/order_comment/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentCopyTo = model.NewStr(`sales_email/order_comment/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailOrderCommentCopyMethod = model.NewStr(`sales_email/order_comment/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceEnabled = model.NewBool(`sales_email/invoice/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceIdentity = model.NewStr(`sales_email/invoice/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceTemplate = model.NewStr(`sales_email/invoice/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceGuestTemplate = model.NewStr(`sales_email/invoice/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCopyTo = model.NewStr(`sales_email/invoice/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCopyMethod = model.NewStr(`sales_email/invoice/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentEnabled = model.NewBool(`sales_email/invoice_comment/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentIdentity = model.NewStr(`sales_email/invoice_comment/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentTemplate = model.NewStr(`sales_email/invoice_comment/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentGuestTemplate = model.NewStr(`sales_email/invoice_comment/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentCopyTo = model.NewStr(`sales_email/invoice_comment/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailInvoiceCommentCopyMethod = model.NewStr(`sales_email/invoice_comment/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentEnabled = model.NewBool(`sales_email/shipment/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentIdentity = model.NewStr(`sales_email/shipment/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentTemplate = model.NewStr(`sales_email/shipment/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentGuestTemplate = model.NewStr(`sales_email/shipment/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCopyTo = model.NewStr(`sales_email/shipment/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCopyMethod = model.NewStr(`sales_email/shipment/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentEnabled = model.NewBool(`sales_email/shipment_comment/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentIdentity = model.NewStr(`sales_email/shipment_comment/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentTemplate = model.NewStr(`sales_email/shipment_comment/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentGuestTemplate = model.NewStr(`sales_email/shipment_comment/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentCopyTo = model.NewStr(`sales_email/shipment_comment/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailShipmentCommentCopyMethod = model.NewStr(`sales_email/shipment_comment/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoEnabled = model.NewBool(`sales_email/creditmemo/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoIdentity = model.NewStr(`sales_email/creditmemo/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoTemplate = model.NewStr(`sales_email/creditmemo/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoGuestTemplate = model.NewStr(`sales_email/creditmemo/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCopyTo = model.NewStr(`sales_email/creditmemo/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCopyMethod = model.NewStr(`sales_email/creditmemo/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentEnabled = model.NewBool(`sales_email/creditmemo_comment/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentIdentity = model.NewStr(`sales_email/creditmemo_comment/identity`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentTemplate = model.NewStr(`sales_email/creditmemo_comment/template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentGuestTemplate = model.NewStr(`sales_email/creditmemo_comment/guest_template`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentCopyTo = model.NewStr(`sales_email/creditmemo_comment/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.SalesEmailCreditmemoCommentCopyMethod = model.NewStr(`sales_email/creditmemo_comment/copy_method`, model.WithPkgCfg(pkgCfg))
	pp.SalesPdfInvoicePutOrderId = model.NewBool(`sales_pdf/invoice/put_order_id`, model.WithPkgCfg(pkgCfg))
	pp.SalesPdfShipmentPutOrderId = model.NewBool(`sales_pdf/shipment/put_order_id`, model.WithPkgCfg(pkgCfg))
	pp.SalesPdfCreditmemoPutOrderId = model.NewBool(`sales_pdf/creditmemo/put_order_id`, model.WithPkgCfg(pkgCfg))
	pp.RssOrderStatus = model.NewBool(`rss/order/status`, model.WithPkgCfg(pkgCfg))
	pp.DevGridAsyncIndexing = model.NewBool(`dev/grid/async_indexing`, model.WithPkgCfg(pkgCfg))

	return pp
}
