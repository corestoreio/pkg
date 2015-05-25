package multishipping

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "multishipping",
		Label:     "Multishipping Settings",
		SortOrder: 311,
		Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     `Options`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `multishipping/options/checkout_multiple`,
						ID:           "checkout_multiple",
						Label:        `Allow Shipping to Multiple Addresses`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `multishipping/options/checkout_multiple_maximum_qty`,
						ID:           "checkout_multiple_maximum_qty",
						Label:        `Maximum Qty Allowed for Shipping to Multiple Addresses`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      100,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
