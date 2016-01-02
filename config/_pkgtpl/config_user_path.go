// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// AdminEmailsResetPasswordTemplate => Reset Password Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: admin/emails/reset_password_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	AdminEmailsResetPasswordTemplate model.Str

	// AdminSecurityLockoutFailures => Maximum Login Failures to Lockout Account.
	// We will disable this feature if the value is empty.
	// Path: admin/security/lockout_failures
	AdminSecurityLockoutFailures model.Str

	// AdminSecurityLockoutThreshold => Lockout Time (minutes).
	// Path: admin/security/lockout_threshold
	AdminSecurityLockoutThreshold model.Str

	// AdminSecurityPasswordLifetime => Password Lifetime (days).
	// We will disable this feature if the value is empty.
	// Path: admin/security/password_lifetime
	AdminSecurityPasswordLifetime model.Str

	// AdminSecurityPasswordIsForced => Password Change.
	// Path: admin/security/password_is_forced
	// SourceModel: Otnegam\User\Model\System\Config\Source\Password
	AdminSecurityPasswordIsForced model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminEmailsResetPasswordTemplate = model.NewStr(`admin/emails/reset_password_template`, model.WithPkgCfg(pkgCfg))
	pp.AdminSecurityLockoutFailures = model.NewStr(`admin/security/lockout_failures`, model.WithPkgCfg(pkgCfg))
	pp.AdminSecurityLockoutThreshold = model.NewStr(`admin/security/lockout_threshold`, model.WithPkgCfg(pkgCfg))
	pp.AdminSecurityPasswordLifetime = model.NewStr(`admin/security/password_lifetime`, model.WithPkgCfg(pkgCfg))
	pp.AdminSecurityPasswordIsForced = model.NewStr(`admin/security/password_is_forced`, model.WithPkgCfg(pkgCfg))

	return pp
}
