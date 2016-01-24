// +build ignore

package theme

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
	// DesignHeadShortcutIcon => Favicon Icon.
	// Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers
	// support all these formats!
	// Path: design/head/shortcut_icon
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Favicon
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
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DesignHeadDemonotice model.Bool

	// DesignSearchEngineRobotsDefaultRobots => Default Robots.
	// This will be included before head closing tag in page HTML.
	// Path: design/search_engine_robots/default_robots
	// SourceModel: Magento\Config\Model\Config\Source\Design\Robots
	DesignSearchEngineRobotsDefaultRobots model.Str

	// DesignSearchEngineRobotsCustomInstructions => Edit custom instruction of robots.txt File.
	// Path: design/search_engine_robots/custom_instructions
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Robots
	DesignSearchEngineRobotsCustomInstructions model.Str

	// DesignSearchEngineRobotsResetToDefaults => Reset to Defaults.
	// This action will delete your custom instructions and reset robots.txt file
	// to system's default settings.
	// Path: design/search_engine_robots/reset_to_defaults
	DesignSearchEngineRobotsResetToDefaults model.Str

	// DesignHeaderLogoSrc => Logo Image.
	// Allowed file types:PNG, GIF, JPG, JPEG, SVG.
	// Path: design/header/logo_src
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Logo
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

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignHeadShortcutIcon = model.NewStr(`design/head/shortcut_icon`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadDefaultTitle = model.NewStr(`design/head/default_title`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadTitlePrefix = model.NewStr(`design/head/title_prefix`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadTitleSuffix = model.NewStr(`design/head/title_suffix`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadDefaultDescription = model.NewStr(`design/head/default_description`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadDefaultKeywords = model.NewStr(`design/head/default_keywords`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadIncludes = model.NewStr(`design/head/includes`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeadDemonotice = model.NewBool(`design/head/demonotice`, model.WithConfigStructure(cfgStruct))
	pp.DesignSearchEngineRobotsDefaultRobots = model.NewStr(`design/search_engine_robots/default_robots`, model.WithConfigStructure(cfgStruct))
	pp.DesignSearchEngineRobotsCustomInstructions = model.NewStr(`design/search_engine_robots/custom_instructions`, model.WithConfigStructure(cfgStruct))
	pp.DesignSearchEngineRobotsResetToDefaults = model.NewStr(`design/search_engine_robots/reset_to_defaults`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeaderLogoSrc = model.NewStr(`design/header/logo_src`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeaderLogoWidth = model.NewStr(`design/header/logo_width`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeaderLogoHeight = model.NewStr(`design/header/logo_height`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeaderLogoAlt = model.NewStr(`design/header/logo_alt`, model.WithConfigStructure(cfgStruct))
	pp.DesignHeaderWelcome = model.NewStr(`design/header/welcome`, model.WithConfigStructure(cfgStruct))
	pp.DesignFooterCopyright = model.NewStr(`design/footer/copyright`, model.WithConfigStructure(cfgStruct))
	pp.DesignFooterAbsoluteFooter = model.NewStr(`design/footer/absolute_footer`, model.WithConfigStructure(cfgStruct))

	return pp
}
