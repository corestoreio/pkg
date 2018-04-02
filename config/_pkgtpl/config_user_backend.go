// +build ignore

package user

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// AdminEmailsResetPasswordTemplate => Reset Password Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: admin/emails/reset_password_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	AdminEmailsResetPasswordTemplate cfgmodel.Str

	// AdminSecurityLockoutFailures => Maximum Login Failures to Lockout Account.
	// We will disable this feature if the value is empty.
	// Path: admin/security/lockout_failures
	AdminSecurityLockoutFailures cfgmodel.Str

	// AdminSecurityLockoutThreshold => Lockout Time (minutes).
	// Path: admin/security/lockout_threshold
	AdminSecurityLockoutThreshold cfgmodel.Str

	// AdminSecurityPasswordLifetime => Password Lifetime (days).
	// We will disable this feature if the value is empty.
	// Path: admin/security/password_lifetime
	AdminSecurityPasswordLifetime cfgmodel.Str

	// AdminSecurityPasswordIsForced => Password Change.
	// Path: admin/security/password_is_forced
	// SourceModel: Magento\User\Model\System\Config\Source\Password
	AdminSecurityPasswordIsForced cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminEmailsResetPasswordTemplate = cfgmodel.NewStr(`admin/emails/reset_password_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityLockoutFailures = cfgmodel.NewStr(`admin/security/lockout_failures`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityLockoutThreshold = cfgmodel.NewStr(`admin/security/lockout_threshold`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityPasswordLifetime = cfgmodel.NewStr(`admin/security/password_lifetime`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityPasswordIsForced = cfgmodel.NewStr(`admin/security/password_is_forced`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
