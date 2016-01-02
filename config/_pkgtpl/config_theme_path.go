// +build ignore

package theme

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
	// DesignHeadShortcutIcon => Favicon Icon.
	// Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers
	// support all these formats!
	// Path: design/head/shortcut_icon
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Favicon
	DesignHeadShortcutIcon model.Str

	// DesignHeadDefaultTitle => Default Title.
	// Path: design/head/default_title
	DesignHeadDefaultTitle model.Str

	// DesignHeadTitlePrefix => Title Prefix.
	// Path: design/head/title_prefix
	DesignHeadTitlePrefix model.Str

	// DesignHeadTitleSuffix => Title Suffix.
	// Path: design/head/title_suffix
	DesignHeadTitleSuffix model.Str

	// DesignHeadDefaultDescription => Default Description.
	// Path: design/head/default_description
	DesignHeadDefaultDescription model.Str

	// DesignHeadDefaultKeywords => Default Keywords.
	// Path: design/head/default_keywords
	DesignHeadDefaultKeywords model.Str

	// DesignHeadIncludes => Miscellaneous Scripts.
	// This will be included before head closing tag in page HTML.
	// Path: design/head/includes
	DesignHeadIncludes model.Str

	// DesignHeadDemonotice => Display Demo Store Notice.
	// Path: design/head/demonotice
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DesignHeadDemonotice model.Bool

	// DesignSearchEngineRobotsDefaultRobots => Default Robots.
	// This will be included before head closing tag in page HTML.
	// Path: design/search_engine_robots/default_robots
	// SourceModel: Otnegam\Config\Model\Config\Source\Design\Robots
	DesignSearchEngineRobotsDefaultRobots model.Str

	// DesignSearchEngineRobotsCustomInstructions => Edit custom instruction of robots.txt File.
	// Path: design/search_engine_robots/custom_instructions
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Robots
	DesignSearchEngineRobotsCustomInstructions model.Str

	// DesignSearchEngineRobotsResetToDefaults => Reset to Defaults.
	// This action will delete your custom instructions and reset robots.txt file
	// to system's default settings.
	// Path: design/search_engine_robots/reset_to_defaults
	DesignSearchEngineRobotsResetToDefaults model.Str

	// DesignHeaderLogoSrc => Logo Image.
	// Allowed file types:PNG, GIF, JPG, JPEG, SVG.
	// Path: design/header/logo_src
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Logo
	DesignHeaderLogoSrc model.Str

	// DesignHeaderLogoWidth => Logo Image Width.
	// Path: design/header/logo_width
	DesignHeaderLogoWidth model.Str

	// DesignHeaderLogoHeight => Logo Image Height.
	// Path: design/header/logo_height
	DesignHeaderLogoHeight model.Str

	// DesignHeaderLogoAlt => Logo Image Alt.
	// Path: design/header/logo_alt
	DesignHeaderLogoAlt model.Str

	// DesignHeaderWelcome => Welcome Text.
	// Path: design/header/welcome
	DesignHeaderWelcome model.Str

	// DesignFooterCopyright => Copyright.
	// Path: design/footer/copyright
	DesignFooterCopyright model.Str

	// DesignFooterAbsoluteFooter => Miscellaneous HTML.
	// This will be displayed just before body closing tag.
	// Path: design/footer/absolute_footer
	DesignFooterAbsoluteFooter model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignHeadShortcutIcon = model.NewStr(`design/head/shortcut_icon`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadDefaultTitle = model.NewStr(`design/head/default_title`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadTitlePrefix = model.NewStr(`design/head/title_prefix`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadTitleSuffix = model.NewStr(`design/head/title_suffix`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadDefaultDescription = model.NewStr(`design/head/default_description`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadDefaultKeywords = model.NewStr(`design/head/default_keywords`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadIncludes = model.NewStr(`design/head/includes`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeadDemonotice = model.NewBool(`design/head/demonotice`, model.WithPkgCfg(pkgCfg))
	pp.DesignSearchEngineRobotsDefaultRobots = model.NewStr(`design/search_engine_robots/default_robots`, model.WithPkgCfg(pkgCfg))
	pp.DesignSearchEngineRobotsCustomInstructions = model.NewStr(`design/search_engine_robots/custom_instructions`, model.WithPkgCfg(pkgCfg))
	pp.DesignSearchEngineRobotsResetToDefaults = model.NewStr(`design/search_engine_robots/reset_to_defaults`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeaderLogoSrc = model.NewStr(`design/header/logo_src`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeaderLogoWidth = model.NewStr(`design/header/logo_width`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeaderLogoHeight = model.NewStr(`design/header/logo_height`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeaderLogoAlt = model.NewStr(`design/header/logo_alt`, model.WithPkgCfg(pkgCfg))
	pp.DesignHeaderWelcome = model.NewStr(`design/header/welcome`, model.WithPkgCfg(pkgCfg))
	pp.DesignFooterCopyright = model.NewStr(`design/footer/copyright`, model.WithPkgCfg(pkgCfg))
	pp.DesignFooterAbsoluteFooter = model.NewStr(`design/footer/absolute_footer`, model.WithPkgCfg(pkgCfg))

	return pp
}
