// +build ignore

package user

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "admin",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "emails",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/emails/reset_password_template`,
						ID:           "reset_password_template",
						Label:        `Reset Password Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      `admin_emails_reset_password_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "admin",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "emails",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/emails/forgot_email_template`,
						ID:      "forgot_email_template",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `admin_emails_forgot_email_template`,
					},

					&config.Field{
						// Path: `admin/emails/forgot_email_identity`,
						ID:      "forgot_email_identity",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `general`,
					},

					&config.Field{
						// Path: `admin/emails/password_reset_link_expiration_period`,
						ID:      "password_reset_link_expiration_period",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: true,
					},
				},
			},
		},
	},
)
