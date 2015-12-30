// +build ignore

package checkout

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCheckoutOptionsOnepageCheckoutEnabled => Enable Onepage Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCheckoutOptionsOnepageCheckoutEnabled = model.NewBool(`checkout/options/onepage_checkout_enabled`)

// PathCheckoutOptionsGuestCheckout => Allow Guest Checkout.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCheckoutOptionsGuestCheckout = model.NewBool(`checkout/options/guest_checkout`)

// PathCheckoutCartDeleteQuoteAfter => Quote Lifetime (days).
var PathCheckoutCartDeleteQuoteAfter = model.NewStr(`checkout/cart/delete_quote_after`)

// PathCheckoutCartRedirectToCart => After Adding a Product Redirect to Shopping Cart.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCheckoutCartRedirectToCart = model.NewBool(`checkout/cart/redirect_to_cart`)

// PathCheckoutCartLinkUseQty => Display Cart Summary.
// SourceModel: Otnegam\Checkout\Model\Config\Source\Cart\Summary
var PathCheckoutCartLinkUseQty = model.NewStr(`checkout/cart_link/use_qty`)

// PathCheckoutSidebarDisplay => Display Shopping Cart Sidebar.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCheckoutSidebarDisplay = model.NewBool(`checkout/sidebar/display`)

// PathCheckoutSidebarCount => Maximum Display Recently Added Item(s).
var PathCheckoutSidebarCount = model.NewStr(`checkout/sidebar/count`)

// PathCheckoutPaymentFailedIdentity => Payment Failed Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCheckoutPaymentFailedIdentity = model.NewStr(`checkout/payment_failed/identity`)

// PathCheckoutPaymentFailedReceiver => Payment Failed Email Receiver.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCheckoutPaymentFailedReceiver = model.NewStr(`checkout/payment_failed/receiver`)

// PathCheckoutPaymentFailedTemplate => Payment Failed Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCheckoutPaymentFailedTemplate = model.NewStr(`checkout/payment_failed/template`)

// PathCheckoutPaymentFailedCopyTo => Send Payment Failed Email Copy To.
// Separate by ",".
var PathCheckoutPaymentFailedCopyTo = model.NewStr(`checkout/payment_failed/copy_to`)

// PathCheckoutPaymentFailedCopyMethod => Send Payment Failed Email Copy Method.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
var PathCheckoutPaymentFailedCopyMethod = model.NewStr(`checkout/payment_failed/copy_method`)
