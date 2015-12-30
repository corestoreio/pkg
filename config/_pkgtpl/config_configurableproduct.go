// +build ignore

package configurableproduct

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
					ID: "cart",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: checkout/cart/configurable_product_image
							ID:        "configurable_product_image",
							Label:     `Configurable Product Image`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `parent`,
							// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
						},
					),
				},
			),
		},
	)
}
