package catalogurlrewrite

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "seo",
				Label:     `Search Engine Optimization`,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/seo/category_url_suffix`,
						ID:           "category_url_suffix",
						Label:        `Category URL Suffix`,
						Comment:      `You need to refresh the cache.`,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/seo/product_url_suffix`,
						ID:           "product_url_suffix",
						Label:        `Product URL Suffix`,
						Comment:      `You need to refresh the cache.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/seo/product_use_categories`,
						ID:           "product_use_categories",
						Label:        `Use Categories Path for Product URLs`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/seo/save_rewrites_history`,
						ID:           "save_rewrites_history",
						Label:        `Create Permanent Redirect for URLs if URL Key Changed`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
