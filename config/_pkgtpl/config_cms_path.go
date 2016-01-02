// +build ignore

package cms

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
	// WebDefaultCmsHomePage => CMS Home Page.
	// Path: web/default/cms_home_page
	// SourceModel: Otnegam\Cms\Model\Config\Source\Page
	WebDefaultCmsHomePage model.Str

	// WebDefaultCmsNoRoute => CMS No Route Page.
	// Path: web/default/cms_no_route
	// SourceModel: Otnegam\Cms\Model\Config\Source\Page
	WebDefaultCmsNoRoute model.Str

	// WebDefaultCmsNoCookies => CMS No Cookies Page.
	// Path: web/default/cms_no_cookies
	// SourceModel: Otnegam\Cms\Model\Config\Source\Page
	WebDefaultCmsNoCookies model.Str

	// WebDefaultShowCmsBreadcrumbs => Show Breadcrumbs for CMS Pages.
	// Path: web/default/show_cms_breadcrumbs
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebDefaultShowCmsBreadcrumbs model.Bool

	// WebBrowserCapabilitiesCookies => Redirect to CMS-page if Cookies are Disabled.
	// Path: web/browser_capabilities/cookies
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesCookies model.Bool

	// WebBrowserCapabilitiesJavascript => Show Notice if JavaScript is Disabled.
	// Path: web/browser_capabilities/javascript
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesJavascript model.Bool

	// WebBrowserCapabilitiesLocalStorage => Show Notice if Local Storage is Disabled.
	// Path: web/browser_capabilities/local_storage
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesLocalStorage model.Bool

	// CmsWysiwygEnabled => Enable WYSIWYG Editor.
	// Path: cms/wysiwyg/enabled
	// SourceModel: Otnegam\Cms\Model\Config\Source\Wysiwyg\Enabled
	CmsWysiwygEnabled model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.WebDefaultCmsHomePage = model.NewStr(`web/default/cms_home_page`, model.WithPkgCfg(pkgCfg))
	pp.WebDefaultCmsNoRoute = model.NewStr(`web/default/cms_no_route`, model.WithPkgCfg(pkgCfg))
	pp.WebDefaultCmsNoCookies = model.NewStr(`web/default/cms_no_cookies`, model.WithPkgCfg(pkgCfg))
	pp.WebDefaultShowCmsBreadcrumbs = model.NewBool(`web/default/show_cms_breadcrumbs`, model.WithPkgCfg(pkgCfg))
	pp.WebBrowserCapabilitiesCookies = model.NewBool(`web/browser_capabilities/cookies`, model.WithPkgCfg(pkgCfg))
	pp.WebBrowserCapabilitiesJavascript = model.NewBool(`web/browser_capabilities/javascript`, model.WithPkgCfg(pkgCfg))
	pp.WebBrowserCapabilitiesLocalStorage = model.NewBool(`web/browser_capabilities/local_storage`, model.WithPkgCfg(pkgCfg))
	pp.CmsWysiwygEnabled = model.NewStr(`cms/wysiwyg/enabled`, model.WithPkgCfg(pkgCfg))

	return pp
}
