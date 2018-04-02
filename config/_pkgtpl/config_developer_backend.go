// +build ignore

package developer

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
	// DevFrontEndDevelopmentWorkflowType => Workflow type.
	// Not available in production mode
	// Path: dev/front_end_development_workflow/type
	// SourceModel: Magento\Developer\Model\Config\Source\WorkflowType
	DevFrontEndDevelopmentWorkflowType cfgmodel.Str

	// DevRestrictAllowIps => Allowed IPs (comma separated).
	// Leave empty for access from any location.
	// Path: dev/restrict/allow_ips
	// BackendModel: Magento\Developer\Model\Config\Backend\AllowedIps
	DevRestrictAllowIps cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.DevFrontEndDevelopmentWorkflowType = cfgmodel.NewStr(`dev/front_end_development_workflow/type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DevRestrictAllowIps = cfgmodel.NewStr(`dev/restrict/allow_ips`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
