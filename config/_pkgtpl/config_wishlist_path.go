// +build ignore

package wishlist

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathWishlistEmailEmailIdentity => Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathWishlistEmailEmailIdentity = model.NewStr(`wishlist/email/email_identity`)

// PathWishlistEmailEmailTemplate => Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathWishlistEmailEmailTemplate = model.NewStr(`wishlist/email/email_template`)

// PathWishlistEmailNumberLimit => Max Emails Allowed to be Sent.
// 10 by default. Max - 10000
var PathWishlistEmailNumberLimit = model.NewStr(`wishlist/email/number_limit`)

// PathWishlistEmailTextLimit => Email Text Length Limit.
// 255 by default
var PathWishlistEmailTextLimit = model.NewStr(`wishlist/email/text_limit`)

// PathWishlistGeneralActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWishlistGeneralActive = model.NewBool(`wishlist/general/active`)

// PathWishlistWishlistLinkUseQty => Display Wish List Summary.
// SourceModel: Otnegam\Wishlist\Model\Config\Source\Summary
var PathWishlistWishlistLinkUseQty = model.NewStr(`wishlist/wishlist_link/use_qty`)

// PathRssWishlistActive => Enable RSS.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssWishlistActive = model.NewBool(`rss/wishlist/active`)
