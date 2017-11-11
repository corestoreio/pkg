// +build ignore

package wishlist

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// WishlistEmailEmailIdentity => Email Sender.
	// Path: wishlist/email/email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	WishlistEmailEmailIdentity cfgmodel.Str

	// WishlistEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: wishlist/email/email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	WishlistEmailEmailTemplate cfgmodel.Str

	// WishlistEmailNumberLimit => Max Emails Allowed to be Sent.
	// 10 by default. Max - 10000
	// Path: wishlist/email/number_limit
	WishlistEmailNumberLimit cfgmodel.Str

	// WishlistEmailTextLimit => Email Text Length Limit.
	// 255 by default
	// Path: wishlist/email/text_limit
	WishlistEmailTextLimit cfgmodel.Str

	// WishlistGeneralActive => Enabled.
	// Path: wishlist/general/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WishlistGeneralActive cfgmodel.Bool

	// WishlistWishlistLinkUseQty => Display Wish List Summary.
	// Path: wishlist/wishlist_link/use_qty
	// SourceModel: Magento\Wishlist\Model\Config\Source\Summary
	WishlistWishlistLinkUseQty cfgmodel.Str

	// RssWishlistActive => Enable RSS.
	// Path: rss/wishlist/active
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssWishlistActive cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.WishlistEmailEmailIdentity = cfgmodel.NewStr(`wishlist/email/email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailEmailTemplate = cfgmodel.NewStr(`wishlist/email/email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailNumberLimit = cfgmodel.NewStr(`wishlist/email/number_limit`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistEmailTextLimit = cfgmodel.NewStr(`wishlist/email/text_limit`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistGeneralActive = cfgmodel.NewBool(`wishlist/general/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WishlistWishlistLinkUseQty = cfgmodel.NewStr(`wishlist/wishlist_link/use_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.RssWishlistActive = cfgmodel.NewBool(`rss/wishlist/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
