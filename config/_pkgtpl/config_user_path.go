// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
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
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminEmailsResetPasswordTemplate = model.NewStr(`admin/emails/reset_password_template`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityLockoutFailures = model.NewStr(`admin/security/lockout_failures`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityLockoutThreshold = model.NewStr(`admin/security/lockout_threshold`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityPasswordLifetime = model.NewStr(`admin/security/password_lifetime`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityPasswordIsForced = model.NewStr(`admin/security/password_is_forced`, model.WithConfigStructure(cfgStruct))

	return pp
}
