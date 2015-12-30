// +build ignore

package newrelicreporting

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "newrelicreporting",
		Label:     `New Relic Reporting`,
		SortOrder: 1100,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_NewRelicReporting::config_newrelicreporting
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "general",
				Label:     `General`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: newrelicreporting/general/enable
						ID:        "enable",
						Label:     `Enable New Relic Integration`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: newrelicreporting/general/api_url
						ID:        "api_url",
						Label:     `New Relic API URL`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `https://api.newrelic.com/deployments.xml`,
					},

					&config.Field{
						// Path: newrelicreporting/general/insights_api_url
						ID:        "insights_api_url",
						Label:     `Insights API URL`,
						Comment:   element.LongText(`Use %s to replace the account ID in the URL`),
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `https://insights-collector.newrelic.com/v1/accounts/%s/events`,
					},

					&config.Field{
						// Path: newrelicreporting/general/account_id
						ID:        "account_id",
						Label:     `New Relic Account ID`,
						Comment:   element.LongText(`"Need a New Relic account? <a href="http://www.newrelic.com/magento" target="_blank">Click here to get one`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: newrelicreporting/general/app_id
						ID:        "app_id",
						Label:     `New Relic Application ID`,
						Comment:   element.LongText(`This can commonly be found at the end of the URL when viewing the APM after "/applications/"`),
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: newrelicreporting/general/api
						ID:        "api",
						Label:     `New Relic API Key`,
						Comment:   element.LongText(`This is located by navigating to Events -> Deployments from the New Relic APM website`),
						Type:      config.TypeObscure,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: newrelicreporting/general/insights_insert_key
						ID:        "insights_insert_key",
						Label:     `Insights API Key`,
						Comment:   element.LongText(`Generated under Insights in Manage data -> API Keys -> Insert Keys`),
						Type:      config.TypeObscure,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
					},

					&config.Field{
						// Path: newrelicreporting/general/app_name
						ID:        "app_name",
						Label:     `New Relic Application Name`,
						Comment:   element.LongText(`This is located by navigating to Settings from the New Relic APM website`),
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},

			&config.Group{
				ID:        "cron",
				Label:     `Cron`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: newrelicreporting/cron/enable_cron
						ID:        "enable_cron",
						Label:     `Enable Cron`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
