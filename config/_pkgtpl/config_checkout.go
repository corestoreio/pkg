// +build ignore

package checkout

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     `Checkout`,
		SortOrder: 305,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Checkout::checkout
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "options",
				Label:     `Checkout Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/options/onepage_checkout_enabled
						ID:        "onepage_checkout_enabled",
						Label:     `Enable Onepage Checkout`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: checkout/options/guest_checkout
						ID:        "guest_checkout",
						Label:     `Allow Guest Checkout`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "cart",
				Label:     `Shopping Cart`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/cart/delete_quote_after
						ID:        "delete_quote_after",
						Label:     `Quote Lifetime (days)`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   30,
					},

					&config.Field{
						// Path: checkout/cart/redirect_to_cart
						ID:        "redirect_to_cart",
						Label:     `After Adding a Product Redirect to Shopping Cart`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "cart_link",
				Label:     `My Cart Link`,
				SortOrder: 3,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/cart_link/use_qty
						ID:        "use_qty",
						Label:     `Display Cart Summary`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Checkout\Model\Config\Source\Cart\Summary
					},
				),
			},

			&config.Group{
				ID:        "sidebar",
				Label:     `Shopping Cart Sidebar`,
				SortOrder: 4,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/sidebar/display
						ID:        "display",
						Label:     `Display Shopping Cart Sidebar`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: checkout/sidebar/count
						ID:        "count",
						Label:     `Maximum Display Recently Added Item(s)`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   5,
					},
				),
			},

			&config.Group{
				ID:        "payment_failed",
				Label:     `Payment Failed Emails`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/payment_failed/identity
						ID:        "identity",
						Label:     `Payment Failed Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: checkout/payment_failed/receiver
						ID:        "receiver",
						Label:     `Payment Failed Email Receiver`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: checkout/payment_failed/template
						ID:        "template",
						Label:     `Payment Failed Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `checkout_payment_failed_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: checkout/payment_failed/copy_to
						ID:        "copy_to",
						Label:     `Send Payment Failed Email Copy To`,
						Comment:   element.LongText(`Separate by ",".`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: checkout/payment_failed/copy_method
						ID:        "copy_method",
						Label:     `Send Payment Failed Email Copy Method`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
					},
				),
			},
		),
	},
)
