// +build ignore

package newrelicreporting

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
	// NewrelicreportingGeneralEnable => Enable New Relic Integration.
	// Path: newrelicreporting/general/enable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	NewrelicreportingGeneralEnable model.Bool

	// NewrelicreportingGeneralApiUrl => New Relic API URL.
	// Path: newrelicreporting/general/api_url
	NewrelicreportingGeneralApiUrl model.Str

	// NewrelicreportingGeneralInsightsApiUrl => Insights API URL.
	// Use %s to replace the account ID in the URL
	// Path: newrelicreporting/general/insights_api_url
	NewrelicreportingGeneralInsightsApiUrl model.Str

	// NewrelicreportingGeneralAccountId => New Relic Account ID.
	// "Need a New Relic account? Click here to get one
	// Path: newrelicreporting/general/account_id
	NewrelicreportingGeneralAccountId model.Str

	// NewrelicreportingGeneralAppId => New Relic Application ID.
	// This can commonly be found at the end of the URL when viewing the APM after
	// "/applications/"
	// Path: newrelicreporting/general/app_id
	NewrelicreportingGeneralAppId model.Str

	// NewrelicreportingGeneralApi => New Relic API Key.
	// This is located by navigating to Events -> Deployments from the New Relic
	// APM website
	// Path: newrelicreporting/general/api
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralApi model.Str

	// NewrelicreportingGeneralInsightsInsertKey => Insights API Key.
	// Generated under Insights in Manage data -> API Keys -> Insert Keys
	// Path: newrelicreporting/general/insights_insert_key
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralInsightsInsertKey model.Str

	// NewrelicreportingGeneralAppName => New Relic Application Name.
	// This is located by navigating to Settings from the New Relic APM website
	// Path: newrelicreporting/general/app_name
	NewrelicreportingGeneralAppName model.Str

	// NewrelicreportingCronEnableCron => Enable Cron.
	// Path: newrelicreporting/cron/enable_cron
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	NewrelicreportingCronEnableCron model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.NewrelicreportingGeneralEnable = model.NewBool(`newrelicreporting/general/enable`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralApiUrl = model.NewStr(`newrelicreporting/general/api_url`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralInsightsApiUrl = model.NewStr(`newrelicreporting/general/insights_api_url`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralAccountId = model.NewStr(`newrelicreporting/general/account_id`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralAppId = model.NewStr(`newrelicreporting/general/app_id`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralApi = model.NewStr(`newrelicreporting/general/api`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralInsightsInsertKey = model.NewStr(`newrelicreporting/general/insights_insert_key`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingGeneralAppName = model.NewStr(`newrelicreporting/general/app_name`, model.WithPkgCfg(pkgCfg))
	pp.NewrelicreportingCronEnableCron = model.NewBool(`newrelicreporting/cron/enable_cron`, model.WithPkgCfg(pkgCfg))

	return pp
}
