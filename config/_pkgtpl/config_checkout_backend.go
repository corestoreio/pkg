// +build ignore

package checkout

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
	// CheckoutOptionsOnepageCheckoutEnabled => Enable Onepage Checkout.
	// Path: checkout/options/onepage_checkout_enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutOptionsOnepageCheckoutEnabled cfgmodel.Bool

	// CheckoutOptionsGuestCheckout => Allow Guest Checkout.
	// Path: checkout/options/guest_checkout
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutOptionsGuestCheckout cfgmodel.Bool

	// CheckoutCartDeleteQuoteAfter => Quote Lifetime (days).
	// Path: checkout/cart/delete_quote_after
	CheckoutCartDeleteQuoteAfter cfgmodel.Str

	// CheckoutCartRedirectToCart => After Adding a Product Redirect to Shopping Cart.
	// Path: checkout/cart/redirect_to_cart
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutCartRedirectToCart cfgmodel.Bool

	// CheckoutCartLinkUseQty => Display Cart Summary.
	// Path: checkout/cart_link/use_qty
	// SourceModel: Magento\Checkout\Model\Config\Source\Cart\Summary
	CheckoutCartLinkUseQty cfgmodel.Str

	// CheckoutSidebarDisplay => Display Shopping Cart Sidebar.
	// Path: checkout/sidebar/display
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutSidebarDisplay cfgmodel.Bool

	// CheckoutSidebarCount => Maximum Display Recently Added Item(s).
	// Path: checkout/sidebar/count
	CheckoutSidebarCount cfgmodel.Str

	// CheckoutPaymentFailedIdentity => Payment Failed Email Sender.
	// Path: checkout/payment_failed/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedIdentity cfgmodel.Str

	// CheckoutPaymentFailedReceiver => Payment Failed Email Receiver.
	// Path: checkout/payment_failed/receiver
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedReceiver cfgmodel.Str

	// CheckoutPaymentFailedTemplate => Payment Failed Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: checkout/payment_failed/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CheckoutPaymentFailedTemplate cfgmodel.Str

	// CheckoutPaymentFailedCopyTo => Send Payment Failed Email Copy To.
	// Separate by ",".
	// Path: checkout/payment_failed/copy_to
	CheckoutPaymentFailedCopyTo cfgmodel.Str

	// CheckoutPaymentFailedCopyMethod => Send Payment Failed Email Copy Method.
	// Path: checkout/payment_failed/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	CheckoutPaymentFailedCopyMethod cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutOptionsOnepageCheckoutEnabled = cfgmodel.NewBool(`checkout/options/onepage_checkout_enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutOptionsGuestCheckout = cfgmodel.NewBool(`checkout/options/guest_checkout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutCartDeleteQuoteAfter = cfgmodel.NewStr(`checkout/cart/delete_quote_after`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutCartRedirectToCart = cfgmodel.NewBool(`checkout/cart/redirect_to_cart`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutCartLinkUseQty = cfgmodel.NewStr(`checkout/cart_link/use_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutSidebarDisplay = cfgmodel.NewBool(`checkout/sidebar/display`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutSidebarCount = cfgmodel.NewStr(`checkout/sidebar/count`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutPaymentFailedIdentity = cfgmodel.NewStr(`checkout/payment_failed/identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutPaymentFailedReceiver = cfgmodel.NewStr(`checkout/payment_failed/receiver`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutPaymentFailedTemplate = cfgmodel.NewStr(`checkout/payment_failed/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutPaymentFailedCopyTo = cfgmodel.NewStr(`checkout/payment_failed/copy_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CheckoutPaymentFailedCopyMethod = cfgmodel.NewStr(`checkout/payment_failed/copy_method`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
