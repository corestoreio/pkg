// +build ignore

package newsletter

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathNewsletterSubscriptionAllowGuestSubscribe => Allow Guest Subscription.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathNewsletterSubscriptionAllowGuestSubscribe = model.NewBool(`newsletter/subscription/allow_guest_subscribe`)

// PathNewsletterSubscriptionConfirm => Need to Confirm.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathNewsletterSubscriptionConfirm = model.NewBool(`newsletter/subscription/confirm`)

// PathNewsletterSubscriptionConfirmEmailIdentity => Confirmation Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathNewsletterSubscriptionConfirmEmailIdentity = model.NewStr(`newsletter/subscription/confirm_email_identity`)

// PathNewsletterSubscriptionConfirmEmailTemplate => Confirmation Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathNewsletterSubscriptionConfirmEmailTemplate = model.NewStr(`newsletter/subscription/confirm_email_template`)

// PathNewsletterSubscriptionSuccessEmailIdentity => Success Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathNewsletterSubscriptionSuccessEmailIdentity = model.NewStr(`newsletter/subscription/success_email_identity`)

// PathNewsletterSubscriptionSuccessEmailTemplate => Success Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathNewsletterSubscriptionSuccessEmailTemplate = model.NewStr(`newsletter/subscription/success_email_template`)

// PathNewsletterSubscriptionUnEmailIdentity => Unsubscription Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathNewsletterSubscriptionUnEmailIdentity = model.NewStr(`newsletter/subscription/un_email_identity`)

// PathNewsletterSubscriptionUnEmailTemplate => Unsubscription Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathNewsletterSubscriptionUnEmailTemplate = model.NewStr(`newsletter/subscription/un_email_template`)
