// +build ignore

package groupedproduct

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     "",
		SortOrder: 305,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cart",
				Label:     ``,
				Comment:   ``,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/cart/grouped_product_image`,
						ID:           "grouped_product_image",
						Label:        `Grouped Product Image`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `itself`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Product\Thumbnail
					},
				},
			},
		},
	},
)
