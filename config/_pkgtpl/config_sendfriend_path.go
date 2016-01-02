// +build ignore

package sendfriend

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
	// SendfriendEmailEnabled => Enabled.
	// Path: sendfriend/email/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SendfriendEmailEnabled model.Bool

	// SendfriendEmailTemplate => Select Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sendfriend/email/template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	SendfriendEmailTemplate model.Str

	// SendfriendEmailAllowGuest => Allow for Guests.
	// Path: sendfriend/email/allow_guest
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SendfriendEmailAllowGuest model.Bool

	// SendfriendEmailMaxRecipients => Max Recipients.
	// Path: sendfriend/email/max_recipients
	SendfriendEmailMaxRecipients model.Str

	// SendfriendEmailMaxPerHour => Max Products Sent in 1 Hour.
	// Path: sendfriend/email/max_per_hour
	SendfriendEmailMaxPerHour model.Str

	// SendfriendEmailCheckBy => Limit Sending By.
	// Path: sendfriend/email/check_by
	// SourceModel: Otnegam\SendFriend\Model\Source\Checktype
	SendfriendEmailCheckBy model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SendfriendEmailEnabled = model.NewBool(`sendfriend/email/enabled`, model.WithPkgCfg(pkgCfg))
	pp.SendfriendEmailTemplate = model.NewStr(`sendfriend/email/template`, model.WithPkgCfg(pkgCfg))
	pp.SendfriendEmailAllowGuest = model.NewBool(`sendfriend/email/allow_guest`, model.WithPkgCfg(pkgCfg))
	pp.SendfriendEmailMaxRecipients = model.NewStr(`sendfriend/email/max_recipients`, model.WithPkgCfg(pkgCfg))
	pp.SendfriendEmailMaxPerHour = model.NewStr(`sendfriend/email/max_per_hour`, model.WithPkgCfg(pkgCfg))
	pp.SendfriendEmailCheckBy = model.NewStr(`sendfriend/email/check_by`, model.WithPkgCfg(pkgCfg))

	return pp
}
