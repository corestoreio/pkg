// +build ignore

package contact

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "contact",
		Label:     "Contacts",
		SortOrder: 100,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "contact",
				Label:     `Contact Us`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `contact/contact/enabled`,
						ID:           "enabled",
						Label:        `Enable Contact Us`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil, // Magento\Contact\Model\System\Config\Backend\Links
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "email",
				Label:     `Email Options`,
				Comment:   ``,
				SortOrder: 50,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `contact/email/recipient_email`,
						ID:           "recipient_email",
						Label:        `Send Emails To`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `hello@example.com`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `contact/email/sender_email_identity`,
						ID:           "sender_email_identity",
						Label:        `Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `custom2`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `contact/email/email_template`,
						ID:           "email_template",
						Label:        `Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `contact_email_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},
)
