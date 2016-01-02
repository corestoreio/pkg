// +build ignore

package newsletter

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
	// NewsletterSubscriptionAllowGuestSubscribe => Allow Guest Subscription.
	// Path: newsletter/subscription/allow_guest_subscribe
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionAllowGuestSubscribe model.Bool

	// NewsletterSubscriptionConfirm => Need to Confirm.
	// Path: newsletter/subscription/confirm
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionConfirm model.Bool

	// NewsletterSubscriptionConfirmEmailIdentity => Confirmation Email Sender.
	// Path: newsletter/subscription/confirm_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionConfirmEmailIdentity model.Str

	// NewsletterSubscriptionConfirmEmailTemplate => Confirmation Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/confirm_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionConfirmEmailTemplate model.Str

	// NewsletterSubscriptionSuccessEmailIdentity => Success Email Sender.
	// Path: newsletter/subscription/success_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionSuccessEmailIdentity model.Str

	// NewsletterSubscriptionSuccessEmailTemplate => Success Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/success_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionSuccessEmailTemplate model.Str

	// NewsletterSubscriptionUnEmailIdentity => Unsubscription Email Sender.
	// Path: newsletter/subscription/un_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionUnEmailIdentity model.Str

	// NewsletterSubscriptionUnEmailTemplate => Unsubscription Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/un_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionUnEmailTemplate model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.NewsletterSubscriptionAllowGuestSubscribe = model.NewBool(`newsletter/subscription/allow_guest_subscribe`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionConfirm = model.NewBool(`newsletter/subscription/confirm`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionConfirmEmailIdentity = model.NewStr(`newsletter/subscription/confirm_email_identity`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionConfirmEmailTemplate = model.NewStr(`newsletter/subscription/confirm_email_template`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionSuccessEmailIdentity = model.NewStr(`newsletter/subscription/success_email_identity`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionSuccessEmailTemplate = model.NewStr(`newsletter/subscription/success_email_template`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionUnEmailIdentity = model.NewStr(`newsletter/subscription/un_email_identity`, model.WithPkgCfg(pkgCfg))
	pp.NewsletterSubscriptionUnEmailTemplate = model.NewStr(`newsletter/subscription/un_email_template`, model.WithPkgCfg(pkgCfg))

	return pp
}
