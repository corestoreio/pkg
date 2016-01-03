// +build ignore

package newrelicreporting

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        "newrelicreporting",
			Label:     `New Relic Reporting`,
			SortOrder: 1100,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_NewRelicReporting::config_newrelicreporting
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "general",
					Label:     `General`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: newrelicreporting/general/enable
							ID:        "enable",
							Label:     `Enable New Relic Integration`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: newrelicreporting/general/api_url
							ID:        "api_url",
							Label:     `New Relic API URL`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `https://api.newrelic.com/deployments.xml`,
						},

						&element.Field{
							// Path: newrelicreporting/general/insights_api_url
							ID:        "insights_api_url",
							Label:     `Insights API URL`,
							Comment:   element.LongText(`Use %s to replace the account ID in the URL`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `https://insights-collector.newrelic.com/v1/accounts/%s/events`,
						},

						&element.Field{
							// Path: newrelicreporting/general/account_id
							ID:        "account_id",
							Label:     `New Relic Account ID`,
							Comment:   element.LongText(`"Need a New Relic account? <a href="http://www.newrelic.com/magento" target="_blank">Click here to get one`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: newrelicreporting/general/app_id
							ID:        "app_id",
							Label:     `New Relic Application ID`,
							Comment:   element.LongText(`This can commonly be found at the end of the URL when viewing the APM after "/applications/"`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: newrelicreporting/general/api
							ID:        "api",
							Label:     `New Relic API Key`,
							Comment:   element.LongText(`This is located by navigating to Events -> Deployments from the New Relic APM website`),
							Type:      element.TypeObscure,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: newrelicreporting/general/insights_insert_key
							ID:        "insights_insert_key",
							Label:     `Insights API Key`,
							Comment:   element.LongText(`Generated under Insights in Manage data -> API Keys -> Insert Keys`),
							Type:      element.TypeObscure,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: newrelicreporting/general/app_name
							ID:        "app_name",
							Label:     `New Relic Application Name`,
							Comment:   element.LongText(`This is located by navigating to Settings from the New Relic APM website`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "cron",
					Label:     `Cron`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: newrelicreporting/cron/enable_cron
							ID:        "enable_cron",
							Label:     `Enable Cron`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
