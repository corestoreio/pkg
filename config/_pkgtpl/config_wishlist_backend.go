// +build ignore

package wishlist

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
	// WishlistEmailEmailIdentity => Email Sender.
	// Path: wishlist/email/email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	WishlistEmailEmailIdentity model.Str

	// WishlistEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: wishlist/email/email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	WishlistEmailEmailTemplate model.Str

	// WishlistEmailNumberLimit => Max Emails Allowed to be Sent.
	// 10 by default. Max - 10000
	// Path: wishlist/email/number_limit
	WishlistEmailNumberLimit model.Str

	// WishlistEmailTextLimit => Email Text Length Limit.
	// 255 by default
	// Path: wishlist/email/text_limit
	WishlistEmailTextLimit model.Str

	// WishlistGeneralActive => Enabled.
	// Path: wishlist/general/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WishlistGeneralActive model.Bool

	// WishlistWishlistLinkUseQty => Display Wish List Summary.
	// Path: wishlist/wishlist_link/use_qty
	// SourceModel: Magento\Wishlist\Model\Config\Source\Summary
	WishlistWishlistLinkUseQty model.Str

	// RssWishlistActive => Enable RSS.
	// Path: rss/wishlist/active
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssWishlistActive model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.WishlistEmailEmailIdentity = model.NewStr(`wishlist/email/email_identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailEmailTemplate = model.NewStr(`wishlist/email/email_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailNumberLimit = model.NewStr(`wishlist/email/number_limit`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailTextLimit = model.NewStr(`wishlist/email/text_limit`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistGeneralActive = model.NewBool(`wishlist/general/active`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistWishlistLinkUseQty = model.NewStr(`wishlist/wishlist_link/use_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.RssWishlistActive = model.NewBool(`rss/wishlist/active`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
