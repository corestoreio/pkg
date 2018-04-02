// +build ignore

package configurableproduct

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
			ID: "checkout",
			Groups: element.MakeGroups(
				element.Group{
					ID: "cart",
					Fields: element.MakeFields(
						element.Field{
							// Path: checkout/cart/configurable_product_image
							ID:        "configurable_product_image",
							Label:     `Configurable Product Image`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `parent`,
							// SourceModel: Magento\Catalog\Model\Config\Source\Product\Thumbnail
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
