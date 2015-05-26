package webapi

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "webapi",
		Label:     "Magento Web API",
		SortOrder: 102,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "soap",
				Label:     `SOAP Settings`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `webapi/soap/charset`,
						ID:           "charset",
						Label:        `Default Response Charset`,
						Comment:      `If empty, UTF-8 will be used.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
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
