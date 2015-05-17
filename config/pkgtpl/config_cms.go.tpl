package cms

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "web",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "default",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     config.NewScopePerm(),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/default/cms_home_page`,
						ID:           "cms_home_page",
						Label:        `CMS Home Page`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `home`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: `web/default/cms_no_route`,
						ID:           "cms_no_route",
						Label:        `CMS No Route Page`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `no-route`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: `web/default/cms_no_cookies`,
						ID:           "cms_no_cookies",
						Label:        `CMS No Cookies Page`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `enable-cookies`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Cms\Model\Config\Source\Page
					},

					&config.Field{
						// Path: `web/default/show_cms_breadcrumbs`,
						ID:           "show_cms_breadcrumbs",
						Label:        `Show Breadcrumbs for CMS Pages`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "browser_capabilities",
				Label:     `Browser Capabilities Detection`,
				Comment:   ``,
				SortOrder: 200,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/browser_capabilities/cookies`,
						ID:           "cookies",
						Label:        `Redirect to CMS-page if Cookies are Disabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/browser_capabilities/javascript`,
						ID:           "javascript",
						Label:        `Show Notice if JavaScript is Disabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    200,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/browser_capabilities/local_storage`,
						ID:           "local_storage",
						Label:        `Show Notice if Local Storage is Disabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    300,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "cms",
		Label:     "Content Management",
		SortOrder: 1001,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "wysiwyg",
				Label:     `WYSIWYG Options`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `cms/wysiwyg/enabled`,
						ID:           "enabled",
						Label:        `Enable WYSIWYG Editor`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      `enabled`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Cms\Model\Config\Source\Wysiwyg\Enabled
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "web",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "default",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/default/front`,
						ID:      "front",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `cms`,
					},

					&config.Field{
						// Path: `web/default/no_route`,
						ID:      "no_route",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `cms/noroute/index`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/media_storage_configuration/allowed_resources`,
						ID:      "allowed_resources",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"wysiwyg_image_folder":"wysiwyg"}`,
					},
				},
			},
		},
	},
)
