package salesrule

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "promo",
		Label:     "Promotions",
		SortOrder: 400,
		Scope:     config.NewScopePerm(config.ScopeDefault),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "auto_generated_coupon_codes",
				Label:     `Auto Generated Specific Coupon Codes`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/length`,
						ID:           "length",
						Label:        `Code Length`,
						Comment:      `Excluding prefix, suffix and separators.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      12,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/format`,
						ID:           "format",
						Label:        `Code Format`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\SalesRule\Model\System\Config\Source\Coupon\Format
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/prefix`,
						ID:           "prefix",
						Label:        `Code Prefix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/suffix`,
						ID:           "suffix",
						Label:        `Code Suffix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `promo/auto_generated_coupon_codes/dash`,
						ID:           "dash",
						Label:        `Dash Every X Characters`,
						Comment:      `If empty no separation.`,
						Type:         config.TypeText,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "rss",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "catalog",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/catalog/discounts`,
						ID:           "discounts",
						Label:        `Coupons/Discounts`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    12,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},
)
