// +build ignore

package googleanalytics

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "google",
		Label:     `Google API`,
		SortOrder: 340,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_GoogleAnalytics::google
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "analytics",
				Label:     `Google Analytics`,
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: google/analytics/active
						ID:        "active",
						Label:     `Enable`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: google/analytics/account
						ID:        "account",
						Label:     `Account Number`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
)
