package googleanalytics

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "google",
		Label:     "Google API",
		SortOrder: 340,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "analytics",
				Label:     `Google Analytics`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/analytics/active`,
						ID:           "active",
						Label:        `Enable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `google/analytics/account`,
						ID:           "account",
						Label:        `Account Number`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
