// +build ignore

package configurableproduct

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID: "checkout",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "cart",
					Fields: element.NewFieldSlice(
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
