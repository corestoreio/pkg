// +build ignore

package offlinepayments

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "payment",
		SortOrder: 400,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "checkmo",
				Label:     `Check / Money Order`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/checkmo/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: payment/checkmo/order_status
						ID:        "order_status",
						Label:     `New Order Status`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `pending`,
						// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: payment/checkmo/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/checkmo/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Check / Money order`,
					},

					&config.Field{
						// Path: payment/checkmo/allowspecific
						ID:        "allowspecific",
						Label:     `Payment from Applicable Countries`,
						Type:      config.TypeAllowspecific,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: payment/checkmo/specificcountry
						ID:         "specificcountry",
						Label:      `Payment from Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  51,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: payment/checkmo/payable_to
						ID:        "payable_to",
						Label:     `Make Check Payable to`,
						Type:      config.Type,
						SortOrder: 61,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: payment/checkmo/mailing_address
						ID:        "mailing_address",
						Label:     `Send Check to`,
						Type:      config.TypeTextarea,
						SortOrder: 62,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: payment/checkmo/min_order_total
						ID:        "min_order_total",
						Label:     `Minimum Order Total`,
						Type:      config.TypeText,
						SortOrder: 98,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/checkmo/max_order_total
						ID:        "max_order_total",
						Label:     `Maximum Order Total`,
						Type:      config.TypeText,
						SortOrder: 99,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/checkmo/model
						ID:      "model",
						Type:    config.Type,
						Visible: config.VisibleYes,
						Default: `Otnegam\OfflinePayments\Model\Checkmo`,
					},
				),
			},

			&config.Group{
				ID:        "purchaseorder",
				Label:     `Purchase Order`,
				SortOrder: 32,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/purchaseorder/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: payment/purchaseorder/order_status
						ID:        "order_status",
						Label:     `New Order Status`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `pending`,
						// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: payment/purchaseorder/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/purchaseorder/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Purchase Order`,
					},

					&config.Field{
						// Path: payment/purchaseorder/allowspecific
						ID:        "allowspecific",
						Label:     `Payment from Applicable Countries`,
						Type:      config.TypeAllowspecific,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: payment/purchaseorder/specificcountry
						ID:         "specificcountry",
						Label:      `Payment from Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  51,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: payment/purchaseorder/min_order_total
						ID:        "min_order_total",
						Label:     `Minimum Order Total`,
						Type:      config.TypeText,
						SortOrder: 98,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/purchaseorder/max_order_total
						ID:        "max_order_total",
						Label:     `Maximum Order Total`,
						Type:      config.TypeText,
						SortOrder: 99,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/purchaseorder/model
						ID:      "model",
						Type:    config.Type,
						Visible: config.VisibleYes,
						Default: `Otnegam\OfflinePayments\Model\Purchaseorder`,
					},
				),
			},

			&config.Group{
				ID:        "banktransfer",
				Label:     `Bank Transfer Payment`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/banktransfer/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: payment/banktransfer/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Bank Transfer Payment`,
					},

					&config.Field{
						// Path: payment/banktransfer/order_status
						ID:        "order_status",
						Label:     `New Order Status`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `pending`,
						// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: payment/banktransfer/allowspecific
						ID:        "allowspecific",
						Label:     `Payment from Applicable Countries`,
						Type:      config.TypeAllowspecific,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: payment/banktransfer/specificcountry
						ID:         "specificcountry",
						Label:      `Payment from Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  51,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: payment/banktransfer/instructions
						ID:        "instructions",
						Label:     `Instructions`,
						Type:      config.TypeTextarea,
						SortOrder: 62,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: payment/banktransfer/min_order_total
						ID:        "min_order_total",
						Label:     `Minimum Order Total`,
						Type:      config.TypeText,
						SortOrder: 98,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/banktransfer/max_order_total
						ID:        "max_order_total",
						Label:     `Maximum Order Total`,
						Type:      config.TypeText,
						SortOrder: 99,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/banktransfer/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},
				),
			},

			&config.Group{
				ID:        "cashondelivery",
				Label:     `Cash On Delivery Payment`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/cashondelivery/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: payment/cashondelivery/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Cash On Delivery`,
					},

					&config.Field{
						// Path: payment/cashondelivery/order_status
						ID:        "order_status",
						Label:     `New Order Status`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `pending`,
						// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
					},

					&config.Field{
						// Path: payment/cashondelivery/allowspecific
						ID:        "allowspecific",
						Label:     `Payment from Applicable Countries`,
						Type:      config.TypeAllowspecific,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: payment/cashondelivery/specificcountry
						ID:         "specificcountry",
						Label:      `Payment from Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  51,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: payment/cashondelivery/instructions
						ID:        "instructions",
						Label:     `Instructions`,
						Type:      config.TypeTextarea,
						SortOrder: 62,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: payment/cashondelivery/min_order_total
						ID:        "min_order_total",
						Label:     `Minimum Order Total`,
						Type:      config.TypeText,
						SortOrder: 98,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/cashondelivery/max_order_total
						ID:        "max_order_total",
						Label:     `Maximum Order Total`,
						Type:      config.TypeText,
						SortOrder: 99,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/cashondelivery/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},
				),
			},

			&config.Group{
				ID:        "free",
				Label:     `Zero Subtotal Checkout`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/free/active
						ID:        "active",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: payment/free/order_status
						ID:        "order_status",
						Label:     `New Order Status`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Newprocessing
					},

					&config.Field{
						// Path: payment/free/payment_action
						ID:        "payment_action",
						Label:     `Automatically Invoice All Items`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Source\Invoice
					},

					&config.Field{
						// Path: payment/free/sort_order
						ID:        "sort_order",
						Label:     `Sort Order`,
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: payment/free/title
						ID:        "title",
						Label:     `Title`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: payment/free/allowspecific
						ID:        "allowspecific",
						Label:     `Payment from Applicable Countries`,
						Type:      config.TypeAllowspecific,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: payment/free/specificcountry
						ID:         "specificcountry",
						Label:      `Payment from Specific Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  51,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: payment/free/model
						ID:      "model",
						Type:    config.Type,
						Visible: config.VisibleYes,
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "checkmo",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/checkmo/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `offline`,
					},
				),
			},

			&config.Group{
				ID: "purchaseorder",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/purchaseorder/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `offline`,
					},
				),
			},

			&config.Group{
				ID: "banktransfer",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/banktransfer/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\OfflinePayments\Model\Banktransfer`,
					},

					&config.Field{
						// Path: payment/banktransfer/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `offline`,
					},
				),
			},

			&config.Group{
				ID: "cashondelivery",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/cashondelivery/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\OfflinePayments\Model\Cashondelivery`,
					},

					&config.Field{
						// Path: payment/cashondelivery/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `offline`,
					},
				),
			},

			&config.Group{
				ID: "free",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/free/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `offline`,
					},
				),
			},
		),
	},
)
