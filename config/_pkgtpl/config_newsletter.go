// +build ignore

package newsletter

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "newsletter",
			Label:     `Newsletter`,
			SortOrder: 110,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Newsletter::newsletter
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "subscription",
					Label:     `Subscription Options`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: newsletter/subscription/allow_guest_subscribe
							ID:        "allow_guest_subscribe",
							Label:     `Allow Guest Subscription`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: newsletter/subscription/confirm
							ID:        "confirm",
							Label:     `Need to Confirm`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: newsletter/subscription/confirm_email_identity
							ID:        "confirm_email_identity",
							Label:     `Confirmation Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `support`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: newsletter/subscription/confirm_email_template
							ID:        "confirm_email_template",
							Label:     `Confirmation Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `newsletter_subscription_confirm_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: newsletter/subscription/success_email_identity
							ID:        "success_email_identity",
							Label:     `Success Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: newsletter/subscription/success_email_template
							ID:        "success_email_template",
							Label:     `Success Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `newsletter_subscription_success_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: newsletter/subscription/un_email_identity
							ID:        "un_email_identity",
							Label:     `Unsubscription Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `support`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: newsletter/subscription/un_email_template
							ID:        "un_email_template",
							Label:     `Unsubscription Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `newsletter_subscription_un_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "newsletter",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "sending",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: newsletter/sending/set_return_path
							ID:      `set_return_path`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
