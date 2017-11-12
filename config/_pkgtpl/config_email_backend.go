// +build ignore

package email

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
	// DesignEmailLogo => Logo Image.
	// Allowed file types: jpg, jpeg, gif, png. To optimize logo for
	// high-resolution displays, upload an image that is 3x normal size and then
	// specify 1x dimensions in width/height fields below.
	// Path: design/email/logo
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Logo
	DesignEmailLogo cfgmodel.Str

	// DesignEmailLogoAlt => Logo Image Alt.
	// Path: design/email/logo_alt
	DesignEmailLogoAlt cfgmodel.Str

	// DesignEmailLogoWidth => Logo Width.
	// Only necessary if image has been uploaded above. Enter number of pixels,
	// without appending "px".
	// Path: design/email/logo_width
	DesignEmailLogoWidth cfgmodel.Str

	// DesignEmailLogoHeight => Logo Height.
	// Only necessary if image has been uploaded above. Enter number of pixels,
	// without appending "px".
	// Path: design/email/logo_height
	DesignEmailLogoHeight cfgmodel.Str

	// DesignEmailHeaderTemplate => Header Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: design/email/header_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	DesignEmailHeaderTemplate cfgmodel.Str

	// DesignEmailFooterTemplate => Footer Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: design/email/footer_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	DesignEmailFooterTemplate cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DesignEmailLogo = cfgmodel.NewStr(`design/email/logo`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignEmailLogoAlt = cfgmodel.NewStr(`design/email/logo_alt`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignEmailLogoWidth = cfgmodel.NewStr(`design/email/logo_width`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignEmailLogoHeight = cfgmodel.NewStr(`design/email/logo_height`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignEmailHeaderTemplate = cfgmodel.NewStr(`design/email/header_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignEmailFooterTemplate = cfgmodel.NewStr(`design/email/footer_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
