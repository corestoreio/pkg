// +build ignore

package integration

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "oauth",
		Label:     `OAuth`,
		SortOrder: 300,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Integration::config_oauth
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "cleanup",
				Label:     `Cleanup Settings`,
				SortOrder: 300,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: oauth/cleanup/cleanup_probability
						ID:        "cleanup_probability",
						Label:     `Cleanup Probability`,
						Comment:   element.LongText(`Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to disable cleanup`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   100,
					},

					&config.Field{
						// Path: oauth/cleanup/expiration_period
						ID:        "expiration_period",
						Label:     `Expiration Period`,
						Comment:   element.LongText(`Cleanup entries older than X minutes.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   120,
					},
				),
			},

			&config.Group{
				ID:        "consumer",
				Label:     `Consumer Settings`,
				SortOrder: 400,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: oauth/consumer/expiration_period
						ID:        "expiration_period",
						Label:     `Expiration Period`,
						Comment:   element.LongText(`Consumer key/secret will expire if not used within X seconds after Oauth token exchange starts.`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   300,
					},

					&config.Field{
						// Path: oauth/consumer/post_maxredirects
						ID:        "post_maxredirects",
						Label:     `OAuth consumer credentials HTTP Post maxredirects`,
						Comment:   element.LongText(`Number of maximum redirects for OAuth consumer credentials Post request.`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: oauth/consumer/post_timeout
						ID:        "post_timeout",
						Label:     `OAuth consumer credentials HTTP Post timeout`,
						Comment:   element.LongText(`Timeout for OAuth consumer credentials Post request within X seconds.`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   5,
					},
				),
			},
		),
	},
)
