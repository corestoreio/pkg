// +build ignore

package checkout

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
	// CheckoutOptionsOnepageCheckoutEnabled => Enable Onepage Checkout.
	// Path: checkout/options/onepage_checkout_enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CheckoutOptionsOnepageCheckoutEnabled model.Bool

	// CheckoutOptionsGuestCheckout => Allow Guest Checkout.
	// Path: checkout/options/guest_checkout
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CheckoutOptionsGuestCheckout model.Bool

	// CheckoutCartDeleteQuoteAfter => Quote Lifetime (days).
	// Path: checkout/cart/delete_quote_after
	CheckoutCartDeleteQuoteAfter model.Str

	// CheckoutCartRedirectToCart => After Adding a Product Redirect to Shopping Cart.
	// Path: checkout/cart/redirect_to_cart
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CheckoutCartRedirectToCart model.Bool

	// CheckoutCartLinkUseQty => Display Cart Summary.
	// Path: checkout/cart_link/use_qty
	// SourceModel: Otnegam\Checkout\Model\Config\Source\Cart\Summary
	CheckoutCartLinkUseQty model.Str

	// CheckoutSidebarDisplay => Display Shopping Cart Sidebar.
	// Path: checkout/sidebar/display
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CheckoutSidebarDisplay model.Bool

	// CheckoutSidebarCount => Maximum Display Recently Added Item(s).
	// Path: checkout/sidebar/count
	CheckoutSidebarCount model.Str

	// CheckoutPaymentFailedIdentity => Payment Failed Email Sender.
	// Path: checkout/payment_failed/identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedIdentity model.Str

	// CheckoutPaymentFailedReceiver => Payment Failed Email Receiver.
	// Path: checkout/payment_failed/receiver
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedReceiver model.Str

	// CheckoutPaymentFailedTemplate => Payment Failed Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: checkout/payment_failed/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	CheckoutPaymentFailedTemplate model.Str

	// CheckoutPaymentFailedCopyTo => Send Payment Failed Email Copy To.
	// Separate by ",".
	// Path: checkout/payment_failed/copy_to
	CheckoutPaymentFailedCopyTo model.Str

	// CheckoutPaymentFailedCopyMethod => Send Payment Failed Email Copy Method.
	// Path: checkout/payment_failed/copy_method
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Method
	CheckoutPaymentFailedCopyMethod model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutOptionsOnepageCheckoutEnabled = model.NewBool(`checkout/options/onepage_checkout_enabled`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutOptionsGuestCheckout = model.NewBool(`checkout/options/guest_checkout`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutCartDeleteQuoteAfter = model.NewStr(`checkout/cart/delete_quote_after`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutCartRedirectToCart = model.NewBool(`checkout/cart/redirect_to_cart`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutCartLinkUseQty = model.NewStr(`checkout/cart_link/use_qty`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutSidebarDisplay = model.NewBool(`checkout/sidebar/display`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutSidebarCount = model.NewStr(`checkout/sidebar/count`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutPaymentFailedIdentity = model.NewStr(`checkout/payment_failed/identity`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutPaymentFailedReceiver = model.NewStr(`checkout/payment_failed/receiver`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutPaymentFailedTemplate = model.NewStr(`checkout/payment_failed/template`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutPaymentFailedCopyTo = model.NewStr(`checkout/payment_failed/copy_to`, model.WithPkgCfg(pkgCfg))
	pp.CheckoutPaymentFailedCopyMethod = model.NewStr(`checkout/payment_failed/copy_method`, model.WithPkgCfg(pkgCfg))

	return pp
}
