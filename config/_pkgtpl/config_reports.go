// +build ignore

package reports

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "recently_products",
				Label:     `Recently Viewed/Compared Products`,
				Comment:   ``,
				SortOrder: 350,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/recently_products/scope`,
						ID:           "scope",
						Label:        `Show for Current`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `website`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Reports\Scope
					},

					&config.Field{
						// Path: `catalog/recently_products/viewed_count`,
						ID:           "viewed_count",
						Label:        `Default Recently Viewed Products Count`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      5,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/recently_products/compared_count`,
						ID:           "compared_count",
						Label:        `Default Recently Compared Products Count`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      5,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "reports",
		Label:     "Reports",
		SortOrder: 1000,
		Scope:     scope.NewPerm(scope.DefaultID),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "dashboard",
				Label:     `Dashboard`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `reports/dashboard/ytd_start`,
						ID:           "ytd_start",
						Label:        `Year-To-Date Starts`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `1,1`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `reports/dashboard/mtd_start`,
						ID:           "mtd_start",
						Label:        `Current Month Starts`,
						Comment:      `Select day of the month.`,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
