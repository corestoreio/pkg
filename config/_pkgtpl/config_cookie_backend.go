// +build ignore

package cookie

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
	// WebCookieCookieLifetime => Cookie Lifetime.
	// Path: web/cookie/cookie_lifetime
	// BackendModel: Otnegam\Cookie\Model\Config\Backend\Lifetime
	WebCookieCookieLifetime model.Str

	// WebCookieCookiePath => Cookie Path.
	// Path: web/cookie/cookie_path
	// BackendModel: Otnegam\Cookie\Model\Config\Backend\Path
	WebCookieCookiePath model.Str

	// WebCookieCookieDomain => Cookie Domain.
	// Path: web/cookie/cookie_domain
	// BackendModel: Otnegam\Cookie\Model\Config\Backend\Domain
	WebCookieCookieDomain model.Str

	// WebCookieCookieHttponly => Use HTTP Only.
	// Warning: Do not set to "No". User security could be compromised.
	// Path: web/cookie/cookie_httponly
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebCookieCookieHttponly model.Bool

	// WebCookieCookieRestriction => Cookie Restriction Mode.
	// Path: web/cookie/cookie_restriction
	// BackendModel: Otnegam\Cookie\Model\Config\Backend\Cookie
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebCookieCookieRestriction model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.WebCookieCookieLifetime = model.NewStr(`web/cookie/cookie_lifetime`, model.WithConfigStructure(cfgStruct))
	pp.WebCookieCookiePath = model.NewStr(`web/cookie/cookie_path`, model.WithConfigStructure(cfgStruct))
	pp.WebCookieCookieDomain = model.NewStr(`web/cookie/cookie_domain`, model.WithConfigStructure(cfgStruct))
	pp.WebCookieCookieHttponly = model.NewBool(`web/cookie/cookie_httponly`, model.WithConfigStructure(cfgStruct))
	pp.WebCookieCookieRestriction = model.NewBool(`web/cookie/cookie_restriction`, model.WithConfigStructure(cfgStruct))

	return pp
}
