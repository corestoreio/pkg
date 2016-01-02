// +build ignore

package wishlist

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathWishlistEmailEmailIdentity => Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathWishlistEmailEmailIdentity = model.NewStr(`wishlist/email/email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathWishlistEmailEmailTemplate => Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathWishlistEmailEmailTemplate = model.NewStr(`wishlist/email/email_template`, model.WithPkgCfg(PackageConfiguration))

// PathWishlistEmailNumberLimit => Max Emails Allowed to be Sent.
// 10 by default. Max - 10000
var PathWishlistEmailNumberLimit = model.NewStr(`wishlist/email/number_limit`, model.WithPkgCfg(PackageConfiguration))

// PathWishlistEmailTextLimit => Email Text Length Limit.
// 255 by default
var PathWishlistEmailTextLimit = model.NewStr(`wishlist/email/text_limit`, model.WithPkgCfg(PackageConfiguration))

// PathWishlistGeneralActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWishlistGeneralActive = model.NewBool(`wishlist/general/active`, model.WithPkgCfg(PackageConfiguration))

// PathWishlistWishlistLinkUseQty => Display Wish List Summary.
// SourceModel: Otnegam\Wishlist\Model\Config\Source\Summary
var PathWishlistWishlistLinkUseQty = model.NewStr(`wishlist/wishlist_link/use_qty`, model.WithPkgCfg(PackageConfiguration))

// PathRssWishlistActive => Enable RSS.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssWishlistActive = model.NewBool(`rss/wishlist/active`, model.WithPkgCfg(PackageConfiguration))
