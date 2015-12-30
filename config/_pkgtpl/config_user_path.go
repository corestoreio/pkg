// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathAdminEmailsResetPasswordTemplate => Reset Password Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathAdminEmailsResetPasswordTemplate = model.NewStr(`admin/emails/reset_password_template`)

// PathAdminSecurityLockoutFailures => Maximum Login Failures to Lockout Account.
// We will disable this feature if the value is empty.
var PathAdminSecurityLockoutFailures = model.NewStr(`admin/security/lockout_failures`)

// PathAdminSecurityLockoutThreshold => Lockout Time (minutes).
var PathAdminSecurityLockoutThreshold = model.NewStr(`admin/security/lockout_threshold`)

// PathAdminSecurityPasswordLifetime => Password Lifetime (days).
// We will disable this feature if the value is empty.
var PathAdminSecurityPasswordLifetime = model.NewStr(`admin/security/password_lifetime`)

// PathAdminSecurityPasswordIsForced => Password Change.
// SourceModel: Otnegam\User\Model\System\Config\Source\Password
var PathAdminSecurityPasswordIsForced = model.NewStr(`admin/security/password_is_forced`)
