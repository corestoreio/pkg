// +build ignore

package integration

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "oauth",
		Label:     "OAuth",
		SortOrder: 300,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cleanup",
				Label:     `Cleanup Settings`,
				Comment:   ``,
				SortOrder: 300,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `oauth/cleanup/cleanup_probability`,
						ID:           "cleanup_probability",
						Label:        `Cleanup Probability`,
						Comment:      `Integer. Launch cleanup in X OAuth requests. 0 (not recommended) - to disable cleanup`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      100,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `oauth/cleanup/expiration_period`,
						ID:           "expiration_period",
						Label:        `Expiration Period`,
						Comment:      `Cleanup entries older than X minutes.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      120,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "consumer",
				Label:     `Consumer Settings`,
				Comment:   ``,
				SortOrder: 400,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `oauth/consumer/expiration_period`,
						ID:           "expiration_period",
						Label:        `Expiration Period`,
						Comment:      `Consumer key/secret will expire if not used within X seconds after Oauth token exchange starts.`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      300,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `oauth/consumer/post_maxredirects`,
						ID:           "post_maxredirects",
						Label:        `OAuth consumer credentials HTTP Post maxredirects`,
						Comment:      `Number of maximum redirects for OAuth consumer credentials Post request.`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      0,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `oauth/consumer/post_timeout`,
						ID:           "post_timeout",
						Label:        `OAuth consumer credentials HTTP Post timeout`,
						Comment:      `Timeout for OAuth consumer credentials Post request within X seconds.`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      5,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
