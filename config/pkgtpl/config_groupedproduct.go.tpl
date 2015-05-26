package groupedproduct

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "checkout",
		Label:     "",
		SortOrder: 305,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "cart",
				Label:     ``,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `checkout/cart/grouped_product_image`,
						ID:           "grouped_product_image",
						Label:        `Grouped Product Image`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `itself`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Product\Thumbnail
					},
				},
			},
		},
	},
)
