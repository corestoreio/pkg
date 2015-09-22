// +build ignore

package developer

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
				ID:        "front_end_development_workflow",
				Label:     `Front-end development workflow`,
				Comment:   ``,
				SortOrder: 8,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/front_end_development_workflow/type`,
						ID:           "type",
						Label:        `Workflow type`,
						Comment:      `Not available in production mode`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `server_side_compilation`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Developer\Model\Config\Source\WorkflowType
					},
				},
			},

			&config.Group{
				ID:        "restrict",
				Label:     `Developer Client Restrictions`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/restrict/allow_ips`,
						ID:           "allow_ips",
						Label:        `Allowed IPs (comma separated)`,
						Comment:      `Leave empty for access from any location.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Developer\Model\Config\Backend\AllowedIps
						SourceModel:  nil,
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
				ID: "restrict",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/restrict/allow_ips`,
						ID:      "allow_ips",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},
				},
			},
		},
	},
)
