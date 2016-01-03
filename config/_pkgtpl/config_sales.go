// +build ignore

package sales

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        "sales",
			Label:     `Sales`,
			SortOrder: 300,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Sales::config_sales
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "general",
					Label:     `General`,
					SortOrder: 5,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/general/hide_customer_ip
							ID:      "hide_customer_ip",
							Label:   `Hide Customer IP`,
							Comment: element.LongText(`Choose whether a customer IP is shown in orders, invoices, shipments, and credit memos.`),
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "totals_sort",
					Label:     `Checkout Totals Sort Order`,
					SortOrder: 10,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/totals_sort/discount
							ID:        "discount",
							Label:     `Discount`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   20,
						},

						&element.Field{
							// Path: sales/totals_sort/grand_total
							ID:        "grand_total",
							Label:     `Grand Total`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   100,
						},

						&element.Field{
							// Path: sales/totals_sort/shipping
							ID:        "shipping",
							Label:     `Shipping`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   30,
						},

						&element.Field{
							// Path: sales/totals_sort/subtotal
							ID:        "subtotal",
							Label:     `Subtotal`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   10,
						},

						&element.Field{
							// Path: sales/totals_sort/tax
							ID:        "tax",
							Label:     `Tax`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   40,
						},
					),
				},

				&element.Group{
					ID:        "reorder",
					Label:     `Reorder`,
					SortOrder: 20,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/reorder/allow
							ID:        "allow",
							Label:     `Allow Reorder`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "identity",
					Label:     `Invoice and Packing Slip Design`,
					SortOrder: 40,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/identity/logo
							ID:        "logo",
							Label:     `Logo for PDF Print-outs (200x50)`,
							Comment:   element.LongText(`Your default logo will be used in PDF and HTML documents.<br />(jpeg, tiff, png) If your pdf image is distorted, try to use larger file-size image.`),
							Type:      element.TypeImage,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Pdf
						},

						&element.Field{
							// Path: sales/identity/logo_html
							ID:        "logo_html",
							Label:     `Logo for HTML Print View`,
							Comment:   element.LongText(`Logo for HTML documents only. If empty, default will be used.<br />(jpeg, gif, png)`),
							Type:      element.TypeImage,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Image
						},

						&element.Field{
							// Path: sales/identity/address
							ID:        "address",
							Label:     `Address`,
							Type:      element.TypeTextarea,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "minimum_order",
					Label:     `Minimum Order Amount`,
					SortOrder: 50,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/minimum_order/active
							ID:        "active",
							Label:     `Enable`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales/minimum_order/amount
							ID:        "amount",
							Label:     `Minimum Amount`,
							Comment:   element.LongText(`Subtotal after discount`),
							Type:      element.Type,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: sales/minimum_order/tax_including
							ID:        "tax_including",
							Label:     `Include Tax to Amount`,
							Type:      element.TypeSelect,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales/minimum_order/description
							ID:        "description",
							Label:     `Description Message`,
							Comment:   element.LongText(`This message will be shown in the shopping cart when the subtotal (after discount) is lower than the minimum allowed amount.`),
							Type:      element.TypeTextarea,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales/minimum_order/error_message
							ID:        "error_message",
							Label:     `Error to Show in Shopping Cart`,
							Type:      element.TypeTextarea,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales/minimum_order/multi_address
							ID:        "multi_address",
							Label:     `Validate Each Address Separately in Multi-address Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales/minimum_order/multi_address_description
							ID:        "multi_address_description",
							Label:     `Multi-address Description Message`,
							Comment:   element.LongText(`We'll use the default description above if you leave this empty.`),
							Type:      element.TypeTextarea,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales/minimum_order/multi_address_error_message
							ID:        "multi_address_error_message",
							Label:     `Multi-address Error to Show in Shopping Cart`,
							Comment:   element.LongText(`We'll use the default error above if you leave this empty.`),
							Type:      element.TypeTextarea,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "dashboard",
					Label:     `Dashboard`,
					SortOrder: 60,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/dashboard/use_aggregated_data
							ID:        "use_aggregated_data",
							Label:     `Use Aggregated Data (beta)`,
							Comment:   element.LongText(`Improves dashboard performance but provides non-realtime data.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "orders",
					Label:     `Orders Cron Settings`,
					SortOrder: 70,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/orders/delete_pending_after
							ID:        "delete_pending_after",
							Label:     `Pending Payment Order Lifetime (minutes)`,
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   480,
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "sales_email",
			Label:     `Sales Emails`,
			SortOrder: 301,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Sales::sales_email
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:    "general",
					Label: `General Settings`,
					Scope: scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/general/async_sending
							ID:        "async_sending",
							Label:     `Asynchronous sending`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   false,
							// BackendModel: Otnegam\Sales\Model\Config\Backend\Email\AsyncSending
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},
					),
				},

				&element.Group{
					ID:        "order",
					Label:     `Order`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/order/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/order/identity
							ID:        "identity",
							Label:     `New Order Confirmation Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/order/template
							ID:        "template",
							Label:     `New Order Confirmation Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_order_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/order/guest_template
							ID:        "guest_template",
							Label:     `New Order Confirmation Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_order_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/order/copy_to
							ID:        "copy_to",
							Label:     `Send Order Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/order/copy_method
							ID:        "copy_method",
							Label:     `Send Order Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "order_comment",
					Label:     `Order Comments`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/order_comment/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/order_comment/identity
							ID:        "identity",
							Label:     `Order Comment Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/order_comment/template
							ID:        "template",
							Label:     `Order Comment Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_order_comment_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/order_comment/guest_template
							ID:        "guest_template",
							Label:     `Order Comment Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_order_comment_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/order_comment/copy_to
							ID:        "copy_to",
							Label:     `Send Order Comment Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/order_comment/copy_method
							ID:        "copy_method",
							Label:     `Send Order Comments Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "invoice",
					Label:     `Invoice`,
					SortOrder: 3,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/invoice/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/invoice/identity
							ID:        "identity",
							Label:     `Invoice Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/invoice/template
							ID:        "template",
							Label:     `Invoice Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_invoice_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/invoice/guest_template
							ID:        "guest_template",
							Label:     `Invoice Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_invoice_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/invoice/copy_to
							ID:        "copy_to",
							Label:     `Send Invoice Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/invoice/copy_method
							ID:        "copy_method",
							Label:     `Send Invoice Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "invoice_comment",
					Label:     `Invoice Comments`,
					SortOrder: 4,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/invoice_comment/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/invoice_comment/identity
							ID:        "identity",
							Label:     `Invoice Comment Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/invoice_comment/template
							ID:        "template",
							Label:     `Invoice Comment Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_invoice_comment_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/invoice_comment/guest_template
							ID:        "guest_template",
							Label:     `Invoice Comment Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_invoice_comment_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/invoice_comment/copy_to
							ID:        "copy_to",
							Label:     `Send Invoice Comment Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/invoice_comment/copy_method
							ID:        "copy_method",
							Label:     `Send Invoice Comments Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "shipment",
					Label:     `Shipment`,
					SortOrder: 5,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/shipment/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/shipment/identity
							ID:        "identity",
							Label:     `Shipment Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/shipment/template
							ID:        "template",
							Label:     `Shipment Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_shipment_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/shipment/guest_template
							ID:        "guest_template",
							Label:     `Shipment Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_shipment_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/shipment/copy_to
							ID:        "copy_to",
							Label:     `Send Shipment Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/shipment/copy_method
							ID:        "copy_method",
							Label:     `Send Shipment Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "shipment_comment",
					Label:     `Shipment Comments`,
					SortOrder: 6,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/shipment_comment/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/shipment_comment/identity
							ID:        "identity",
							Label:     `Shipment Comment Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/shipment_comment/template
							ID:        "template",
							Label:     `Shipment Comment Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_shipment_comment_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/shipment_comment/guest_template
							ID:        "guest_template",
							Label:     `Shipment Comment Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_shipment_comment_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/shipment_comment/copy_to
							ID:        "copy_to",
							Label:     `Send Shipment Comment Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/shipment_comment/copy_method
							ID:        "copy_method",
							Label:     `Send Shipment Comments Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "creditmemo",
					Label:     `Credit Memo`,
					SortOrder: 7,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/creditmemo/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/creditmemo/identity
							ID:        "identity",
							Label:     `Credit Memo Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/creditmemo/template
							ID:        "template",
							Label:     `Credit Memo Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_creditmemo_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/creditmemo/guest_template
							ID:        "guest_template",
							Label:     `Credit Memo Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_creditmemo_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/creditmemo/copy_to
							ID:        "copy_to",
							Label:     `Send Credit Memo Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/creditmemo/copy_method
							ID:        "copy_method",
							Label:     `Send Credit Memo Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},

				&element.Group{
					ID:        "creditmemo_comment",
					Label:     `Credit Memo Comments`,
					SortOrder: 8,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_email/creditmemo_comment/enabled
							ID:      "enabled",
							Label:   `Enabled`,
							Type:    element.TypeSelect,
							Visible: element.VisibleYes,
							Scope:   scope.PermAll,
							Default: true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales_email/creditmemo_comment/identity
							ID:        "identity",
							Label:     `Credit Memo Comment Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sales_email/creditmemo_comment/template
							ID:        "template",
							Label:     `Credit Memo Comment Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_creditmemo_comment_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/creditmemo_comment/guest_template
							ID:        "guest_template",
							Label:     `Credit Memo Comment Email Template for Guest`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `sales_email_creditmemo_comment_guest_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sales_email/creditmemo_comment/copy_to
							ID:        "copy_to",
							Label:     `Send Credit Memo Comment Email Copy To`,
							Comment:   element.LongText(`Comma-separated`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sales_email/creditmemo_comment/copy_method
							ID:        "copy_method",
							Label:     `Send Credit Memo Comments Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `bcc`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "sales_pdf",
			Label:     `PDF Print-outs`,
			SortOrder: 302,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Sales::sales_pdf
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "invoice",
					Label:     `Invoice`,
					SortOrder: 10,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_pdf/invoice/put_order_id
							ID:        "put_order_id",
							Label:     `Display Order ID in Header`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "shipment",
					Label:     `Shipment`,
					SortOrder: 20,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_pdf/shipment/put_order_id
							ID:        "put_order_id",
							Label:     `Display Order ID in Header`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "creditmemo",
					Label:     `Credit Memo`,
					SortOrder: 30,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales_pdf/creditmemo/put_order_id
							ID:        "put_order_id",
							Label:     `Display Order ID in Header`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID: "rss",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "order",
					Label:     `Order`,
					SortOrder: 4,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: rss/order/status
							ID:        "status",
							Label:     `Customer Order Status Notification`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
		&element.Section{
			ID: "dev",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "grid",
					Label:     `Grid Settings`,
					SortOrder: 131,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/grid/async_indexing
							ID:        "async_indexing",
							Label:     `Asynchronous indexing`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   false,
							// BackendModel: Otnegam\Sales\Model\Config\Backend\Grid\AsyncIndexing
							// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
