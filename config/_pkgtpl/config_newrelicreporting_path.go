// +build ignore

package newrelicreporting

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathNewrelicreportingGeneralEnable => Enable New Relic Integration.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathNewrelicreportingGeneralEnable = model.NewBool(`newrelicreporting/general/enable`)

// PathNewrelicreportingGeneralApiUrl => New Relic API URL.
var PathNewrelicreportingGeneralApiUrl = model.NewStr(`newrelicreporting/general/api_url`)

// PathNewrelicreportingGeneralInsightsApiUrl => Insights API URL.
// Use %s to replace the account ID in the URL
var PathNewrelicreportingGeneralInsightsApiUrl = model.NewStr(`newrelicreporting/general/insights_api_url`)

// PathNewrelicreportingGeneralAccountId => New Relic Account ID.
// "Need a New Relic account? Click here to get one
var PathNewrelicreportingGeneralAccountId = model.NewStr(`newrelicreporting/general/account_id`)

// PathNewrelicreportingGeneralAppId => New Relic Application ID.
// This can commonly be found at the end of the URL when viewing the APM after
// "/applications/"
var PathNewrelicreportingGeneralAppId = model.NewStr(`newrelicreporting/general/app_id`)

// PathNewrelicreportingGeneralApi => New Relic API Key.
// This is located by navigating to Events -> Deployments from the New Relic
// APM website
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathNewrelicreportingGeneralApi = model.NewStr(`newrelicreporting/general/api`)

// PathNewrelicreportingGeneralInsightsInsertKey => Insights API Key.
// Generated under Insights in Manage data -> API Keys -> Insert Keys
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathNewrelicreportingGeneralInsightsInsertKey = model.NewStr(`newrelicreporting/general/insights_insert_key`)

// PathNewrelicreportingGeneralAppName => New Relic Application Name.
// This is located by navigating to Settings from the New Relic APM website
var PathNewrelicreportingGeneralAppName = model.NewStr(`newrelicreporting/general/app_name`)

// PathNewrelicreportingCronEnableCron => Enable Cron.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathNewrelicreportingCronEnableCron = model.NewBool(`newrelicreporting/cron/enable_cron`)
