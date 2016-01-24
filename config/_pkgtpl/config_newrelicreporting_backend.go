// +build ignore

package newrelicreporting

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
	// NewrelicreportingGeneralEnable => Enable New Relic Integration.
	// Path: newrelicreporting/general/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
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
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralApi model.Str

	// NewrelicreportingGeneralInsightsInsertKey => Insights API Key.
	// Generated under Insights in Manage data -> API Keys -> Insert Keys
	// Path: newrelicreporting/general/insights_insert_key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralInsightsInsertKey model.Str

	// NewrelicreportingGeneralAppName => New Relic Application Name.
	// This is located by navigating to Settings from the New Relic APM website
	// Path: newrelicreporting/general/app_name
	NewrelicreportingGeneralAppName model.Str

	// NewrelicreportingCronEnableCron => Enable Cron.
	// Path: newrelicreporting/cron/enable_cron
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewrelicreportingCronEnableCron model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.NewrelicreportingGeneralEnable = model.NewBool(`newrelicreporting/general/enable`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralApiUrl = model.NewStr(`newrelicreporting/general/api_url`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralInsightsApiUrl = model.NewStr(`newrelicreporting/general/insights_api_url`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralAccountId = model.NewStr(`newrelicreporting/general/account_id`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralAppId = model.NewStr(`newrelicreporting/general/app_id`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralApi = model.NewStr(`newrelicreporting/general/api`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralInsightsInsertKey = model.NewStr(`newrelicreporting/general/insights_insert_key`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingGeneralAppName = model.NewStr(`newrelicreporting/general/app_name`, model.WithConfigStructure(cfgStruct))
	pp.NewrelicreportingCronEnableCron = model.NewBool(`newrelicreporting/cron/enable_cron`, model.WithConfigStructure(cfgStruct))

	return pp
}
