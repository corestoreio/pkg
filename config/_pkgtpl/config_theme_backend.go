// +build ignore

package theme

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
	// DesignHeadShortcutIcon => Favicon Icon.
	// Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers
	// support all these formats!
	// Path: design/head/shortcut_icon
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Favicon
	DesignHeadShortcutIcon cfgmodel.Str

	// DesignHeadDefaultTitle => Default Title.
	// Path: design/head/default_title
	DesignHeadDefaultTitle cfgmodel.Str

	// DesignHeadTitlePrefix => Title Prefix.
	// Path: design/head/title_prefix
	DesignHeadTitlePrefix cfgmodel.Str

	// DesignHeadTitleSuffix => Title Suffix.
	// Path: design/head/title_suffix
	DesignHeadTitleSuffix cfgmodel.Str

	// DesignHeadDefaultDescription => Default Description.
	// Path: design/head/default_description
	DesignHeadDefaultDescription cfgmodel.Str

	// DesignHeadDefaultKeywords => Default Keywords.
	// Path: design/head/default_keywords
	DesignHeadDefaultKeywords cfgmodel.Str

	// DesignHeadIncludes => Miscellaneous Scripts.
	// This will be included before head closing tag in page HTML.
	// Path: design/head/includes
	DesignHeadIncludes cfgmodel.Str

	// DesignHeadDemonotice => Display Demo Store Notice.
	// Path: design/head/demonotice
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DesignHeadDemonotice cfgmodel.Bool

	// DesignSearchEngineRobotsDefaultRobots => Default Robots.
	// This will be included before head closing tag in page HTML.
	// Path: design/search_engine_robots/default_robots
	// SourceModel: Magento\Config\Model\Config\Source\Design\Robots
	DesignSearchEngineRobotsDefaultRobots cfgmodel.Str

	// DesignSearchEngineRobotsCustomInstructions => Edit custom instruction of robots.txt File.
	// Path: design/search_engine_robots/custom_instructions
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Robots
	DesignSearchEngineRobotsCustomInstructions cfgmodel.Str

	// DesignSearchEngineRobotsResetToDefaults => Reset to Defaults.
	// This action will delete your custom instructions and reset robots.txt file
	// to system's default settings.
	// Path: design/search_engine_robots/reset_to_defaults
	DesignSearchEngineRobotsResetToDefaults cfgmodel.Str

	// DesignHeaderLogoSrc => Logo Image.
	// Allowed file types:PNG, GIF, JPG, JPEG, SVG.
	// Path: design/header/logo_src
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Logo
	DesignHeaderLogoSrc cfgmodel.Str

	// DesignHeaderLogoWidth => Logo Image Width.
	// Path: design/header/logo_width
	DesignHeaderLogoWidth cfgmodel.Str

	// DesignHeaderLogoHeight => Logo Image Height.
	// Path: design/header/logo_height
	DesignHeaderLogoHeight cfgmodel.Str

	// DesignHeaderLogoAlt => Logo Image Alt.
	// Path: design/header/logo_alt
	DesignHeaderLogoAlt cfgmodel.Str

	// DesignHeaderWelcome => Welcome Text.
	// Path: design/header/welcome
	DesignHeaderWelcome cfgmodel.Str

	// DesignFooterCopyright => Copyright.
	// Path: design/footer/copyright
	DesignFooterCopyright cfgmodel.Str

	// DesignFooterAbsoluteFooter => Miscellaneous HTML.
	// This will be displayed just before body closing tag.
	// Path: design/footer/absolute_footer
	DesignFooterAbsoluteFooter cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignHeadShortcutIcon = cfgmodel.NewStr(`design/head/shortcut_icon`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadDefaultTitle = cfgmodel.NewStr(`design/head/default_title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadTitlePrefix = cfgmodel.NewStr(`design/head/title_prefix`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadTitleSuffix = cfgmodel.NewStr(`design/head/title_suffix`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadDefaultDescription = cfgmodel.NewStr(`design/head/default_description`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadDefaultKeywords = cfgmodel.NewStr(`design/head/default_keywords`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadIncludes = cfgmodel.NewStr(`design/head/includes`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeadDemonotice = cfgmodel.NewBool(`design/head/demonotice`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignSearchEngineRobotsDefaultRobots = cfgmodel.NewStr(`design/search_engine_robots/default_robots`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignSearchEngineRobotsCustomInstructions = cfgmodel.NewStr(`design/search_engine_robots/custom_instructions`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignSearchEngineRobotsResetToDefaults = cfgmodel.NewStr(`design/search_engine_robots/reset_to_defaults`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeaderLogoSrc = cfgmodel.NewStr(`design/header/logo_src`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeaderLogoWidth = cfgmodel.NewStr(`design/header/logo_width`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeaderLogoHeight = cfgmodel.NewStr(`design/header/logo_height`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeaderLogoAlt = cfgmodel.NewStr(`design/header/logo_alt`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignHeaderWelcome = cfgmodel.NewStr(`design/header/welcome`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignFooterCopyright = cfgmodel.NewStr(`design/footer/copyright`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignFooterAbsoluteFooter = cfgmodel.NewStr(`design/footer/absolute_footer`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
