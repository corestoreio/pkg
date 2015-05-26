package googleoptimizer

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "google",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "analytics",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/analytics/experiments`,
						ID:           "experiments",
						Label:        `Enable Content Experiments`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "google",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "optimizer",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/optimizer/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: false,
					},
				},
			},
		},
	},
)
