// +build ignore

package ui

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
						// Path: `dev/js/session_storage_logging`,
						ID:           "session_storage_logging",
						Label:        `Log JS Errors to Session Storage`,
						Comment:      `If enabled, can be used by functional tests for extended reporting`,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `dev/js/session_storage_key`,
						ID:           "session_storage_key",
						Label:        `Log JS Errors to Session Storage Key`,
						Comment:      `Use this key to retrieve collected js errors`,
						Type:         config.TypeText,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `collected_errors`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
