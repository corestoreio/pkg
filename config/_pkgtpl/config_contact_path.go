// +build ignore

package contact

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
	// ContactContactEnabled => Enable Contact Us.
	// Path: contact/contact/enabled
	// BackendModel: Otnegam\Contact\Model\System\Config\Backend\Links
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	ContactContactEnabled model.Bool

	// ContactEmailRecipientEmail => Send Emails To.
	// Path: contact/email/recipient_email
	ContactEmailRecipientEmail model.Str

	// ContactEmailSenderEmailIdentity => Email Sender.
	// Path: contact/email/sender_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	ContactEmailSenderEmailIdentity model.Str

	// ContactEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: contact/email/email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	ContactEmailEmailTemplate model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.ContactContactEnabled = model.NewBool(`contact/contact/enabled`, model.WithPkgCfg(pkgCfg))
	pp.ContactEmailRecipientEmail = model.NewStr(`contact/email/recipient_email`, model.WithPkgCfg(pkgCfg))
	pp.ContactEmailSenderEmailIdentity = model.NewStr(`contact/email/sender_email_identity`, model.WithPkgCfg(pkgCfg))
	pp.ContactEmailEmailTemplate = model.NewStr(`contact/email/email_template`, model.WithPkgCfg(pkgCfg))

	return pp
}
