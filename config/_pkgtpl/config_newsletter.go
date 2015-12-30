// +build ignore

package newsletter

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID:        "newsletter",
		Label:     `Newsletter`,
		SortOrder: 110,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Newsletter::newsletter
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "subscription",
				Label:     `Subscription Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: newsletter/subscription/allow_guest_subscribe
						ID:        "allow_guest_subscribe",
						Label:     `Allow Guest Subscription`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: newsletter/subscription/confirm
						ID:        "confirm",
						Label:     `Need to Confirm`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: newsletter/subscription/confirm_email_identity
						ID:        "confirm_email_identity",
						Label:     `Confirmation Email Sender`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `support`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&element.Field{
						// Path: newsletter/subscription/confirm_email_template
						ID:        "confirm_email_template",
						Label:     `Confirmation Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `newsletter_subscription_confirm_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&element.Field{
						// Path: newsletter/subscription/success_email_identity
						ID:        "success_email_identity",
						Label:     `Success Email Sender`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&element.Field{
						// Path: newsletter/subscription/success_email_template
						ID:        "success_email_template",
						Label:     `Success Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `newsletter_subscription_success_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&element.Field{
						// Path: newsletter/subscription/un_email_identity
						ID:        "un_email_identity",
						Label:     `Unsubscription Email Sender`,
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `support`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&element.Field{
						// Path: newsletter/subscription/un_email_template
						ID:        "un_email_template",
						Label:     `Unsubscription Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `newsletter_subscription_un_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&element.Section{
		ID: "newsletter",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID: "sending",
				Fields: element.NewFieldSlice(
					&element.Field{
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
