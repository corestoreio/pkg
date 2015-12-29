// +build ignore

package newsletter

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "newsletter",
		Label:     "Newsletter",
		SortOrder: 110,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "subscription",
				Label:     `Subscription Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `newsletter/subscription/allow_guest_subscribe`,
						ID:           "allow_guest_subscribe",
						Label:        `Allow Guest Subscription`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `newsletter/subscription/confirm`,
						ID:           "confirm",
						Label:        `Need to Confirm`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `newsletter/subscription/confirm_email_identity`,
						ID:           "confirm_email_identity",
						Label:        `Confirmation Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `support`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `newsletter/subscription/confirm_email_template`,
						ID:           "confirm_email_template",
						Label:        `Confirmation Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `newsletter_subscription_confirm_email_template`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `newsletter/subscription/success_email_identity`,
						ID:           "success_email_identity",
						Label:        `Success Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `general`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `newsletter/subscription/success_email_template`,
						ID:           "success_email_template",
						Label:        `Success Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `newsletter_subscription_success_email_template`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `newsletter/subscription/un_email_identity`,
						ID:           "un_email_identity",
						Label:        `Unsubscription Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `support`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `newsletter/subscription/un_email_template`,
						ID:           "un_email_template",
						Label:        `Unsubscription Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `newsletter_subscription_un_email_template`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "newsletter",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "sending",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `newsletter/sending/set_return_path`,
						ID:      "set_return_path",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},
				},
			},
		},
	},
)
