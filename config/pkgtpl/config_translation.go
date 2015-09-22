// +build ignore

package translation

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "dev",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "js",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/js/translate_strategy`,
						ID:           "translate_strategy",
						Label:        `Translation Strategy`,
						Comment:      `Please put your store into maintenance mode and redeploy static files after changing strategy`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `dictionary`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Translation\Model\Js\Config\Source\Strategy
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "dev",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "translate_inline",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/translate_inline/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `dev/translate_inline/active_admin`,
						ID:      "active_admin",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `dev/translate_inline/invalid_caches`,
						ID:      "invalid_caches",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"block_html":null}`,
					},
				},
			},
		},
	},
)
