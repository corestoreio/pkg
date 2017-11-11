// +build ignore

package contact

import (
	"github.com/corestoreio/cspkg/config/element"
	"github.com/corestoreio/cspkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "contact",
			Label:     `Contacts`,
			SortOrder: 100,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Contact::contact
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "contact",
					Label:     `Contact Us`,
					SortOrder: 10,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: contact/contact/enabled
							ID:        "enabled",
							Label:     `Enable Contact Us`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// BackendModel: Magento\Contact\Model\System\Config\Backend\Links
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "email",
					Label:     `Email Options`,
					SortOrder: 50,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: contact/email/recipient_email
							ID:        "recipient_email",
							Label:     `Send Emails To`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `hello@example.com`,
						},

						element.Field{
							// Path: contact/email/sender_email_identity
							ID:        "sender_email_identity",
							Label:     `Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `custom2`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: contact/email/email_template
							ID:        "email_template",
							Label:     `Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `contact_email_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
