// +build ignore

package captcha

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
	// AdminCaptchaEnable => Enable CAPTCHA in Admin.
	// Path: admin/captcha/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminCaptchaEnable model.Bool

	// AdminCaptchaFont => Font.
	// Path: admin/captcha/font
	// SourceModel: Magento\Captcha\Model\Config\Font
	AdminCaptchaFont model.Str

	// AdminCaptchaForms => Forms.
	// Path: admin/captcha/forms
	// SourceModel: Magento\Captcha\Model\Config\Form\Backend
	AdminCaptchaForms model.StringCSV

	// AdminCaptchaMode => Displaying Mode.
	// Path: admin/captcha/mode
	// SourceModel: Magento\Captcha\Model\Config\Mode
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
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminCaptchaCaseSensitive model.Bool

	// CustomerCaptchaEnable => Enable CAPTCHA on Storefront.
	// Path: customer/captcha/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCaptchaEnable model.Bool

	// CustomerCaptchaFont => Font.
	// Path: customer/captcha/font
	// SourceModel: Magento\Captcha\Model\Config\Font
	CustomerCaptchaFont model.Str

	// CustomerCaptchaForms => Forms.
	// CAPTCHA for "Create user" and "Forgot password" forms is always enabled if
	// chosen.
	// Path: customer/captcha/forms
	// SourceModel: Magento\Captcha\Model\Config\Form\Frontend
	CustomerCaptchaForms model.StringCSV

	// CustomerCaptchaMode => Displaying Mode.
	// Path: customer/captcha/mode
	// SourceModel: Magento\Captcha\Model\Config\Mode
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
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCaptchaCaseSensitive model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminCaptchaEnable = model.NewBool(`admin/captcha/enable`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaFont = model.NewStr(`admin/captcha/font`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaForms = model.NewStringCSV(`admin/captcha/forms`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaMode = model.NewStr(`admin/captcha/mode`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaFailedAttemptsLogin = model.NewStr(`admin/captcha/failed_attempts_login`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaTimeout = model.NewStr(`admin/captcha/timeout`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaLength = model.NewStr(`admin/captcha/length`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaSymbols = model.NewStr(`admin/captcha/symbols`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaCaseSensitive = model.NewBool(`admin/captcha/case_sensitive`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaEnable = model.NewBool(`customer/captcha/enable`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaFont = model.NewStr(`customer/captcha/font`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaForms = model.NewStringCSV(`customer/captcha/forms`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaMode = model.NewStr(`customer/captcha/mode`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaFailedAttemptsLogin = model.NewStr(`customer/captcha/failed_attempts_login`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaTimeout = model.NewStr(`customer/captcha/timeout`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaLength = model.NewStr(`customer/captcha/length`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaSymbols = model.NewStr(`customer/captcha/symbols`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaCaseSensitive = model.NewBool(`customer/captcha/case_sensitive`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
