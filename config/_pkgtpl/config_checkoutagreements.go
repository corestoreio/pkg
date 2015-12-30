// +build ignore

package checkoutagreements

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
			ID: "checkout",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "options",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: checkout/options/enable_agreements
							ID:        "enable_agreements",
							Label:     `Enable Terms and Conditions`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
}
