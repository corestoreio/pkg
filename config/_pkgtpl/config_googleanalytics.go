// +build ignore

package googleanalytics

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "google",
			Label:     `Google API`,
			SortOrder: 340,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_GoogleAnalytics::google
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "analytics",
					Label:     `Google Analytics`,
					SortOrder: 10,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: google/analytics/active
							ID:        "active",
							Label:     `Enable`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: google/analytics/account
							ID:        "account",
							Label:     `Account Number`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
