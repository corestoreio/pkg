// +build ignore

package reports

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
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "recently_products",
					Label:     `Recently Viewed/Compared Products`,
					SortOrder: 350,
					Scope:     scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/recently_products/scope
							ID:        "scope",
							Label:     `Show for Current`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `website`,
							// SourceModel: Magento\Config\Model\Config\Source\Reports\Scope
						},

						&element.Field{
							// Path: catalog/recently_products/viewed_count
							ID:        "viewed_count",
							Label:     `Default Recently Viewed Products Count`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   5,
						},

						&element.Field{
							// Path: catalog/recently_products/compared_count
							ID:        "compared_count",
							Label:     `Default Recently Compared Products Count`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
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
			Scope:     scope.PermDefault,
			Resource:  0, // Magento_Reports::reports
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "dashboard",
					Label:     `Dashboard`,
					SortOrder: 1,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: reports/dashboard/ytd_start
							ID:        "ytd_start",
							Label:     `Year-To-Date Starts`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `1,1`,
						},

						&element.Field{
							// Path: reports/dashboard/mtd_start
							ID:        "mtd_start",
							Label:     `Current Month Starts`,
							Comment:   text.Long(`Select day of the month.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   true,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
