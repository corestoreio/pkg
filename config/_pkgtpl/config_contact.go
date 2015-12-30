// +build ignore

package contact

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "contact",
		Label:     `Contacts`,
		SortOrder: 100,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Contact::contact
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "contact",
				Label:     `Contact Us`,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: contact/contact/enabled
						ID:        "enabled",
						Label:     `Enable Contact Us`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// BackendModel: Otnegam\Contact\Model\System\Config\Backend\Links
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "email",
				Label:     `Email Options`,
				SortOrder: 50,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: contact/email/recipient_email
						ID:        "recipient_email",
						Label:     `Send Emails To`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `hello@example.com`,
					},

					&config.Field{
						// Path: contact/email/sender_email_identity
						ID:        "sender_email_identity",
						Label:     `Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `custom2`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: contact/email/email_template
						ID:        "email_template",
						Label:     `Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `contact_email_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},
				),
			},
		),
	},
)
