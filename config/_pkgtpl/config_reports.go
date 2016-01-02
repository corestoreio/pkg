// +build ignore

package reports

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
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "recently_products",
					Label:     `Recently Viewed/Compared Products`,
					SortOrder: 350,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/recently_products/scope
							ID:        "scope",
							Label:     `Show for Current`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `website`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Reports\Scope
						},

						&element.Field{
							// Path: catalog/recently_products/viewed_count
							ID:        "viewed_count",
							Label:     `Default Recently Viewed Products Count`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   5,
						},

						&element.Field{
							// Path: catalog/recently_products/compared_count
							ID:        "compared_count",
							Label:     `Default Recently Compared Products Count`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   5,
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "reports",
			Label:     `Reports`,
			SortOrder: 1000,
			Scope:     scope.NewPerm(scope.DefaultID),
			Resource:  0, // Otnegam_Reports::reports
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "dashboard",
					Label:     `Dashboard`,
					SortOrder: 1,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: reports/dashboard/ytd_start
							ID:        "ytd_start",
							Label:     `Year-To-Date Starts`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `1,1`,
						},

						&element.Field{
							// Path: reports/dashboard/mtd_start
							ID:        "mtd_start",
							Label:     `Current Month Starts`,
							Comment:   element.LongText(`Select day of the month.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   true,
						},
					),
				},
			),
		},
	)
	Path = NewPath(PackageConfiguration)
}
