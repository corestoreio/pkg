// +build ignore

package integration

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
			ID:        "oauth",
			Label:     `OAuth`,
			SortOrder: 300,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Integration::config_oauth
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "cleanup",
					Label:     `Cleanup Settings`,
					SortOrder: 300,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: oauth/cleanup/cleanup_probability
							ID:        "cleanup_probability",
							Label:     `Cleanup Probability`,
							Comment:   text.Long(`Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to disable cleanup`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   100,
						},

						&element.Field{
							// Path: oauth/cleanup/expiration_period
							ID:        "expiration_period",
							Label:     `Expiration Period`,
							Comment:   text.Long(`Cleanup entries older than X minutes.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   120,
						},
					),
				},

				&element.Group{
					ID:        "consumer",
					Label:     `Consumer Settings`,
					SortOrder: 400,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: oauth/consumer/expiration_period
							ID:        "expiration_period",
							Label:     `Expiration Period`,
							Comment:   text.Long(`Consumer key/secret will expire if not used within X seconds after Oauth token exchange starts.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   300,
						},

						&element.Field{
							// Path: oauth/consumer/post_maxredirects
							ID:        "post_maxredirects",
							Label:     `OAuth consumer credentials HTTP Post maxredirects`,
							Comment:   text.Long(`Number of maximum redirects for OAuth consumer credentials Post request.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						&element.Field{
							// Path: oauth/consumer/post_timeout
							ID:        "post_timeout",
							Label:     `OAuth consumer credentials HTTP Post timeout`,
							Comment:   text.Long(`Timeout for OAuth consumer credentials Post request within X seconds.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   5,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
