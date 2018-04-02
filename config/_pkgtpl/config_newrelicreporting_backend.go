// +build ignore

package newrelicreporting

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
	// NewrelicreportingGeneralEnable => Enable New Relic Integration.
	// Path: newrelicreporting/general/enable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewrelicreportingGeneralEnable cfgmodel.Bool

	// NewrelicreportingGeneralApiUrl => New Relic API URL.
	// Path: newrelicreporting/general/api_url
	NewrelicreportingGeneralApiUrl cfgmodel.Str

	// NewrelicreportingGeneralInsightsApiUrl => Insights API URL.
	// Use %s to replace the account ID in the URL
	// Path: newrelicreporting/general/insights_api_url
	NewrelicreportingGeneralInsightsApiUrl cfgmodel.Str

	// NewrelicreportingGeneralAccountId => New Relic Account ID.
	// "Need a New Relic account? Click here to get one
	// Path: newrelicreporting/general/account_id
	NewrelicreportingGeneralAccountId cfgmodel.Str

	// NewrelicreportingGeneralAppId => New Relic Application ID.
	// This can commonly be found at the end of the URL when viewing the APM after
	// "/applications/"
	// Path: newrelicreporting/general/app_id
	NewrelicreportingGeneralAppId cfgmodel.Str

	// NewrelicreportingGeneralApi => New Relic API Key.
	// This is located by navigating to Events -> Deployments from the New Relic
	// APM website
	// Path: newrelicreporting/general/api
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralApi cfgmodel.Str

	// NewrelicreportingGeneralInsightsInsertKey => Insights API Key.
	// Generated under Insights in Manage data -> API Keys -> Insert Keys
	// Path: newrelicreporting/general/insights_insert_key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	NewrelicreportingGeneralInsightsInsertKey cfgmodel.Str

	// NewrelicreportingGeneralAppName => New Relic Application Name.
	// This is located by navigating to Settings from the New Relic APM website
	// Path: newrelicreporting/general/app_name
	NewrelicreportingGeneralAppName cfgmodel.Str

	// NewrelicreportingCronEnableCron => Enable Cron.
	// Path: newrelicreporting/cron/enable_cron
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	NewrelicreportingCronEnableCron cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.NewrelicreportingGeneralEnable = cfgmodel.NewBool(`newrelicreporting/general/enable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralApiUrl = cfgmodel.NewStr(`newrelicreporting/general/api_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralInsightsApiUrl = cfgmodel.NewStr(`newrelicreporting/general/insights_api_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralAccountId = cfgmodel.NewStr(`newrelicreporting/general/account_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralAppId = cfgmodel.NewStr(`newrelicreporting/general/app_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralApi = cfgmodel.NewStr(`newrelicreporting/general/api`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralInsightsInsertKey = cfgmodel.NewStr(`newrelicreporting/general/insights_insert_key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingGeneralAppName = cfgmodel.NewStr(`newrelicreporting/general/app_name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.NewrelicreportingCronEnableCron = cfgmodel.NewBool(`newrelicreporting/cron/enable_cron`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
