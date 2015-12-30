// +build ignore

package cookie

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathWebCookieCookieLifetime => Cookie Lifetime.
// BackendModel: Otnegam\Cookie\Model\Config\Backend\Lifetime
var PathWebCookieCookieLifetime = model.NewStr(`web/cookie/cookie_lifetime`)

// PathWebCookieCookiePath => Cookie Path.
// BackendModel: Otnegam\Cookie\Model\Config\Backend\Path
var PathWebCookieCookiePath = model.NewStr(`web/cookie/cookie_path`)

// PathWebCookieCookieDomain => Cookie Domain.
// BackendModel: Otnegam\Cookie\Model\Config\Backend\Domain
var PathWebCookieCookieDomain = model.NewStr(`web/cookie/cookie_domain`)

// PathWebCookieCookieHttponly => Use HTTP Only.
// Warning: Do not set to "No". User security could be compromised.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebCookieCookieHttponly = model.NewBool(`web/cookie/cookie_httponly`)

// PathWebCookieCookieRestriction => Cookie Restriction Mode.
// BackendModel: Otnegam\Cookie\Model\Config\Backend\Cookie
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebCookieCookieRestriction = model.NewBool(`web/cookie/cookie_restriction`)
