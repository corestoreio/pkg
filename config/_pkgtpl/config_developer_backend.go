// +build ignore

package developer

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
	// DevFrontEndDevelopmentWorkflowType => Workflow type.
	// Not available in production mode
	// Path: dev/front_end_development_workflow/type
	// SourceModel: Otnegam\Developer\Model\Config\Source\WorkflowType
	DevFrontEndDevelopmentWorkflowType model.Str

	// DevRestrictAllowIps => Allowed IPs (comma separated).
	// Leave empty for access from any location.
	// Path: dev/restrict/allow_ips
	// BackendModel: Otnegam\Developer\Model\Config\Backend\AllowedIps
	DevRestrictAllowIps model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevFrontEndDevelopmentWorkflowType = model.NewStr(`dev/front_end_development_workflow/type`, model.WithConfigStructure(cfgStruct))
	pp.DevRestrictAllowIps = model.NewStr(`dev/restrict/allow_ips`, model.WithConfigStructure(cfgStruct))

	return pp
}
