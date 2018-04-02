// +build ignore

package newrelicreporting

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID:        "newrelicreporting",
			Label:     `New Relic Reporting`,
			SortOrder: 1100,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_NewRelicReporting::config_newrelicreporting
			Groups: element.MakeGroups(
				element.Group{
					ID:        "general",
					Label:     `General`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: newrelicreporting/general/enable
							ID:        "enable",
							Label:     `Enable New Relic Integration`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: newrelicreporting/general/api_url
							ID:        "api_url",
							Label:     `New Relic API URL`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `https://api.newrelic.com/deployments.xml`,
						},

						element.Field{
							// Path: newrelicreporting/general/insights_api_url
							ID:        "insights_api_url",
							Label:     `Insights API URL`,
							Comment:   text.Long(`Use %s to replace the account ID in the URL`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `https://insights-collector.newrelic.com/v1/accounts/%s/events`,
						},

						element.Field{
							// Path: newrelicreporting/general/account_id
							ID:        "account_id",
							Label:     `New Relic Account ID`,
							Comment:   text.Long(`"Need a New Relic account? <a href="http://www.newrelic.com/magento" target="_blank">Click here to get one`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: newrelicreporting/general/app_id
							ID:        "app_id",
							Label:     `New Relic Application ID`,
							Comment:   text.Long(`This can commonly be found at the end of the URL when viewing the APM after "/applications/"`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: newrelicreporting/general/api
							ID:        "api",
							Label:     `New Relic API Key`,
							Comment:   text.Long(`This is located by navigating to Events -> Deployments from the New Relic APM website`),
							Type:      element.TypeObscure,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: newrelicreporting/general/insights_insert_key
							ID:        "insights_insert_key",
							Label:     `Insights API Key`,
							Comment:   text.Long(`Generated under Insights in Manage data -> API Keys -> Insert Keys`),
							Type:      element.TypeObscure,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
						},

						element.Field{
							// Path: newrelicreporting/general/app_name
							ID:        "app_name",
							Label:     `New Relic Application Name`,
							Comment:   text.Long(`This is located by navigating to Settings from the New Relic APM website`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},

				element.Group{
					ID:        "cron",
					Label:     `Cron`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: newrelicreporting/cron/enable_cron
							ID:        "enable_cron",
							Label:     `Enable Cron`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
