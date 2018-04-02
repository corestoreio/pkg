// +build ignore

package cms

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
	// WebDefaultCmsHomePage => CMS Home Page.
	// Path: web/default/cms_home_page
	// SourceModel: Magento\Cms\Model\Config\Source\Page
	WebDefaultCmsHomePage cfgmodel.Str

	// WebDefaultCmsNoRoute => CMS No Route Page.
	// Path: web/default/cms_no_route
	// SourceModel: Magento\Cms\Model\Config\Source\Page
	WebDefaultCmsNoRoute cfgmodel.Str

	// WebDefaultCmsNoCookies => CMS No Cookies Page.
	// Path: web/default/cms_no_cookies
	// SourceModel: Magento\Cms\Model\Config\Source\Page
	WebDefaultCmsNoCookies cfgmodel.Str

	// WebDefaultShowCmsBreadcrumbs => Show Breadcrumbs for CMS Pages.
	// Path: web/default/show_cms_breadcrumbs
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebDefaultShowCmsBreadcrumbs cfgmodel.Bool

	// WebBrowserCapabilitiesCookies => Redirect to CMS-page if Cookies are Disabled.
	// Path: web/browser_capabilities/cookies
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesCookies cfgmodel.Bool

	// WebBrowserCapabilitiesJavascript => Show Notice if JavaScript is Disabled.
	// Path: web/browser_capabilities/javascript
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesJavascript cfgmodel.Bool

	// WebBrowserCapabilitiesLocalStorage => Show Notice if Local Storage is Disabled.
	// Path: web/browser_capabilities/local_storage
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebBrowserCapabilitiesLocalStorage cfgmodel.Bool

	// CmsWysiwygEnabled => Enable WYSIWYG Editor.
	// Path: cms/wysiwyg/enabled
	// SourceModel: Magento\Cms\Model\Config\Source\Wysiwyg\Enabled
	CmsWysiwygEnabled cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.WebDefaultCmsHomePage = cfgmodel.NewStr(`web/default/cms_home_page`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebDefaultCmsNoRoute = cfgmodel.NewStr(`web/default/cms_no_route`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebDefaultCmsNoCookies = cfgmodel.NewStr(`web/default/cms_no_cookies`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebDefaultShowCmsBreadcrumbs = cfgmodel.NewBool(`web/default/show_cms_breadcrumbs`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebBrowserCapabilitiesCookies = cfgmodel.NewBool(`web/browser_capabilities/cookies`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebBrowserCapabilitiesJavascript = cfgmodel.NewBool(`web/browser_capabilities/javascript`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebBrowserCapabilitiesLocalStorage = cfgmodel.NewBool(`web/browser_capabilities/local_storage`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CmsWysiwygEnabled = cfgmodel.NewStr(`cms/wysiwyg/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
