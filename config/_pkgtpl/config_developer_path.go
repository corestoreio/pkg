// +build ignore

package developer

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathDevFrontEndDevelopmentWorkflowType => Workflow type.
// Not available in production mode
// SourceModel: Otnegam\Developer\Model\Config\Source\WorkflowType
var PathDevFrontEndDevelopmentWorkflowType = model.NewStr(`dev/front_end_development_workflow/type`, model.WithPkgCfg(PackageConfiguration))

// PathDevRestrictAllowIps => Allowed IPs (comma separated).
// Leave empty for access from any location.
// BackendModel: Otnegam\Developer\Model\Config\Backend\AllowedIps
var PathDevRestrictAllowIps = model.NewStr(`dev/restrict/allow_ips`, model.WithPkgCfg(PackageConfiguration))
