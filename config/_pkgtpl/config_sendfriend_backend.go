// +build ignore

package sendfriend

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

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SendfriendEmailEnabled = model.NewBool(`sendfriend/email/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SendfriendEmailTemplate = model.NewStr(`sendfriend/email/template`, model.WithConfigStructure(cfgStruct))
	pp.SendfriendEmailAllowGuest = model.NewBool(`sendfriend/email/allow_guest`, model.WithConfigStructure(cfgStruct))
	pp.SendfriendEmailMaxRecipients = model.NewStr(`sendfriend/email/max_recipients`, model.WithConfigStructure(cfgStruct))
	pp.SendfriendEmailMaxPerHour = model.NewStr(`sendfriend/email/max_per_hour`, model.WithConfigStructure(cfgStruct))
	pp.SendfriendEmailCheckBy = model.NewStr(`sendfriend/email/check_by`, model.WithConfigStructure(cfgStruct))

	return pp
}
