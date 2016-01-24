// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "admin",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "emails",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/emails/reset_password_template
							ID:        "reset_password_template",
							Label:     `Reset Password Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `admin_emails_reset_password_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},
					),
				},

				&element.Group{
					ID: "security",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/security/lockout_failures
							ID:        "lockout_failures",
							Label:     `Maximum Login Failures to Lockout Account`,
							Comment:   text.Long(`We will disable this feature if the value is empty.`),
							Type:      element.Type,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   6,
						},

						&element.Field{
							// Path: admin/security/lockout_threshold
							ID:        "lockout_threshold",
							Label:     `Lockout Time (minutes)`,
							Type:      element.Type,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   30,
						},

						&element.Field{
							// Path: admin/security/password_lifetime
							ID:        "password_lifetime",
							Label:     `Password Lifetime (days)`,
							Comment:   text.Long(`We will disable this feature if the value is empty.`),
							Type:      element.Type,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   90,
						},

						&element.Field{
							// Path: admin/security/password_is_forced
							ID:        "password_is_forced",
							Label:     `Password Change`,
							Type:      element.TypeSelect,
							SortOrder: 130,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   true,
							// SourceModel: Magento\User\Model\System\Config\Source\Password
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "admin",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "emails",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/emails/forgot_email_template
							ID:      `forgot_email_template`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `admin_emails_forgot_email_template`,
						},

						&element.Field{
							// Path: admin/emails/forgot_email_identity
							ID:      `forgot_email_identity`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `general`,
						},

						&element.Field{
							// Path: admin/emails/password_reset_link_expiration_period
							ID:      `password_reset_link_expiration_period`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
