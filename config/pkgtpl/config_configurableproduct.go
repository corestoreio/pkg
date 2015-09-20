// +build ignore

package configurableproduct

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cart",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/cart/configurable_product_image`,
						ID:           "configurable_product_image",
						Label:        `Configurable Product Image`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `parent`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Product\Thumbnail
					},
				},
			},
		},
	},
)
