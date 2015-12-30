// +build ignore

package captcha

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathAdminCaptchaEnable => Enable CAPTCHA in Admin.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminCaptchaEnable = model.NewBool(`admin/captcha/enable`)

// PathAdminCaptchaFont => Font.
// SourceModel: Otnegam\Captcha\Model\Config\Font
var PathAdminCaptchaFont = model.NewStr(`admin/captcha/font`)

// PathAdminCaptchaForms => Forms.
// SourceModel: Otnegam\Captcha\Model\Config\Form\Backend
var PathAdminCaptchaForms = model.NewStringCSV(`admin/captcha/forms`)

// PathAdminCaptchaMode => Displaying Mode.
// SourceModel: Otnegam\Captcha\Model\Config\Mode
var PathAdminCaptchaMode = model.NewStr(`admin/captcha/mode`)

// PathAdminCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
// If 0 is specified, CAPTCHA on the Login form will be always available.
var PathAdminCaptchaFailedAttemptsLogin = model.NewStr(`admin/captcha/failed_attempts_login`)

// PathAdminCaptchaTimeout => CAPTCHA Timeout (minutes).
var PathAdminCaptchaTimeout = model.NewStr(`admin/captcha/timeout`)

// PathAdminCaptchaLength => Number of Symbols.
// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
var PathAdminCaptchaLength = model.NewStr(`admin/captcha/length`)

// PathAdminCaptchaSymbols => Symbols Used in CAPTCHA.
// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
// spaces or other characters are allowed.Similar looking characters (e.g.
// "i", "l", "1") decrease chance of correct recognition by customer.
var PathAdminCaptchaSymbols = model.NewStr(`admin/captcha/symbols`)

// PathAdminCaptchaCaseSensitive => Case Sensitive.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminCaptchaCaseSensitive = model.NewBool(`admin/captcha/case_sensitive`)

// PathCustomerCaptchaEnable => Enable CAPTCHA on Storefront.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCaptchaEnable = model.NewBool(`customer/captcha/enable`)

// PathCustomerCaptchaFont => Font.
// SourceModel: Otnegam\Captcha\Model\Config\Font
var PathCustomerCaptchaFont = model.NewStr(`customer/captcha/font`)

// PathCustomerCaptchaForms => Forms.
// CAPTCHA for "Create user" and "Forgot password" forms is always enabled if
// chosen.
// SourceModel: Otnegam\Captcha\Model\Config\Form\Frontend
var PathCustomerCaptchaForms = model.NewStringCSV(`customer/captcha/forms`)

// PathCustomerCaptchaMode => Displaying Mode.
// SourceModel: Otnegam\Captcha\Model\Config\Mode
var PathCustomerCaptchaMode = model.NewStr(`customer/captcha/mode`)

// PathCustomerCaptchaFailedAttemptsLogin => Number of Unsuccessful Attempts to Login.
// If 0 is specified, CAPTCHA on the Login form will be always available.
var PathCustomerCaptchaFailedAttemptsLogin = model.NewStr(`customer/captcha/failed_attempts_login`)

// PathCustomerCaptchaTimeout => CAPTCHA Timeout (minutes).
var PathCustomerCaptchaTimeout = model.NewStr(`customer/captcha/timeout`)

// PathCustomerCaptchaLength => Number of Symbols.
// Please specify 8 symbols at the most. Range allowed (e.g. 3-5)
var PathCustomerCaptchaLength = model.NewStr(`customer/captcha/length`)

// PathCustomerCaptchaSymbols => Symbols Used in CAPTCHA.
// Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No
// spaces or other characters are allowed.Similar looking characters (e.g.
// "i", "l", "1") decrease chance of correct recognition by customer.
var PathCustomerCaptchaSymbols = model.NewStr(`customer/captcha/symbols`)

// PathCustomerCaptchaCaseSensitive => Case Sensitive.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCaptchaCaseSensitive = model.NewBool(`customer/captcha/case_sensitive`)
