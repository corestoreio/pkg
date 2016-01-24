// +build ignore

package checkout

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CheckoutOptionsOnepageCheckoutEnabled => Enable Onepage Checkout.
	// Path: checkout/options/onepage_checkout_enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutOptionsOnepageCheckoutEnabled model.Bool

	// CheckoutOptionsGuestCheckout => Allow Guest Checkout.
	// Path: checkout/options/guest_checkout
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutOptionsGuestCheckout model.Bool

	// CheckoutCartDeleteQuoteAfter => Quote Lifetime (days).
	// Path: checkout/cart/delete_quote_after
	CheckoutCartDeleteQuoteAfter model.Str

	// CheckoutCartRedirectToCart => After Adding a Product Redirect to Shopping Cart.
	// Path: checkout/cart/redirect_to_cart
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutCartRedirectToCart model.Bool

	// CheckoutCartLinkUseQty => Display Cart Summary.
	// Path: checkout/cart_link/use_qty
	// SourceModel: Magento\Checkout\Model\Config\Source\Cart\Summary
	CheckoutCartLinkUseQty model.Str

	// CheckoutSidebarDisplay => Display Shopping Cart Sidebar.
	// Path: checkout/sidebar/display
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CheckoutSidebarDisplay model.Bool

	// CheckoutSidebarCount => Maximum Display Recently Added Item(s).
	// Path: checkout/sidebar/count
	CheckoutSidebarCount model.Str

	// CheckoutPaymentFailedIdentity => Payment Failed Email Sender.
	// Path: checkout/payment_failed/identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedIdentity model.Str

	// CheckoutPaymentFailedReceiver => Payment Failed Email Receiver.
	// Path: checkout/payment_failed/receiver
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CheckoutPaymentFailedReceiver model.Str

	// CheckoutPaymentFailedTemplate => Payment Failed Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: checkout/payment_failed/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CheckoutPaymentFailedTemplate model.Str

	// CheckoutPaymentFailedCopyTo => Send Payment Failed Email Copy To.
	// Separate by ",".
	// Path: checkout/payment_failed/copy_to
	CheckoutPaymentFailedCopyTo model.Str

	// CheckoutPaymentFailedCopyMethod => Send Payment Failed Email Copy Method.
	// Path: checkout/payment_failed/copy_method
	// SourceModel: Magento\Config\Model\Config\Source\Email\Method
	CheckoutPaymentFailedCopyMethod model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CheckoutOptionsOnepageCheckoutEnabled = model.NewBool(`checkout/options/onepage_checkout_enabled`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutOptionsGuestCheckout = model.NewBool(`checkout/options/guest_checkout`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutCartDeleteQuoteAfter = model.NewStr(`checkout/cart/delete_quote_after`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutCartRedirectToCart = model.NewBool(`checkout/cart/redirect_to_cart`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutCartLinkUseQty = model.NewStr(`checkout/cart_link/use_qty`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutSidebarDisplay = model.NewBool(`checkout/sidebar/display`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutSidebarCount = model.NewStr(`checkout/sidebar/count`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutPaymentFailedIdentity = model.NewStr(`checkout/payment_failed/identity`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutPaymentFailedReceiver = model.NewStr(`checkout/payment_failed/receiver`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutPaymentFailedTemplate = model.NewStr(`checkout/payment_failed/template`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutPaymentFailedCopyTo = model.NewStr(`checkout/payment_failed/copy_to`, model.WithConfigStructure(cfgStruct))
	pp.CheckoutPaymentFailedCopyMethod = model.NewStr(`checkout/payment_failed/copy_method`, model.WithConfigStructure(cfgStruct))

	return pp
}
