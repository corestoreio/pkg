// +build ignore

package googleanalytics

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
			ID:        "google",
			Label:     `Google API`,
			SortOrder: 340,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_GoogleAnalytics::google
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "analytics",
					Label:     `Google Analytics`,
					SortOrder: 10,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: google/analytics/active
							ID:        "active",
							Label:     `Enable`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: google/analytics/account
							ID:        "account",
							Label:     `Account Number`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},
			),
		},
	)
}
