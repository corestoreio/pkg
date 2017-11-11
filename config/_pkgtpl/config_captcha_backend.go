// +build ignore

package captcha

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// AdminCaptchaEnable => Enable CAPTCHA in Admin.
	// Path: admin/captcha/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminCaptchaEnable cfgmodel.Bool

	// AdminCaptchaFont => Font.
	// Path: admin/captcha/font
	// SourceModel: Magento\Captcha\Model\Config\Font
	AdminCaptchaFont cfgmodel.Str

	// AdminCaptchaForms => Forms.
	// Path: admin/captcha/forms
	// SourceModel: Magento\Captcha\Model\Config\Form\Backend
	AdminCaptchaForms cfgmodel.StringCSV

	// AdminCaptchaMode => Displaying Mode.
	// Path: admin/captcha/mode
	// SourceModel: Magento\Captcha\Model\Config\Mode
	AdminCaptchaMode cfgmodel.Str

	// AdminCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
	// If 0 is specified, CAPTCHA on the Login form will be always available.
	// Path: admin/captcha/failed_attempts_login
	AdminCaptchaFailedAttemptsLogin cfgmodel.Str

	// AdminCaptchaTimeout => CAPTCHA Timeout (minutes).
	// Path: admin/captcha/timeout
	AdminCaptchaTimeout cfgmodel.Str

	// AdminCaptchaLength => Number of Symbols.
	// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
	// Path: admin/captcha/length
	AdminCaptchaLength cfgmodel.Str

	// AdminCaptchaSymbols => Symbols Used in CAPTCHA.
	// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
	// spaces or other characters are allowed.Similar looking characters (e.g.
	// "i", "l", "1") decrease chance of correct recognition by customer.
	// Path: admin/captcha/symbols
	AdminCaptchaSymbols cfgmodel.Str

	// AdminCaptchaCaseSensitive => Case Sensitive.
	// Path: admin/captcha/case_sensitive
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminCaptchaCaseSensitive cfgmodel.Bool

	// CustomerCaptchaEnable => Enable CAPTCHA on Storefront.
	// Path: customer/captcha/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCaptchaEnable cfgmodel.Bool

	// CustomerCaptchaFont => Font.
	// Path: customer/captcha/font
	// SourceModel: Magento\Captcha\Model\Config\Font
	CustomerCaptchaFont cfgmodel.Str

	// CustomerCaptchaForms => Forms.
	// CAPTCHA for "Create user" and "Forgot password" forms is always enabled if
	// chosen.
	// Path: customer/captcha/forms
	// SourceModel: Magento\Captcha\Model\Config\Form\Frontend
	CustomerCaptchaForms cfgmodel.StringCSV

	// CustomerCaptchaMode => Displaying Mode.
	// Path: customer/captcha/mode
	// SourceModel: Magento\Captcha\Model\Config\Mode
	CustomerCaptchaMode cfgmodel.Str

	// CustomerCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
	// If 0 is specified, CAPTCHA on the Login form will be always available.
	// Path: customer/captcha/failed_attempts_login
	CustomerCaptchaFailedAttemptsLogin cfgmodel.Str

	// CustomerCaptchaTimeout => CAPTCHA Timeout (minutes).
	// Path: customer/captcha/timeout
	CustomerCaptchaTimeout cfgmodel.Str

	// CustomerCaptchaLength => Number of Symbols.
	// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
	// Path: customer/captcha/length
	CustomerCaptchaLength cfgmodel.Str

	// CustomerCaptchaSymbols => Symbols Used in CAPTCHA.
	// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
	// spaces or other characters are allowed.Similar looking characters (e.g.
	// "i", "l", "1") decrease chance of correct recognition by customer.
	// Path: customer/captcha/symbols
	CustomerCaptchaSymbols cfgmodel.Str

	// CustomerCaptchaCaseSensitive => Case Sensitive.
	// Path: customer/captcha/case_sensitive
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CustomerCaptchaCaseSensitive cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.AdminCaptchaEnable = cfgmodel.NewBool(`admin/captcha/enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaFont = cfgmodel.NewStr(`admin/captcha/font`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaForms = cfgmodel.NewStringCSV(`admin/captcha/forms`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaMode = cfgmodel.NewStr(`admin/captcha/mode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaFailedAttemptsLogin = cfgmodel.NewStr(`admin/captcha/failed_attempts_login`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaTimeout = cfgmodel.NewStr(`admin/captcha/timeout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaLength = cfgmodel.NewStr(`admin/captcha/length`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaSymbols = cfgmodel.NewStr(`admin/captcha/symbols`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminCaptchaCaseSensitive = cfgmodel.NewBool(`admin/captcha/case_sensitive`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaEnable = cfgmodel.NewBool(`customer/captcha/enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaFont = cfgmodel.NewStr(`customer/captcha/font`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaForms = cfgmodel.NewStringCSV(`customer/captcha/forms`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaMode = cfgmodel.NewStr(`customer/captcha/mode`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaFailedAttemptsLogin = cfgmodel.NewStr(`customer/captcha/failed_attempts_login`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaTimeout = cfgmodel.NewStr(`customer/captcha/timeout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaLength = cfgmodel.NewStr(`customer/captcha/length`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaSymbols = cfgmodel.NewStr(`customer/captcha/symbols`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CustomerCaptchaCaseSensitive = cfgmodel.NewBool(`customer/captcha/case_sensitive`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
