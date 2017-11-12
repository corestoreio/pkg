// +build ignore

package sales

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SalesGeneralHideCustomerIp => Hide Customer IP.
	// Choose whether a customer IP is shown in orders, invoices, shipments, and
	// credit memos.
	// Path: sales/general/hide_customer_ip
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesGeneralHideCustomerIp cfgmodel.Bool

	// SalesTotalsSortDiscount => Discount.
	// Path: sales/totals_sort/discount
	SalesTotalsSortDiscount cfgmodel.Str

	// SalesTotalsSortGrandTotal => Grand Total.
	// Path: sales/totals_sort/grand_total
	SalesTotalsSortGrandTotal cfgmodel.Str

	// SalesTotalsSortShipping => Shipping.
	// Path: sales/totals_sort/shipping
	SalesTotalsSortShipping cfgmodel.Str

	// SalesTotalsSortSubtotal => Subtotal.
	// Path: sales/totals_sort/subtotal
	SalesTotalsSortSubtotal cfgmodel.Str

	// SalesTotalsSortTax => Tax.
	// Path: sales/totals_sort/tax
	SalesTotalsSortTax cfgmodel.Str

	// SalesReorderAllow => Allow Reorder.
	// Path: sales/reorder/allow
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesReorderAllow cfgmodel.Bool

	// SalesIdentityLogo => Logo for PDF Print-outs (200x50).
	// Your default logo will be used in PDF and HTML documents.(jpeg, tiff, png)
	// If your pdf image is distorted, try to use larger file-size image.
	// Path: sales/identity/logo
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Pdf
	SalesIdentityLogo cfgmodel.Str

	// SalesIdentityLogoHtml => Logo for HTML Print View.
	// Logo for HTML documents only. If empty, default will be used.(jpeg, gif,
	// png)
	// Path: sales/identity/logo_html
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	SalesIdentityLogoHtml cfgmodel.Str

	// SalesIdentityAddress => Address.
	// Path: sales/identity/address
	SalesIdentityAddress cfgmodel.Str

	// SalesMinimumOrderActive => Enable.
	// Path: sales/minimum_order/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderActive cfgmodel.Bool

	// SalesMinimumOrderAmount => Minimum Amount.
	// Subtotal after discount
	// Path: sales/minimum_order/amount
	SalesMinimumOrderAmount cfgmodel.Str

	// SalesMinimumOrderTaxIncluding => Include Tax to Amount.
	// Path: sales/minimum_order/tax_including
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderTaxIncluding cfgmodel.Bool

	// SalesMinimumOrderDescription => Description Message.
	// This message will be shown in the shopping cart when the subtotal (after
	// discount) is lower than the minimum allowed amount.
	// Path: sales/minimum_order/description
	SalesMinimumOrderDescription cfgmodel.Str

	// SalesMinimumOrderErrorMessage => Error to Show in Shopping Cart.
	// Path: sales/minimum_order/error_message
	SalesMinimumOrderErrorMessage cfgmodel.Str

	// SalesMinimumOrderMultiAddress => Validate Each Address Separately in Multi-address Checkout.
	// Path: sales/minimum_order/multi_address
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesMinimumOrderMultiAddress cfgmodel.Bool

	// SalesMinimumOrderMultiAddressDescription => Multi-address Description Message.
	// We'll use the default description above if you leave this empty.
	// Path: sales/minimum_order/multi_address_description
	SalesMinimumOrderMultiAddressDescription cfgmodel.Str

	// SalesMinimumOrderMultiAddressErrorMessage => Multi-address Error to Show in Shopping Cart.
	// We'll use the default error above if you leave this empty.
	// Path: sales/minimum_order/multi_address_error_message
	SalesMinimumOrderMultiAddressErrorMessage cfgmodel.Str

	// SalesDashboardUseAggregatedData => Use Aggregated Data (beta).
	// Improves dashboard performance but provides non-realtime data.
	// Path: sales/dashboard/use_aggregated_data
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesDashboardUseAggregatedData cfgmodel.Bool

	// SalesOrdersDeletePendingAfter => Pending Payment Order Lifetime (minutes).
	// Path: sales/orders/delete_pending_after
	SalesOrdersDeletePendingAfter cfgmodel.Str

	// SalesEmailGeneralAsyncSending => Asynchronous sending.
	// Path: sales_email/general/async_sending
	// BackendModel: Magento\Sales\Model\Config\Backend\Email\AsyncSending
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	SalesEmailGeneralAsyncSending cfgmodel.Bool

	// SalesEmailOrderEnabled => Enabled.
	// Path: sales_email/order/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailOrderEnabled cfgmodel.Bool

	// SalesEmailOrderIdentity => New Order Confirmation Email Sender.
	// Path: sales_email/order/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderIdentity cfgmodel.Str

	// SalesEmailOrderTemplate => New Order Confirmation Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderTemplate cfgmodel.Str

	// SalesEmailOrderGuestTemplate => New Order Confirmation Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderGuestTemplate cfgmodel.Str

	// SalesEmailOrderCopyTo => Send Order Email Copy To.
	// Comma-separated
	// Path: sales_email/order/copy_to
	SalesEmailOrderCopyTo cfgmodel.Str

	// SalesEmailOrderCopyMethod => Send Order Email Copy Method.
	// Path: sales_email/order/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCopyMethod cfgmodel.Str

	// SalesEmailOrderCommentEnabled => Enabled.
	// Path: sales_email/order_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailOrderCommentEnabled cfgmodel.Bool

	// SalesEmailOrderCommentIdentity => Order Comment Email Sender.
	// Path: sales_email/order_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailOrderCommentIdentity cfgmodel.Str

	// SalesEmailOrderCommentTemplate => Order Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentTemplate cfgmodel.Str

	// SalesEmailOrderCommentGuestTemplate => Order Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/order_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailOrderCommentGuestTemplate cfgmodel.Str

	// SalesEmailOrderCommentCopyTo => Send Order Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/order_comment/copy_to
	SalesEmailOrderCommentCopyTo cfgmodel.Str

	// SalesEmailOrderCommentCopyMethod => Send Order Comments Email Copy Method.
	// Path: sales_email/order_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailOrderCommentCopyMethod cfgmodel.Str

	// SalesEmailInvoiceEnabled => Enabled.
	// Path: sales_email/invoice/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceEnabled cfgmodel.Bool

	// SalesEmailInvoiceIdentity => Invoice Email Sender.
	// Path: sales_email/invoice/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceIdentity cfgmodel.Str

	// SalesEmailInvoiceTemplate => Invoice Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceTemplate cfgmodel.Str

	// SalesEmailInvoiceGuestTemplate => Invoice Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceGuestTemplate cfgmodel.Str

	// SalesEmailInvoiceCopyTo => Send Invoice Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice/copy_to
	SalesEmailInvoiceCopyTo cfgmodel.Str

	// SalesEmailInvoiceCopyMethod => Send Invoice Email Copy Method.
	// Path: sales_email/invoice/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCopyMethod cfgmodel.Str

	// SalesEmailInvoiceCommentEnabled => Enabled.
	// Path: sales_email/invoice_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailInvoiceCommentEnabled cfgmodel.Bool

	// SalesEmailInvoiceCommentIdentity => Invoice Comment Email Sender.
	// Path: sales_email/invoice_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailInvoiceCommentIdentity cfgmodel.Str

	// SalesEmailInvoiceCommentTemplate => Invoice Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentTemplate cfgmodel.Str

	// SalesEmailInvoiceCommentGuestTemplate => Invoice Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/invoice_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailInvoiceCommentGuestTemplate cfgmodel.Str

	// SalesEmailInvoiceCommentCopyTo => Send Invoice Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/invoice_comment/copy_to
	SalesEmailInvoiceCommentCopyTo cfgmodel.Str

	// SalesEmailInvoiceCommentCopyMethod => Send Invoice Comments Email Copy Method.
	// Path: sales_email/invoice_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailInvoiceCommentCopyMethod cfgmodel.Str

	// SalesEmailShipmentEnabled => Enabled.
	// Path: sales_email/shipment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailShipmentEnabled cfgmodel.Bool

	// SalesEmailShipmentIdentity => Shipment Email Sender.
	// Path: sales_email/shipment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentIdentity cfgmodel.Str

	// SalesEmailShipmentTemplate => Shipment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentTemplate cfgmodel.Str

	// SalesEmailShipmentGuestTemplate => Shipment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentGuestTemplate cfgmodel.Str

	// SalesEmailShipmentCopyTo => Send Shipment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment/copy_to
	SalesEmailShipmentCopyTo cfgmodel.Str

	// SalesEmailShipmentCopyMethod => Send Shipment Email Copy Method.
	// Path: sales_email/shipment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCopyMethod cfgmodel.Str

	// SalesEmailShipmentCommentEnabled => Enabled.
	// Path: sales_email/shipment_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailShipmentCommentEnabled cfgmodel.Bool

	// SalesEmailShipmentCommentIdentity => Shipment Comment Email Sender.
	// Path: sales_email/shipment_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailShipmentCommentIdentity cfgmodel.Str

	// SalesEmailShipmentCommentTemplate => Shipment Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentTemplate cfgmodel.Str

	// SalesEmailShipmentCommentGuestTemplate => Shipment Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/shipment_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailShipmentCommentGuestTemplate cfgmodel.Str

	// SalesEmailShipmentCommentCopyTo => Send Shipment Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/shipment_comment/copy_to
	SalesEmailShipmentCommentCopyTo cfgmodel.Str

	// SalesEmailShipmentCommentCopyMethod => Send Shipment Comments Email Copy Method.
	// Path: sales_email/shipment_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailShipmentCommentCopyMethod cfgmodel.Str

	// SalesEmailCreditmemoEnabled => Enabled.
	// Path: sales_email/creditmemo/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoEnabled cfgmodel.Bool

	// SalesEmailCreditmemoIdentity => Credit Memo Email Sender.
	// Path: sales_email/creditmemo/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoIdentity cfgmodel.Str

	// SalesEmailCreditmemoTemplate => Credit Memo Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoTemplate cfgmodel.Str

	// SalesEmailCreditmemoGuestTemplate => Credit Memo Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoGuestTemplate cfgmodel.Str

	// SalesEmailCreditmemoCopyTo => Send Credit Memo Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo/copy_to
	SalesEmailCreditmemoCopyTo cfgmodel.Str

	// SalesEmailCreditmemoCopyMethod => Send Credit Memo Email Copy Method.
	// Path: sales_email/creditmemo/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCopyMethod cfgmodel.Str

	// SalesEmailCreditmemoCommentEnabled => Enabled.
	// Path: sales_email/creditmemo_comment/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesEmailCreditmemoCommentEnabled cfgmodel.Bool

	// SalesEmailCreditmemoCommentIdentity => Credit Memo Comment Email Sender.
	// Path: sales_email/creditmemo_comment/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SalesEmailCreditmemoCommentIdentity cfgmodel.Str

	// SalesEmailCreditmemoCommentTemplate => Credit Memo Comment Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentTemplate cfgmodel.Str

	// SalesEmailCreditmemoCommentGuestTemplate => Credit Memo Comment Email Template for Guest.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sales_email/creditmemo_comment/guest_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SalesEmailCreditmemoCommentGuestTemplate cfgmodel.Str

	// SalesEmailCreditmemoCommentCopyTo => Send Credit Memo Comment Email Copy To.
	// Comma-separated
	// Path: sales_email/creditmemo_comment/copy_to
	SalesEmailCreditmemoCommentCopyTo cfgmodel.Str

	// SalesEmailCreditmemoCommentCopyMethod => Send Credit Memo Comments Email Copy Method.
	// Path: sales_email/creditmemo_comment/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	SalesEmailCreditmemoCommentCopyMethod cfgmodel.Str

	// SalesPdfInvoicePutOrderId => Display Order ID in Header.
	// Path: sales_pdf/invoice/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfInvoicePutOrderId cfgmodel.Bool

	// SalesPdfShipmentPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/shipment/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfShipmentPutOrderId cfgmodel.Bool

	// SalesPdfCreditmemoPutOrderId => Display Order ID in Header.
	// Path: sales_pdf/creditmemo/put_order_id
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesPdfCreditmemoPutOrderId cfgmodel.Bool

	// RssOrderStatus => Customer Order Status Notification.
	// Path: rss/order/status
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssOrderStatus cfgmodel.Bool

	// DevGridAsyncIndexing => Asynchronous indexing.
	// Path: dev/grid/async_indexing
	// BackendModel: Magento\Sales\Model\Config\Backend\Grid\AsyncIndexing
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	DevGridAsyncIndexing cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesGeneralHideCustomerIp = cfgmodel.NewBool(`sales/general/hide_customer_ip`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortDiscount = cfgmodel.NewStr(`sales/totals_sort/discount`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortGrandTotal = cfgmodel.NewStr(`sales/totals_sort/grand_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortShipping = cfgmodel.NewStr(`sales/totals_sort/shipping`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortSubtotal = cfgmodel.NewStr(`sales/totals_sort/subtotal`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesTotalsSortTax = cfgmodel.NewStr(`sales/totals_sort/tax`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesReorderAllow = cfgmodel.NewBool(`sales/reorder/allow`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityLogo = cfgmodel.NewStr(`sales/identity/logo`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityLogoHtml = cfgmodel.NewStr(`sales/identity/logo_html`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesIdentityAddress = cfgmodel.NewStr(`sales/identity/address`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderActive = cfgmodel.NewBool(`sales/minimum_order/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderAmount = cfgmodel.NewStr(`sales/minimum_order/amount`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderTaxIncluding = cfgmodel.NewBool(`sales/minimum_order/tax_including`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderDescription = cfgmodel.NewStr(`sales/minimum_order/description`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderErrorMessage = cfgmodel.NewStr(`sales/minimum_order/error_message`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddress = cfgmodel.NewBool(`sales/minimum_order/multi_address`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddressDescription = cfgmodel.NewStr(`sales/minimum_order/multi_address_description`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesMinimumOrderMultiAddressErrorMessage = cfgmodel.NewStr(`sales/minimum_order/multi_address_error_message`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesDashboardUseAggregatedData = cfgmodel.NewBool(`sales/dashboard/use_aggregated_data`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesOrdersDeletePendingAfter = cfgmodel.NewStr(`sales/orders/delete_pending_after`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailGeneralAsyncSending = cfgmodel.NewBool(`sales_email/general/async_sending`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderEnabled = cfgmodel.NewBool(`sales_email/order/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderIdentity = cfgmodel.NewStr(`sales_email/order/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderTemplate = cfgmodel.NewStr(`sales_email/order/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderGuestTemplate = cfgmodel.NewStr(`sales_email/order/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCopyTo = cfgmodel.NewStr(`sales_email/order/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCopyMethod = cfgmodel.NewStr(`sales_email/order/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentEnabled = cfgmodel.NewBool(`sales_email/order_comment/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentIdentity = cfgmodel.NewStr(`sales_email/order_comment/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentTemplate = cfgmodel.NewStr(`sales_email/order_comment/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentGuestTemplate = cfgmodel.NewStr(`sales_email/order_comment/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentCopyTo = cfgmodel.NewStr(`sales_email/order_comment/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailOrderCommentCopyMethod = cfgmodel.NewStr(`sales_email/order_comment/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceEnabled = cfgmodel.NewBool(`sales_email/invoice/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceIdentity = cfgmodel.NewStr(`sales_email/invoice/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceTemplate = cfgmodel.NewStr(`sales_email/invoice/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceGuestTemplate = cfgmodel.NewStr(`sales_email/invoice/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCopyTo = cfgmodel.NewStr(`sales_email/invoice/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCopyMethod = cfgmodel.NewStr(`sales_email/invoice/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentEnabled = cfgmodel.NewBool(`sales_email/invoice_comment/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentIdentity = cfgmodel.NewStr(`sales_email/invoice_comment/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentTemplate = cfgmodel.NewStr(`sales_email/invoice_comment/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentGuestTemplate = cfgmodel.NewStr(`sales_email/invoice_comment/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentCopyTo = cfgmodel.NewStr(`sales_email/invoice_comment/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailInvoiceCommentCopyMethod = cfgmodel.NewStr(`sales_email/invoice_comment/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentEnabled = cfgmodel.NewBool(`sales_email/shipment/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentIdentity = cfgmodel.NewStr(`sales_email/shipment/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentTemplate = cfgmodel.NewStr(`sales_email/shipment/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentGuestTemplate = cfgmodel.NewStr(`sales_email/shipment/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCopyTo = cfgmodel.NewStr(`sales_email/shipment/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCopyMethod = cfgmodel.NewStr(`sales_email/shipment/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentEnabled = cfgmodel.NewBool(`sales_email/shipment_comment/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentIdentity = cfgmodel.NewStr(`sales_email/shipment_comment/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentTemplate = cfgmodel.NewStr(`sales_email/shipment_comment/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentGuestTemplate = cfgmodel.NewStr(`sales_email/shipment_comment/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentCopyTo = cfgmodel.NewStr(`sales_email/shipment_comment/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailShipmentCommentCopyMethod = cfgmodel.NewStr(`sales_email/shipment_comment/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoEnabled = cfgmodel.NewBool(`sales_email/creditmemo/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoIdentity = cfgmodel.NewStr(`sales_email/creditmemo/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoTemplate = cfgmodel.NewStr(`sales_email/creditmemo/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoGuestTemplate = cfgmodel.NewStr(`sales_email/creditmemo/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCopyTo = cfgmodel.NewStr(`sales_email/creditmemo/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCopyMethod = cfgmodel.NewStr(`sales_email/creditmemo/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentEnabled = cfgmodel.NewBool(`sales_email/creditmemo_comment/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentIdentity = cfgmodel.NewStr(`sales_email/creditmemo_comment/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentTemplate = cfgmodel.NewStr(`sales_email/creditmemo_comment/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentGuestTemplate = cfgmodel.NewStr(`sales_email/creditmemo_comment/guest_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentCopyTo = cfgmodel.NewStr(`sales_email/creditmemo_comment/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesEmailCreditmemoCommentCopyMethod = cfgmodel.NewStr(`sales_email/creditmemo_comment/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfInvoicePutOrderId = cfgmodel.NewBool(`sales_pdf/invoice/put_order_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfShipmentPutOrderId = cfgmodel.NewBool(`sales_pdf/shipment/put_order_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesPdfCreditmemoPutOrderId = cfgmodel.NewBool(`sales_pdf/creditmemo/put_order_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.RssOrderStatus = cfgmodel.NewBool(`rss/order/status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DevGridAsyncIndexing = cfgmodel.NewBool(`dev/grid/async_indexing`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
