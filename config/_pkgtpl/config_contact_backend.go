// +build ignore

package contact

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// ContactContactEnabled => Enable Contact Us.
	// Path: contact/contact/enabled
	// BackendModel: Magento\Contact\Model\System\Config\Backend\Links
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	ContactContactEnabled cfgmodel.Bool

	// ContactEmailRecipientEmail => Send Emails To.
	// Path: contact/email/recipient_email
	ContactEmailRecipientEmail cfgmodel.Str

	// ContactEmailSenderEmailIdentity => Email Sender.
	// Path: contact/email/sender_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	ContactEmailSenderEmailIdentity cfgmodel.Str

	// ContactEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: contact/email/email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	ContactEmailEmailTemplate cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.ContactContactEnabled = cfgmodel.NewBool(`contact/contact/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ContactEmailRecipientEmail = cfgmodel.NewStr(`contact/email/recipient_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ContactEmailSenderEmailIdentity = cfgmodel.NewStr(`contact/email/sender_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.ContactEmailEmailTemplate = cfgmodel.NewStr(`contact/email/email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
