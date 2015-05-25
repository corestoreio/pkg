package translation

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "dev",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "js",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/js/translate_strategy`,
						ID:           "translate_strategy",
						Label:        `Translation Strategy`,
						Comment:      `Please put your store into maintenance mode and redeploy static files after changing strategy`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault),
						Default:      `none`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Translation\Model\Js\Config\Source\Strategy
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "dev",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "translate_inline",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/translate_inline/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `dev/translate_inline/active_admin`,
						ID:      "active_admin",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `dev/translate_inline/invalid_caches`,
						ID:      "invalid_caches",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"block_html":null}`,
					},
				},
			},
		},
	},
)
