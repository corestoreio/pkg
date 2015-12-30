// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "admin",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "emails",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/emails/reset_password_template
						ID:        "reset_password_template",
						Label:     `Reset Password Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `admin_emails_reset_password_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},
				),
			},

			&config.Group{
				ID: "security",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/security/lockout_failures
						ID:        "lockout_failures",
						Label:     `Maximum Login Failures to Lockout Account`,
						Comment:   element.LongText(`We will disable this feature if the value is empty.`),
						Type:      config.Type,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   6,
					},

					&config.Field{
						// Path: admin/security/lockout_threshold
						ID:        "lockout_threshold",
						Label:     `Lockout Time (minutes)`,
						Type:      config.Type,
						SortOrder: 110,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   30,
					},

					&config.Field{
						// Path: admin/security/password_lifetime
						ID:        "password_lifetime",
						Label:     `Password Lifetime (days)`,
						Comment:   element.LongText(`We will disable this feature if the value is empty.`),
						Type:      config.Type,
						SortOrder: 120,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   90,
					},

					&config.Field{
						// Path: admin/security/password_is_forced
						ID:        "password_is_forced",
						Label:     `Password Change`,
						Type:      config.TypeSelect,
						SortOrder: 130,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\User\Model\System\Config\Source\Password
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "admin",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "emails",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/emails/forgot_email_template
						ID:      `forgot_email_template`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `admin_emails_forgot_email_template`,
					},

					&config.Field{
						// Path: admin/emails/forgot_email_identity
						ID:      `forgot_email_identity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `general`,
					},

					&config.Field{
						// Path: admin/emails/password_reset_link_expiration_period
						ID:      `password_reset_link_expiration_period`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},
		),
	},
)
