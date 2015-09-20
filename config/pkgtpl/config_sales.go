// +build ignore

package sales

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sales",
		Label:     "Sales",
		SortOrder: 300,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "general",
				Label:     `General`,
				Comment:   ``,
				SortOrder: 5,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/general/hide_customer_ip`,
						ID:           "hide_customer_ip",
						Label:        `Hide Customer IP`,
						Comment:      `Choose whether a customer IP is shown in orders, invoices, shipments, and credit memos.`,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "totals_sort",
				Label:     `Checkout Totals Sort Order`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/totals_sort/discount`,
						ID:           "discount",
						Label:        `Discount`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      20,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/totals_sort/grand_total`,
						ID:           "grand_total",
						Label:        `Grand Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      100,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/totals_sort/shipping`,
						ID:           "shipping",
						Label:        `Shipping`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      30,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/totals_sort/subtotal`,
						ID:           "subtotal",
						Label:        `Subtotal`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      10,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/totals_sort/tax`,
						ID:           "tax",
						Label:        `Tax`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      40,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "reorder",
				Label:     `Reorder`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/reorder/allow`,
						ID:           "allow",
						Label:        `Allow Reorder`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "identity",
				Label:     `Invoice and Packing Slip Design`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/identity/logo`,
						ID:           "logo",
						Label:        `Logo for PDF Print-outs (200x50)`,
						Comment:      `Your default logo will be used in PDF and HTML documents.<br />(jpeg, tiff, png) If your pdf image is distorted, try to use larger file-size image.`,
						Type:         config.TypeImage,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image\Pdf
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/identity/logo_html`,
						ID:           "logo_html",
						Label:        `Logo for HTML Print View`,
						Comment:      `Logo for HTML documents only. If empty, default will be used.<br />(jpeg, gif, png)`,
						Type:         config.TypeImage,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/identity/address`,
						ID:           "address",
						Label:        `Address`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "minimum_order",
				Label:     `Minimum Order Amount`,
				Comment:   ``,
				SortOrder: 50,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/minimum_order/active`,
						ID:           "active",
						Label:        `Enable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales/minimum_order/amount`,
						ID:           "amount",
						Label:        `Minimum Amount`,
						Comment:      `Subtotal after discount`,
						Type:         config.Type,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/minimum_order/tax_including`,
						ID:           "tax_including",
						Label:        `Include Tax to Amount`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    15,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales/minimum_order/description`,
						ID:           "description",
						Label:        `Description Message`,
						Comment:      `This message will be shown in the shopping cart when the subtotal (after discount) is lower than the minimum allowed amount.`,
						Type:         config.TypeTextarea,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/minimum_order/error_message`,
						ID:           "error_message",
						Label:        `Error to Show in Shopping Cart`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/minimum_order/multi_address`,
						ID:           "multi_address",
						Label:        `Validate Each Address Separately in Multi-address Checkout`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales/minimum_order/multi_address_description`,
						ID:           "multi_address_description",
						Label:        `Multi-address Description Message`,
						Comment:      `We'll use the default description above if you leave this empty.`,
						Type:         config.TypeTextarea,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/minimum_order/multi_address_error_message`,
						ID:           "multi_address_error_message",
						Label:        `Multi-address Error to Show in Shopping Cart`,
						Comment:      `We'll use the default error above if you leave this empty.`,
						Type:         config.TypeTextarea,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "dashboard",
				Label:     `Dashboard`,
				Comment:   ``,
				SortOrder: 60,
				Scope:     config.NewScopePerm(config.IDScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/dashboard/use_aggregated_data`,
						ID:           "use_aggregated_data",
						Label:        `Use Aggregated Data (beta)`,
						Comment:      `Improves dashboard performance but provides non-realtime data.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "sales_email",
		Label:     "Sales Emails",
		SortOrder: 301,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "general",
				Label:     `General Settings`,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(config.IDScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/general/async_sending`,
						ID:           "async_sending",
						Label:        `Asynchronous sending`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      false,
						BackendModel: nil, // Magento\Sales\Model\Config\Backend\Email\AsyncSending
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},

			&config.Group{
				ID:        "order",
				Label:     `Order`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/order/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/order/identity`,
						ID:           "identity",
						Label:        `New Order Confirmation Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/order/template`,
						ID:           "template",
						Label:        `New Order Confirmation Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_order_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/order/guest_template`,
						ID:           "guest_template",
						Label:        `New Order Confirmation Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_order_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/order/copy_to`,
						ID:           "copy_to",
						Label:        `Send Order Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/order/copy_method`,
						ID:           "copy_method",
						Label:        `Send Order Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "order_comment",
				Label:     `Order Comments`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/order_comment/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/order_comment/identity`,
						ID:           "identity",
						Label:        `Order Comment Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/order_comment/template`,
						ID:           "template",
						Label:        `Order Comment Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_order_comment_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/order_comment/guest_template`,
						ID:           "guest_template",
						Label:        `Order Comment Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_order_comment_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/order_comment/copy_to`,
						ID:           "copy_to",
						Label:        `Send Order Comment Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/order_comment/copy_method`,
						ID:           "copy_method",
						Label:        `Send Order Comments Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "invoice",
				Label:     `Invoice`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/invoice/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/invoice/identity`,
						ID:           "identity",
						Label:        `Invoice Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/invoice/template`,
						ID:           "template",
						Label:        `Invoice Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_invoice_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/invoice/guest_template`,
						ID:           "guest_template",
						Label:        `Invoice Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_invoice_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/invoice/copy_to`,
						ID:           "copy_to",
						Label:        `Send Invoice Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/invoice/copy_method`,
						ID:           "copy_method",
						Label:        `Send Invoice Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "invoice_comment",
				Label:     `Invoice Comments`,
				Comment:   ``,
				SortOrder: 4,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/invoice_comment/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/invoice_comment/identity`,
						ID:           "identity",
						Label:        `Invoice Comment Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/invoice_comment/template`,
						ID:           "template",
						Label:        `Invoice Comment Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_invoice_comment_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/invoice_comment/guest_template`,
						ID:           "guest_template",
						Label:        `Invoice Comment Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_invoice_comment_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/invoice_comment/copy_to`,
						ID:           "copy_to",
						Label:        `Send Invoice Comment Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/invoice_comment/copy_method`,
						ID:           "copy_method",
						Label:        `Send Invoice Comments Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "shipment",
				Label:     `Shipment`,
				Comment:   ``,
				SortOrder: 5,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/shipment/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/shipment/identity`,
						ID:           "identity",
						Label:        `Shipment Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/shipment/template`,
						ID:           "template",
						Label:        `Shipment Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_shipment_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/shipment/guest_template`,
						ID:           "guest_template",
						Label:        `Shipment Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_shipment_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/shipment/copy_to`,
						ID:           "copy_to",
						Label:        `Send Shipment Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/shipment/copy_method`,
						ID:           "copy_method",
						Label:        `Send Shipment Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "shipment_comment",
				Label:     `Shipment Comments`,
				Comment:   ``,
				SortOrder: 6,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/shipment_comment/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/shipment_comment/identity`,
						ID:           "identity",
						Label:        `Shipment Comment Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/shipment_comment/template`,
						ID:           "template",
						Label:        `Shipment Comment Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_shipment_comment_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/shipment_comment/guest_template`,
						ID:           "guest_template",
						Label:        `Shipment Comment Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_shipment_comment_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/shipment_comment/copy_to`,
						ID:           "copy_to",
						Label:        `Send Shipment Comment Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/shipment_comment/copy_method`,
						ID:           "copy_method",
						Label:        `Send Shipment Comments Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "creditmemo",
				Label:     `Credit Memo`,
				Comment:   ``,
				SortOrder: 7,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/creditmemo/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/creditmemo/identity`,
						ID:           "identity",
						Label:        `Credit Memo Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/creditmemo/template`,
						ID:           "template",
						Label:        `Credit Memo Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_creditmemo_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/creditmemo/guest_template`,
						ID:           "guest_template",
						Label:        `Credit Memo Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_creditmemo_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/creditmemo/copy_to`,
						ID:           "copy_to",
						Label:        `Send Credit Memo Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/creditmemo/copy_method`,
						ID:           "copy_method",
						Label:        `Send Credit Memo Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},

			&config.Group{
				ID:        "creditmemo_comment",
				Label:     `Credit Memo Comments`,
				Comment:   ``,
				SortOrder: 8,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_email/creditmemo_comment/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    0,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales_email/creditmemo_comment/identity`,
						ID:           "identity",
						Label:        `Credit Memo Comment Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sales_email/creditmemo_comment/template`,
						ID:           "template",
						Label:        `Credit Memo Comment Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_creditmemo_comment_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/creditmemo_comment/guest_template`,
						ID:           "guest_template",
						Label:        `Credit Memo Comment Email Template for Guest`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sales_email_creditmemo_comment_guest_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sales_email/creditmemo_comment/copy_to`,
						ID:           "copy_to",
						Label:        `Send Credit Memo Comment Email Copy To`,
						Comment:      `Comma-separated`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales_email/creditmemo_comment/copy_method`,
						ID:           "copy_method",
						Label:        `Send Credit Memo Comments Email Copy Method`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `bcc`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Method
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "sales_pdf",
		Label:     "PDF Print-outs",
		SortOrder: 302,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "invoice",
				Label:     `Invoice`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_pdf/invoice/put_order_id`,
						ID:           "put_order_id",
						Label:        `Display Order ID in Header`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "shipment",
				Label:     `Shipment`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_pdf/shipment/put_order_id`,
						ID:           "put_order_id",
						Label:        `Display Order ID in Header`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "creditmemo",
				Label:     `Credit Memo`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales_pdf/creditmemo/put_order_id`,
						ID:           "put_order_id",
						Label:        `Display Order ID in Header`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "rss",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "order",
				Label:     `Order`,
				Comment:   ``,
				SortOrder: 4,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/order/status`,
						ID:           "status",
						Label:        `Customer Order Status Notification`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "dev",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "grid",
				Label:     `Grid Settings`,
				Comment:   ``,
				SortOrder: 131,
				Scope:     config.NewScopePerm(config.IDScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/grid/async_indexing`,
						ID:           "async_indexing",
						Label:        `Asynchronous indexing`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      false,
						BackendModel: nil, // Magento\Sales\Model\Config\Backend\Grid\AsyncIndexing
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},
)
