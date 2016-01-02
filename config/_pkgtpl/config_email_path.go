// +build ignore

package email

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDesignEmailLogo => Logo Image.
// Allowed file types: jpg, jpeg, gif, png. To optimize logo for
// high-resolution displays, upload an image that is 3x normal size and then
// specify 1x dimensions in width/height fields below.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Logo
var PathDesignEmailLogo = model.NewStr(`design/email/logo`, model.WithPkgCfg(PackageConfiguration))

// PathDesignEmailLogoAlt => Logo Image Alt.
var PathDesignEmailLogoAlt = model.NewStr(`design/email/logo_alt`, model.WithPkgCfg(PackageConfiguration))

// PathDesignEmailLogoWidth => Logo Width.
// Only necessary if image has been uploaded above. Enter number of pixels,
// without appending "px".
var PathDesignEmailLogoWidth = model.NewStr(`design/email/logo_width`, model.WithPkgCfg(PackageConfiguration))

// PathDesignEmailLogoHeight => Logo Height.
// Only necessary if image has been uploaded above. Enter number of pixels,
// without appending "px".
var PathDesignEmailLogoHeight = model.NewStr(`design/email/logo_height`, model.WithPkgCfg(PackageConfiguration))

// PathDesignEmailHeaderTemplate => Header Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathDesignEmailHeaderTemplate = model.NewStr(`design/email/header_template`, model.WithPkgCfg(PackageConfiguration))

// PathDesignEmailFooterTemplate => Footer Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathDesignEmailFooterTemplate = model.NewStr(`design/email/footer_template`, model.WithPkgCfg(PackageConfiguration))
