package search

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "search",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/search/engine`,
						ID:           "engine",
						Label:        `Search Engine`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    19,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Search\Model\Adminhtml\System\Config\Source\Engine
					},

					&config.Field{
						// Path: `catalog/search/search_type`,
						ID:           "search_type",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      true,
						Scope:        config.NewScopePerm(),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/search/use_layered_navigation_count`,
						ID:           "use_layered_navigation_count",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    0,
						Visible:      true,
						Scope:        config.NewScopePerm(),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
