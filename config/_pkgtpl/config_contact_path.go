// +build ignore

package contact

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathContactContactEnabled => Enable Contact Us.
// BackendModel: Otnegam\Contact\Model\System\Config\Backend\Links
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathContactContactEnabled = model.NewBool(`contact/contact/enabled`)

// PathContactEmailRecipientEmail => Send Emails To.
var PathContactEmailRecipientEmail = model.NewStr(`contact/email/recipient_email`)

// PathContactEmailSenderEmailIdentity => Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathContactEmailSenderEmailIdentity = model.NewStr(`contact/email/sender_email_identity`)

// PathContactEmailEmailTemplate => Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathContactEmailEmailTemplate = model.NewStr(`contact/email/email_template`)
