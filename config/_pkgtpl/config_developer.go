// +build ignore

package developer

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "dev",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "front_end_development_workflow",
				Label:     `Frontend Development Workflow`,
				SortOrder: 8,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/front_end_development_workflow/type
						ID:        "type",
						Label:     `Workflow type`,
						Comment:   element.LongText(`Not available in production mode`),
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `server_side_compilation`,
						// SourceModel: Otnegam\Developer\Model\Config\Source\WorkflowType
					},
				),
			},

			&config.Group{
				ID:        "restrict",
				Label:     `Developer Client Restrictions`,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/restrict/allow_ips
						ID:        "allow_ips",
						Label:     `Allowed IPs (comma separated)`,
						Comment:   element.LongText(`Leave empty for access from any location.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Developer\Model\Config\Backend\AllowedIps
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "dev",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "restrict",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/restrict/allow_ips
						ID:      `allow_ips`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},
				),
			},
		),
	},
)
