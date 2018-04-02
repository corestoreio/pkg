// +build ignore

package groupedproduct

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID:        "checkout",
			SortOrder: 305,
			Scopes:    scope.PermStore,
			Groups: element.MakeGroups(
				element.Group{
					ID:        "cart",
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: checkout/cart/grouped_product_image
							ID:        "grouped_product_image",
							Label:     `Grouped Product Image`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `itself`,
							// SourceModel: Magento\Catalog\Model\Config\Source\Product\Thumbnail
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
