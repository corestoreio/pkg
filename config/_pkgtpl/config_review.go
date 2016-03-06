// +build ignore

package review

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
					ID:        "review",
					Label:     `Product Reviews`,
					SortOrder: 100,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/review/allow_guest
							ID:        "allow_guest",
							Label:     `Allow Guests to Write Reviews`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
