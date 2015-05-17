package payment

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "payment",
		Label:     "Payment Methods",
		SortOrder: 400,
		Scope:     config.ScopePermAll,
		Groups:    config.GroupSlice{},
	},

	// Hidden Configuration
	&config.Section{
		ID: "payment",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "free",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/free/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: true,
					},

					&config.Field{
						// Path: `payment/free/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `Magento\Payment\Model\Method\Free`,
					},

					&config.Field{
						// Path: `payment/free/order_status`,
						ID:      "order_status",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `pending`,
					},

					&config.Field{
						// Path: `payment/free/title`,
						ID:      "title",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `No Payment Information Required`,
					},

					&config.Field{
						// Path: `payment/free/allowspecific`,
						ID:      "allowspecific",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `payment/free/sort_order`,
						ID:      "sort_order",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: true,
					},
				},
			},

			&config.Group{
				ID: "substitution",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/substitution/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `payment/substitution/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `Magento\Payment\Model\Method\Substitution`,
					},

					&config.Field{
						// Path: `payment/substitution/allowspecific`,
						ID:      "allowspecific",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: false,
					},
				},
			},
		},
	},
)
