// +build ignore

package sendfriend

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SendfriendEmailEnabled => Enabled.
	// Path: sendfriend/email/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SendfriendEmailEnabled cfgmodel.Bool

	// SendfriendEmailTemplate => Select Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sendfriend/email/template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SendfriendEmailTemplate cfgmodel.Str

	// SendfriendEmailAllowGuest => Allow for Guests.
	// Path: sendfriend/email/allow_guest
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SendfriendEmailAllowGuest cfgmodel.Bool

	// SendfriendEmailMaxRecipients => Max Recipients.
	// Path: sendfriend/email/max_recipients
	SendfriendEmailMaxRecipients cfgmodel.Str

	// SendfriendEmailMaxPerHour => Max Products Sent in 1 Hour.
	// Path: sendfriend/email/max_per_hour
	SendfriendEmailMaxPerHour cfgmodel.Str

	// SendfriendEmailCheckBy => Limit Sending By.
	// Path: sendfriend/email/check_by
	// SourceModel: Magento\SendFriend\Model\Source\Checktype
	SendfriendEmailCheckBy cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SendfriendEmailEnabled = cfgmodel.NewBool(`sendfriend/email/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SendfriendEmailTemplate = cfgmodel.NewStr(`sendfriend/email/template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SendfriendEmailAllowGuest = cfgmodel.NewBool(`sendfriend/email/allow_guest`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SendfriendEmailMaxRecipients = cfgmodel.NewStr(`sendfriend/email/max_recipients`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SendfriendEmailMaxPerHour = cfgmodel.NewStr(`sendfriend/email/max_per_hour`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SendfriendEmailCheckBy = cfgmodel.NewStr(`sendfriend/email/check_by`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
