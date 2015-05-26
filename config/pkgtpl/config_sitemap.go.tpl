package sitemap

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sitemap",
		Label:     "XML Sitemap",
		SortOrder: 70,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "category",
				Label:     `Categories Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/category/changefreq`,
						ID:           "changefreq",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `daily`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `sitemap/category/priority`,
						ID:           "priority",
						Label:        `Priority`,
						Comment:      `Valid values range from 0.0 to 1.0.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      0.5,
						BackendModel: nil, // Magento\Sitemap\Model\Config\Backend\Priority
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "product",
				Label:     `Products Options`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/product/changefreq`,
						ID:           "changefreq",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `daily`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `sitemap/product/priority`,
						ID:           "priority",
						Label:        `Priority`,
						Comment:      `Valid values range from 0.0 to 1.0.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      1,
						BackendModel: nil, // Magento\Sitemap\Model\Config\Backend\Priority
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sitemap/product/image_include`,
						ID:           "image_include",
						Label:        `Add Images into Sitemap`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `all`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sitemap\Model\Source\Product\Image\IncludeImage
					},
				},
			},

			&config.Group{
				ID:        "page",
				Label:     `CMS Pages Options`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/page/changefreq`,
						ID:           "changefreq",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `daily`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `sitemap/page/priority`,
						ID:           "priority",
						Label:        `Priority`,
						Comment:      `Valid values range from 0.0 to 1.0.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      0.25,
						BackendModel: nil, // Magento\Sitemap\Model\Config\Backend\Priority
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "generate",
				Label:     `Generation Settings`,
				Comment:   ``,
				SortOrder: 4,
				Scope:     config.NewScopePerm(config.IDScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/generate/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sitemap/generate/error_email`,
						ID:           "error_email",
						Label:        `Error Email Recipient`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sitemap/generate/error_email_identity`,
						ID:           "error_email_identity",
						Label:        `Error Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `sitemap/generate/error_email_template`,
						ID:           "error_email_template",
						Label:        `Error Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      `sitemap_generate_error_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sitemap/generate/frequency`,
						ID:           "frequency",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Cron\Model\Config\Backend\Sitemap
						SourceModel:  nil, // Magento\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `sitemap/generate/time`,
						ID:           "time",
						Label:        `Start Time`,
						Comment:      ``,
						Type:         config.TypeTime,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "limit",
				Label:     `Sitemap File Limits`,
				Comment:   ``,
				SortOrder: 5,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/limit/max_lines`,
						ID:           "max_lines",
						Label:        `Maximum No of URLs Per File`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      50000,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sitemap/limit/max_file_size`,
						ID:           "max_file_size",
						Label:        `Maximum File Size`,
						Comment:      `File size in bytes.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      10485760,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "search_engines",
				Label:     `Search Engine Submission Settings`,
				Comment:   ``,
				SortOrder: 6,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/search_engines/submission_robots`,
						ID:           "submission_robots",
						Label:        `Enable Submission to Robots.txt`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "sitemap",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "generate",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/generate/error_email`,
						ID:      "error_email",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: nil,
					},
				},
			},

			&config.Group{
				ID: "file",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sitemap/file/valid_paths`,
						ID:      "valid_paths",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"available":{"any_path":"\/*\/*.xml"}}`,
					},
				},
			},
		},
	},
)
