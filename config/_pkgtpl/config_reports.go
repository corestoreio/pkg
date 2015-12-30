// +build ignore

package reports

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "recently_products",
				Label:     `Recently Viewed/Compared Products`,
				SortOrder: 350,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/recently_products/scope
						ID:        "scope",
						Label:     `Show for Current`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `website`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Reports\Scope
					},

					&config.Field{
						// Path: catalog/recently_products/viewed_count
						ID:        "viewed_count",
						Label:     `Default Recently Viewed Products Count`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   5,
					},

					&config.Field{
						// Path: catalog/recently_products/compared_count
						ID:        "compared_count",
						Label:     `Default Recently Compared Products Count`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   5,
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "reports",
		Label:     `Reports`,
		SortOrder: 1000,
		Scope:     scope.NewPerm(scope.DefaultID),
		Resource:  0, // Otnegam_Reports::reports
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "dashboard",
				Label:     `Dashboard`,
				SortOrder: 1,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: reports/dashboard/ytd_start
						ID:        "ytd_start",
						Label:     `Year-To-Date Starts`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `1,1`,
					},

					&config.Field{
						// Path: reports/dashboard/mtd_start
						ID:        "mtd_start",
						Label:     `Current Month Starts`,
						Comment:   element.LongText(`Select day of the month.`),
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
					},
				),
			},
		),
	},
)
