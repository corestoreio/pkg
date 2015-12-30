// +build ignore

package theme

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDesignHeadShortcutIcon => Favicon Icon.
// Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers
// support all these formats!
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Favicon
var PathDesignHeadShortcutIcon = model.NewStr(`design/head/shortcut_icon`)

// PathDesignHeadDefaultTitle => Default Title.
var PathDesignHeadDefaultTitle = model.NewStr(`design/head/default_title`)

// PathDesignHeadTitlePrefix => Title Prefix.
var PathDesignHeadTitlePrefix = model.NewStr(`design/head/title_prefix`)

// PathDesignHeadTitleSuffix => Title Suffix.
var PathDesignHeadTitleSuffix = model.NewStr(`design/head/title_suffix`)

// PathDesignHeadDefaultDescription => Default Description.
var PathDesignHeadDefaultDescription = model.NewStr(`design/head/default_description`)

// PathDesignHeadDefaultKeywords => Default Keywords.
var PathDesignHeadDefaultKeywords = model.NewStr(`design/head/default_keywords`)

// PathDesignHeadIncludes => Miscellaneous Scripts.
// This will be included before head closing tag in page HTML.
var PathDesignHeadIncludes = model.NewStr(`design/head/includes`)

// PathDesignHeadDemonotice => Display Demo Store Notice.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDesignHeadDemonotice = model.NewBool(`design/head/demonotice`)

// PathDesignSearchEngineRobotsDefaultRobots => Default Robots.
// This will be included before head closing tag in page HTML.
// SourceModel: Otnegam\Config\Model\Config\Source\Design\Robots
var PathDesignSearchEngineRobotsDefaultRobots = model.NewStr(`design/search_engine_robots/default_robots`)

// PathDesignSearchEngineRobotsCustomInstructions => Edit custom instruction of robots.txt File.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Robots
var PathDesignSearchEngineRobotsCustomInstructions = model.NewStr(`design/search_engine_robots/custom_instructions`)

// PathDesignSearchEngineRobotsResetToDefaults => Reset to Defaults.
// This action will delete your custom instructions and reset robots.txt file
// to system's default settings.
var PathDesignSearchEngineRobotsResetToDefaults = model.NewStr(`design/search_engine_robots/reset_to_defaults`)

// PathDesignHeaderLogoSrc => Logo Image.
// Allowed file types:PNG, GIF, JPG, JPEG, SVG.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Logo
var PathDesignHeaderLogoSrc = model.NewStr(`design/header/logo_src`)

// PathDesignHeaderLogoWidth => Logo Image Width.
var PathDesignHeaderLogoWidth = model.NewStr(`design/header/logo_width`)

// PathDesignHeaderLogoHeight => Logo Image Height.
var PathDesignHeaderLogoHeight = model.NewStr(`design/header/logo_height`)

// PathDesignHeaderLogoAlt => Logo Image Alt.
var PathDesignHeaderLogoAlt = model.NewStr(`design/header/logo_alt`)

// PathDesignHeaderWelcome => Welcome Text.
var PathDesignHeaderWelcome = model.NewStr(`design/header/welcome`)

// PathDesignFooterCopyright => Copyright.
var PathDesignFooterCopyright = model.NewStr(`design/footer/copyright`)

// PathDesignFooterAbsoluteFooter => Miscellaneous HTML.
// This will be displayed just before body closing tag.
var PathDesignFooterAbsoluteFooter = model.NewStr(`design/footer/absolute_footer`)
