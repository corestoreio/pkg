// +build ignore

package configurableproduct

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "checkout",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "cart",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/cart/configurable_product_image
						ID:        "configurable_product_image",
						Label:     `Configurable Product Image`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `parent`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
					},
				),
			},
		),
	},
)
