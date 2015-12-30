// +build ignore

package sendfriend

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sendfriend",
		Label:     `Email to a Friend`,
		SortOrder: 120,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::sendfriend
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "email",
				Label:     `Email Templates`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sendfriend/email/enabled
						ID:        "enabled",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: sendfriend/email/template
						ID:        "template",
						Label:     `Select Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `sendfriend_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: sendfriend/email/allow_guest
						ID:        "allow_guest",
						Label:     `Allow for Guests`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: sendfriend/email/max_recipients
						ID:        "max_recipients",
						Label:     `Max Recipients`,
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   5,
					},

					&config.Field{
						// Path: sendfriend/email/max_per_hour
						ID:        "max_per_hour",
						Label:     `Max Products Sent in 1 Hour`,
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   5,
					},

					&config.Field{
						// Path: sendfriend/email/check_by
						ID:        "check_by",
						Label:     `Limit Sending By`,
						Type:      config.TypeSelect,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\SendFriend\Model\Source\Checktype
					},
				),
			},
		),
	},
)
