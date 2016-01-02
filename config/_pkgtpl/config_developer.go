// +build ignore

package developer

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package.
// Used in frontend and backend. See init() for details.
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(
		&element.Section{
			ID: "dev",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "front_end_development_workflow",
					Label:     `Frontend Development Workflow`,
					SortOrder: 8,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/front_end_development_workflow/type
							ID:        "type",
							Label:     `Workflow type`,
							Comment:   element.LongText(`Not available in production mode`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `server_side_compilation`,
							// SourceModel: Otnegam\Developer\Model\Config\Source\WorkflowType
						},
					),
				},

				&element.Group{
					ID:        "restrict",
					Label:     `Developer Client Restrictions`,
					SortOrder: 10,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/restrict/allow_ips
							ID:        "allow_ips",
							Label:     `Allowed IPs (comma separated)`,
							Comment:   element.LongText(`Leave empty for access from any location.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Developer\Model\Config\Backend\AllowedIps
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "dev",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "restrict",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/restrict/allow_ips
							ID:      `allow_ips`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},
			),
		},
	)
	Path = NewPath(PackageConfiguration)
}
