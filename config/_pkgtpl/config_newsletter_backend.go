// +build ignore

package newsletter

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
	// NewsletterSubscriptionAllowGuestSubscribe => Allow Guest Subscription.
	// Path: newsletter/subscription/allow_guest_subscribe
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionAllowGuestSubscribe model.Bool

	// NewsletterSubscriptionConfirm => Need to Confirm.
	// Path: newsletter/subscription/confirm
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewsletterSubscriptionConfirm model.Bool

	// NewsletterSubscriptionConfirmEmailIdentity => Confirmation Email Sender.
	// Path: newsletter/subscription/confirm_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionConfirmEmailIdentity model.Str

	// NewsletterSubscriptionConfirmEmailTemplate => Confirmation Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/confirm_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionConfirmEmailTemplate model.Str

	// NewsletterSubscriptionSuccessEmailIdentity => Success Email Sender.
	// Path: newsletter/subscription/success_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionSuccessEmailIdentity model.Str

	// NewsletterSubscriptionSuccessEmailTemplate => Success Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/success_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionSuccessEmailTemplate model.Str

	// NewsletterSubscriptionUnEmailIdentity => Unsubscription Email Sender.
	// Path: newsletter/subscription/un_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	NewsletterSubscriptionUnEmailIdentity model.Str

	// NewsletterSubscriptionUnEmailTemplate => Unsubscription Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: newsletter/subscription/un_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	NewsletterSubscriptionUnEmailTemplate model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.NewsletterSubscriptionAllowGuestSubscribe = model.NewBool(`newsletter/subscription/allow_guest_subscribe`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirm = model.NewBool(`newsletter/subscription/confirm`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirmEmailIdentity = model.NewStr(`newsletter/subscription/confirm_email_identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionConfirmEmailTemplate = model.NewStr(`newsletter/subscription/confirm_email_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionSuccessEmailIdentity = model.NewStr(`newsletter/subscription/success_email_identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionSuccessEmailTemplate = model.NewStr(`newsletter/subscription/success_email_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionUnEmailIdentity = model.NewStr(`newsletter/subscription/un_email_identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.NewsletterSubscriptionUnEmailTemplate = model.NewStr(`newsletter/subscription/un_email_template`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
