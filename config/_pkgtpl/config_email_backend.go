// +build ignore

package email

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
	// DesignEmailLogo => Logo Image.
	// Allowed file types: jpg, jpeg, gif, png. To optimize logo for
	// high-resolution displays, upload an image that is 3x normal size and then
	// specify 1x dimensions in width/height fields below.
	// Path: design/email/logo
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Logo
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
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	DesignEmailHeaderTemplate model.Str

	// DesignEmailFooterTemplate => Footer Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: design/email/footer_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	DesignEmailFooterTemplate model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignEmailLogo = model.NewStr(`design/email/logo`, model.WithConfigStructure(cfgStruct))
	pp.DesignEmailLogoAlt = model.NewStr(`design/email/logo_alt`, model.WithConfigStructure(cfgStruct))
	pp.DesignEmailLogoWidth = model.NewStr(`design/email/logo_width`, model.WithConfigStructure(cfgStruct))
	pp.DesignEmailLogoHeight = model.NewStr(`design/email/logo_height`, model.WithConfigStructure(cfgStruct))
	pp.DesignEmailHeaderTemplate = model.NewStr(`design/email/header_template`, model.WithConfigStructure(cfgStruct))
	pp.DesignEmailFooterTemplate = model.NewStr(`design/email/footer_template`, model.WithConfigStructure(cfgStruct))

	return pp
}
