// +build ignore

package cookie

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
	// WebCookieCookieLifetime => Cookie Lifetime.
	// Path: web/cookie/cookie_lifetime
	// BackendModel: Magento\Cookie\Model\Config\Backend\Lifetime
	WebCookieCookieLifetime cfgmodel.Str

	// WebCookieCookiePath => Cookie Path.
	// Path: web/cookie/cookie_path
	// BackendModel: Magento\Cookie\Model\Config\Backend\Path
	WebCookieCookiePath cfgmodel.Str

	// WebCookieCookieDomain => Cookie Domain.
	// Path: web/cookie/cookie_domain
	// BackendModel: Magento\Cookie\Model\Config\Backend\Domain
	WebCookieCookieDomain cfgmodel.Str

	// WebCookieCookieHttponly => Use HTTP Only.
	// Warning: Do not set to "No". User security could be compromised.
	// Path: web/cookie/cookie_httponly
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebCookieCookieHttponly cfgmodel.Bool

	// WebCookieCookieRestriction => Cookie Restriction Mode.
	// Path: web/cookie/cookie_restriction
	// BackendModel: Magento\Cookie\Model\Config\Backend\Cookie
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebCookieCookieRestriction cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.WebCookieCookieLifetime = cfgmodel.NewStr(`web/cookie/cookie_lifetime`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebCookieCookiePath = cfgmodel.NewStr(`web/cookie/cookie_path`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebCookieCookieDomain = cfgmodel.NewStr(`web/cookie/cookie_domain`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebCookieCookieHttponly = cfgmodel.NewBool(`web/cookie/cookie_httponly`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebCookieCookieRestriction = cfgmodel.NewBool(`web/cookie/cookie_restriction`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
