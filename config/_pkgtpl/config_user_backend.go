// +build ignore

package user

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
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

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminEmailsResetPasswordTemplate = model.NewStr(`admin/emails/reset_password_template`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityLockoutFailures = model.NewStr(`admin/security/lockout_failures`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityLockoutThreshold = model.NewStr(`admin/security/lockout_threshold`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityPasswordLifetime = model.NewStr(`admin/security/password_lifetime`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityPasswordIsForced = model.NewStr(`admin/security/password_is_forced`, model.WithConfigStructure(cfgStruct))

	return pp
}
