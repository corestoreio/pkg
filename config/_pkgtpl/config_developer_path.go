// +build ignore

package developer

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
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.DevFrontEndDevelopmentWorkflowType = model.NewStr(`dev/front_end_development_workflow/type`, model.WithPkgCfg(pkgCfg))
	pp.DevRestrictAllowIps = model.NewStr(`dev/restrict/allow_ips`, model.WithPkgCfg(pkgCfg))

	return pp
}
