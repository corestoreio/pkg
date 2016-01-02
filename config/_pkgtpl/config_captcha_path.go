// +build ignore

package captcha

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
	// AdminCaptchaEnable => Enable CAPTCHA in Admin.
	// Path: admin/captcha/enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminCaptchaEnable model.Bool

	// AdminCaptchaFont => Font.
	// Path: admin/captcha/font
	// SourceModel: Otnegam\Captcha\Model\Config\Font
	AdminCaptchaFont model.Str

	// AdminCaptchaForms => Forms.
	// Path: admin/captcha/forms
	// SourceModel: Otnegam\Captcha\Model\Config\Form\Backend
	AdminCaptchaForms model.StringCSV

	// AdminCaptchaMode => Displaying Mode.
	// Path: admin/captcha/mode
	// SourceModel: Otnegam\Captcha\Model\Config\Mode
	AdminCaptchaMode model.Str

	// AdminCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
	// If 0 is specified, CAPTCHA on the Login form will be always available.
	// Path: admin/captcha/failed_attempts_login
	AdminCaptchaFailedAttemptsLogin model.Str

	// AdminCaptchaTimeout => CAPTCHA Timeout (minutes).
	// Path: admin/captcha/timeout
	AdminCaptchaTimeout model.Str

	// AdminCaptchaLength => Number of Symbols.
	// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
	// Path: admin/captcha/length
	AdminCaptchaLength model.Str

	// AdminCaptchaSymbols => Symbols Used in CAPTCHA.
	// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
	// spaces or other characters are allowed.Similar looking characters (e.g.
	// "i", "l", "1") decrease chance of correct recognition by customer.
	// Path: admin/captcha/symbols
	AdminCaptchaSymbols model.Str

	// AdminCaptchaCaseSensitive => Case Sensitive.
	// Path: admin/captcha/case_sensitive
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminCaptchaCaseSensitive model.Bool

	// CustomerCaptchaEnable => Enable CAPTCHA on Storefront.
	// Path: customer/captcha/enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CustomerCaptchaEnable model.Bool

	// CustomerCaptchaFont => Font.
	// Path: customer/captcha/font
	// SourceModel: Otnegam\Captcha\Model\Config\Font
	CustomerCaptchaFont model.Str

	// CustomerCaptchaForms => Forms.
	// CAPTCHA for "Create user" and "Forgot password" forms is always enabled if
	// chosen.
	// Path: customer/captcha/forms
	// SourceModel: Otnegam\Captcha\Model\Config\Form\Frontend
	CustomerCaptchaForms model.StringCSV

	// CustomerCaptchaMode => Displaying Mode.
	// Path: customer/captcha/mode
	// SourceModel: Otnegam\Captcha\Model\Config\Mode
	CustomerCaptchaMode model.Str

	// CustomerCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
	// If 0 is specified, CAPTCHA on the Login form will be always available.
	// Path: customer/captcha/failed_attempts_login
	CustomerCaptchaFailedAttemptsLogin model.Str

	// CustomerCaptchaTimeout => CAPTCHA Timeout (minutes).
	// Path: customer/captcha/timeout
	CustomerCaptchaTimeout model.Str

	// CustomerCaptchaLength => Number of Symbols.
	// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
	// Path: customer/captcha/length
	CustomerCaptchaLength model.Str

	// CustomerCaptchaSymbols => Symbols Used in CAPTCHA.
	// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
	// spaces or other characters are allowed.Similar looking characters (e.g.
	// "i", "l", "1") decrease chance of correct recognition by customer.
	// Path: customer/captcha/symbols
	CustomerCaptchaSymbols model.Str

	// CustomerCaptchaCaseSensitive => Case Sensitive.
	// Path: customer/captcha/case_sensitive
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CustomerCaptchaCaseSensitive model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminCaptchaEnable = model.NewBool(`admin/captcha/enable`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaFont = model.NewStr(`admin/captcha/font`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaForms = model.NewStringCSV(`admin/captcha/forms`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaMode = model.NewStr(`admin/captcha/mode`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaFailedAttemptsLogin = model.NewStr(`admin/captcha/failed_attempts_login`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaTimeout = model.NewStr(`admin/captcha/timeout`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaLength = model.NewStr(`admin/captcha/length`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaSymbols = model.NewStr(`admin/captcha/symbols`, model.WithPkgCfg(pkgCfg))
	pp.AdminCaptchaCaseSensitive = model.NewBool(`admin/captcha/case_sensitive`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaEnable = model.NewBool(`customer/captcha/enable`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaFont = model.NewStr(`customer/captcha/font`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaForms = model.NewStringCSV(`customer/captcha/forms`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaMode = model.NewStr(`customer/captcha/mode`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaFailedAttemptsLogin = model.NewStr(`customer/captcha/failed_attempts_login`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaTimeout = model.NewStr(`customer/captcha/timeout`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaLength = model.NewStr(`customer/captcha/length`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaSymbols = model.NewStr(`customer/captcha/symbols`, model.WithPkgCfg(pkgCfg))
	pp.CustomerCaptchaCaseSensitive = model.NewBool(`customer/captcha/case_sensitive`, model.WithPkgCfg(pkgCfg))

	return pp
}
