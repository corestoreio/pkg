// +build ignore

package cms

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
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
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.WebDefaultCmsHomePage = model.NewStr(`web/default/cms_home_page`, model.WithConfigStructure(cfgStruct))
	pp.WebDefaultCmsNoRoute = model.NewStr(`web/default/cms_no_route`, model.WithConfigStructure(cfgStruct))
	pp.WebDefaultCmsNoCookies = model.NewStr(`web/default/cms_no_cookies`, model.WithConfigStructure(cfgStruct))
	pp.WebDefaultShowCmsBreadcrumbs = model.NewBool(`web/default/show_cms_breadcrumbs`, model.WithConfigStructure(cfgStruct))
	pp.WebBrowserCapabilitiesCookies = model.NewBool(`web/browser_capabilities/cookies`, model.WithConfigStructure(cfgStruct))
	pp.WebBrowserCapabilitiesJavascript = model.NewBool(`web/browser_capabilities/javascript`, model.WithConfigStructure(cfgStruct))
	pp.WebBrowserCapabilitiesLocalStorage = model.NewBool(`web/browser_capabilities/local_storage`, model.WithConfigStructure(cfgStruct))
	pp.CmsWysiwygEnabled = model.NewStr(`cms/wysiwyg/enabled`, model.WithConfigStructure(cfgStruct))

	return pp
}
