package wishlist

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "wishlist",
		Label:     "Wishlist",
		SortOrder: 140,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "email",
				Label:     `Share Options`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `wishlist/email/email_identity`,
						ID:           "email_identity",
						Label:        `Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `wishlist/email/email_template`,
						ID:           "email_template",
						Label:        `Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `wishlist_email_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `wishlist/email/number_limit`,
						ID:           "number_limit",
						Label:        `Max Emails Allowed to be Sent`,
						Comment:      `10 by default. Max - 10000`,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      10,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `wishlist/email/text_limit`,
						ID:           "text_limit",
						Label:        `Email Text Length Limit`,
						Comment:      `255 by default`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      255,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "general",
				Label:     `General Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `wishlist/general/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "wishlist_link",
				Label:     `My Wish List Link`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `wishlist/wishlist_link/use_qty`,
						ID:           "use_qty",
						Label:        `Display Wish List Summary`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Wishlist\Model\Config\Source\Summary
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
				ID:        "wishlist",
				Label:     `Wishlist`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/wishlist/active`,
						ID:           "active",
						Label:        `Enable RSS`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
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
