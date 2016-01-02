// +build ignore

package theme

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDesignHeadShortcutIcon => Favicon Icon.
// Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers
// support all these formats!
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Favicon
var PathDesignHeadShortcutIcon = model.NewStr(`design/head/shortcut_icon`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadDefaultTitle => Default Title.
var PathDesignHeadDefaultTitle = model.NewStr(`design/head/default_title`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadTitlePrefix => Title Prefix.
var PathDesignHeadTitlePrefix = model.NewStr(`design/head/title_prefix`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadTitleSuffix => Title Suffix.
var PathDesignHeadTitleSuffix = model.NewStr(`design/head/title_suffix`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadDefaultDescription => Default Description.
var PathDesignHeadDefaultDescription = model.NewStr(`design/head/default_description`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadDefaultKeywords => Default Keywords.
var PathDesignHeadDefaultKeywords = model.NewStr(`design/head/default_keywords`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadIncludes => Miscellaneous Scripts.
// This will be included before head closing tag in page HTML.
var PathDesignHeadIncludes = model.NewStr(`design/head/includes`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeadDemonotice => Display Demo Store Notice.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDesignHeadDemonotice = model.NewBool(`design/head/demonotice`, model.WithPkgCfg(PackageConfiguration))

// PathDesignSearchEngineRobotsDefaultRobots => Default Robots.
// This will be included before head closing tag in page HTML.
// SourceModel: Otnegam\Config\Model\Config\Source\Design\Robots
var PathDesignSearchEngineRobotsDefaultRobots = model.NewStr(`design/search_engine_robots/default_robots`, model.WithPkgCfg(PackageConfiguration))

// PathDesignSearchEngineRobotsCustomInstructions => Edit custom instruction of robots.txt File.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Robots
var PathDesignSearchEngineRobotsCustomInstructions = model.NewStr(`design/search_engine_robots/custom_instructions`, model.WithPkgCfg(PackageConfiguration))

// PathDesignSearchEngineRobotsResetToDefaults => Reset to Defaults.
// This action will delete your custom instructions and reset robots.txt file
// to system's default settings.
var PathDesignSearchEngineRobotsResetToDefaults = model.NewStr(`design/search_engine_robots/reset_to_defaults`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeaderLogoSrc => Logo Image.
// Allowed file types:PNG, GIF, JPG, JPEG, SVG.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Logo
var PathDesignHeaderLogoSrc = model.NewStr(`design/header/logo_src`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeaderLogoWidth => Logo Image Width.
var PathDesignHeaderLogoWidth = model.NewStr(`design/header/logo_width`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeaderLogoHeight => Logo Image Height.
var PathDesignHeaderLogoHeight = model.NewStr(`design/header/logo_height`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeaderLogoAlt => Logo Image Alt.
var PathDesignHeaderLogoAlt = model.NewStr(`design/header/logo_alt`, model.WithPkgCfg(PackageConfiguration))

// PathDesignHeaderWelcome => Welcome Text.
var PathDesignHeaderWelcome = model.NewStr(`design/header/welcome`, model.WithPkgCfg(PackageConfiguration))

// PathDesignFooterCopyright => Copyright.
var PathDesignFooterCopyright = model.NewStr(`design/footer/copyright`, model.WithPkgCfg(PackageConfiguration))

// PathDesignFooterAbsoluteFooter => Miscellaneous HTML.
// This will be displayed just before body closing tag.
var PathDesignFooterAbsoluteFooter = model.NewStr(`design/footer/absolute_footer`, model.WithPkgCfg(PackageConfiguration))
