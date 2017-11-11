// +build ignore

package newsletter

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
	// NewsletterSubscriptionAllowGuestSubscribe => Allow Guest Subscription.
	// Path: newsletter/subscription/allow_guest_subscribe
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionAllowGuestSubscribe cfgmodel.Bool

	// NewsletterSubscriptionConfirm => Need to Confirm.
	// Path: newsletter/subscription/confirm
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionConfirm cfgmodel.Bool

	// NewsletterSubscriptionConfirmEmailIdentity => Confirmation Email Sender.
	// Path: newsletter/subscription/confirm_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionConfirmEmailIdentity cfgmodel.Str

	// NewsletterSubscriptionConfirmEmailTemplate => Confirmation Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/confirm_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionConfirmEmailTemplate cfgmodel.Str

	// NewsletterSubscriptionSuccessEmailIdentity => Success Email Sender.
	// Path: newsletter/subscription/success_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionSuccessEmailIdentity cfgmodel.Str

	// NewsletterSubscriptionSuccessEmailTemplate => Success Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/success_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionSuccessEmailTemplate cfgmodel.Str

	// NewsletterSubscriptionUnEmailIdentity => Unsubscription Email Sender.
	// Path: newsletter/subscription/un_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionUnEmailIdentity cfgmodel.Str

	// NewsletterSubscriptionUnEmailTemplate => Unsubscription Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/un_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionUnEmailTemplate cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.NewsletterSubscriptionAllowGuestSubscribe = cfgmodel.NewBool(`newsletter/subscription/allow_guest_subscribe`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirm = cfgmodel.NewBool(`newsletter/subscription/confirm`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirmEmailIdentity = cfgmodel.NewStr(`newsletter/subscription/confirm_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirmEmailTemplate = cfgmodel.NewStr(`newsletter/subscription/confirm_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionSuccessEmailIdentity = cfgmodel.NewStr(`newsletter/subscription/success_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionSuccessEmailTemplate = cfgmodel.NewStr(`newsletter/subscription/success_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionUnEmailIdentity = cfgmodel.NewStr(`newsletter/subscription/un_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionUnEmailTemplate = cfgmodel.NewStr(`newsletter/subscription/un_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
