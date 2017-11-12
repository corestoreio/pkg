// +build ignore

package checkout

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "checkout",
			Label:     `Checkout`,
			SortOrder: 305,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Checkout::checkout
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "options",
					Label:     `Checkout Options`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: checkout/options/onepage_checkout_enabled
							ID:        "onepage_checkout_enabled",
							Label:     `Enable Onepage Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: checkout/options/guest_checkout
							ID:        "guest_checkout",
							Label:     `Allow Guest Checkout`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "cart",
					Label:     `Shopping Cart`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: checkout/cart/delete_quote_after
							ID:        "delete_quote_after",
							Label:     `Quote Lifetime (days)`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   30,
						},

						element.Field{
							// Path: checkout/cart/redirect_to_cart
							ID:        "redirect_to_cart",
							Label:     `After Adding a Product Redirect to Shopping Cart`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "cart_link",
					Label:     `My Cart Link`,
					SortOrder: 3,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: checkout/cart_link/use_qty
							ID:        "use_qty",
							Label:     `Display Cart Summary`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Checkout\Model\Config\Source\Cart\Summary
						},
					),
				},

				element.Group{
					ID:        "sidebar",
					Label:     `Shopping Cart Sidebar`,
					SortOrder: 4,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: checkout/sidebar/display
							ID:        "display",
							Label:     `Display Shopping Cart Sidebar`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: checkout/sidebar/count
							ID:        "count",
							Label:     `Maximum Display Recently Added Item(s)`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   5,
						},
					),
				},

				element.Group{
					ID:        "payment_failed",
					Label:     `Payment Failed Emails`,
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: checkout/payment_failed/identity
							ID:        "identity",
							Label:     `Payment Failed Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: checkout/payment_failed/receiver
							ID:        "receiver",
							Label:     `Payment Failed Email Receiver`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: checkout/payment_failed/template
							ID:        "template",
							Label:     `Payment Failed Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `checkout_payment_failed_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: checkout/payment_failed/copy_to
							ID:        "copy_to",
							Label:     `Send Payment Failed Email Copy To`,
							Comment:   text.Long(`Separate by ",".`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: checkout/payment_failed/copy_method
							ID:        "copy_method",
							Label:     `Send Payment Failed Email Copy Method`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Method
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
