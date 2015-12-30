// +build ignore

package groupedproduct

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		SortOrder: 305,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "cart",
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: checkout/cart/grouped_product_image
						ID:        "grouped_product_image",
						Label:     `Grouped Product Image`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `itself`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\Product\Thumbnail
					},
				),
			},
		),
	},
)
