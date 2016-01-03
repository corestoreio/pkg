// +build ignore

package developer

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
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

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DevFrontEndDevelopmentWorkflowType = model.NewStr(`dev/front_end_development_workflow/type`, model.WithConfigStructure(cfgStruct))
	pp.DevRestrictAllowIps = model.NewStr(`dev/restrict/allow_ips`, model.WithConfigStructure(cfgStruct))

	return pp
}
