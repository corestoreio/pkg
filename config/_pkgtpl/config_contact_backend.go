// +build ignore

package contact

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
	// ContactContactEnabled => Enable Contact Us.
	// Path: contact/contact/enabled
	// BackendModel: Magento\Contact\Model\System\Config\Backend\Links
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	ContactContactEnabled model.Bool

	// ContactEmailRecipientEmail => Send Emails To.
	// Path: contact/email/recipient_email
	ContactEmailRecipientEmail model.Str

	// ContactEmailSenderEmailIdentity => Email Sender.
	// Path: contact/email/sender_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	ContactEmailSenderEmailIdentity model.Str

	// ContactEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: contact/email/email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	ContactEmailEmailTemplate model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.ContactContactEnabled = model.NewBool(`contact/contact/enabled`, model.WithConfigStructure(cfgStruct))
	pp.ContactEmailRecipientEmail = model.NewStr(`contact/email/recipient_email`, model.WithConfigStructure(cfgStruct))
	pp.ContactEmailSenderEmailIdentity = model.NewStr(`contact/email/sender_email_identity`, model.WithConfigStructure(cfgStruct))
	pp.ContactEmailEmailTemplate = model.NewStr(`contact/email/email_template`, model.WithConfigStructure(cfgStruct))

	return pp
}
