// +build ignore

package email

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
	// DesignEmailLogo => Logo Image.
	// Allowed file types: jpg, jpeg, gif, png. To optimize logo for
	// high-resolution displays, upload an image that is 3x normal size and then
	// specify 1x dimensions in width/height fields below.
	// Path: design/email/logo
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Logo
	DesignEmailLogo model.Str

	// DesignEmailLogoAlt => Logo Image Alt.
	// Path: design/email/logo_alt
	DesignEmailLogoAlt model.Str

	// DesignEmailLogoWidth => Logo Width.
	// Only necessary if image has been uploaded above. Enter number of pixels,
	// without appending "px".
	// Path: design/email/logo_width
	DesignEmailLogoWidth model.Str

	// DesignEmailLogoHeight => Logo Height.
	// Only necessary if image has been uploaded above. Enter number of pixels,
	// without appending "px".
	// Path: design/email/logo_height
	DesignEmailLogoHeight model.Str

	// DesignEmailHeaderTemplate => Header Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: design/email/header_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	DesignEmailHeaderTemplate model.Str

	// DesignEmailFooterTemplate => Footer Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: design/email/footer_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	DesignEmailFooterTemplate model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignEmailLogo = model.NewStr(`design/email/logo`, model.WithPkgCfg(pkgCfg))
	pp.DesignEmailLogoAlt = model.NewStr(`design/email/logo_alt`, model.WithPkgCfg(pkgCfg))
	pp.DesignEmailLogoWidth = model.NewStr(`design/email/logo_width`, model.WithPkgCfg(pkgCfg))
	pp.DesignEmailLogoHeight = model.NewStr(`design/email/logo_height`, model.WithPkgCfg(pkgCfg))
	pp.DesignEmailHeaderTemplate = model.NewStr(`design/email/header_template`, model.WithPkgCfg(pkgCfg))
	pp.DesignEmailFooterTemplate = model.NewStr(`design/email/footer_template`, model.WithPkgCfg(pkgCfg))

	return pp
}
