package backend

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "advanced",
		Label:     "Advanced",
		SortOrder: 910,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "modules_disable_output",
				Label:     `Disable Modules Output`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields:    config.FieldSlice{},
			},
		},
	},
	&config.Section{
		ID:        "trans_email",
		Label:     "Store Email Addresses",
		SortOrder: 90,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "ident_custom1",
				Label:     `Custom Email 1`,
				Comment:   ``,
				SortOrder: 4,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `trans_email/ident_custom1/email`,
						ID:           "email",
						Label:        `Sender Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `trans_email/ident_custom1/name`,
						ID:           "name",
						Label:        `Sender Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "ident_custom2",
				Label:     `Custom Email 2`,
				Comment:   ``,
				SortOrder: 5,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `trans_email/ident_custom2/email`,
						ID:           "email",
						Label:        `Sender Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `trans_email/ident_custom2/name`,
						ID:           "name",
						Label:        `Sender Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "ident_general",
				Label:     `General Contact`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `trans_email/ident_general/email`,
						ID:           "email",
						Label:        `Sender Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `trans_email/ident_general/name`,
						ID:           "name",
						Label:        `Sender Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "ident_sales",
				Label:     `Sales Representative`,
				Comment:   ``,
				SortOrder: 2,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `trans_email/ident_sales/email`,
						ID:           "email",
						Label:        `Sender Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `trans_email/ident_sales/name`,
						ID:           "name",
						Label:        `Sender Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "ident_support",
				Label:     `Customer Support`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `trans_email/ident_support/email`,
						ID:           "email",
						Label:        `Sender Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `trans_email/ident_support/name`,
						ID:           "name",
						Label:        `Sender Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "design",
		Label:     "Design",
		SortOrder: 30,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "theme",
				Label:     `Design Theme`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/theme/theme_id`,
						ID:           "theme_id",
						Label:        `Design Theme`,
						Comment:      `If no value is specified, the system default will be used. The system default may be modified by third party extensions.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Theme\Model\Design\Backend\Theme
						SourceModel:  nil, // Magento\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
					},

					&config.Field{
						// Path: `design/theme/ua_regexp`,
						ID:           "ua_regexp",
						Label:        `User-Agent Exceptions`,
						Comment:      `Search strings are either normal strings or regular exceptions (PCRE). They are matched in the same order as entered. Examples:<br /><span style="font-family:monospace">Firefox<br />/^mozilla/i</span>`,
						Type:         config.Type,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Theme\Model\Design\Backend\Exceptions
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "pagination",
				Label:     `Pagination`,
				Comment:   ``,
				SortOrder: 500,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/pagination/pagination_frame`,
						ID:           "pagination_frame",
						Label:        `Pagination Frame`,
						Comment:      `How many links to display at once.`,
						Type:         config.TypeText,
						SortOrder:    7,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/pagination/pagination_frame_skip`,
						ID:           "pagination_frame_skip",
						Label:        `Pagination Frame Skip`,
						Comment:      `If the current frame position does not cover utmost pages, will render link to current position plus/minus this value.`,
						Type:         config.TypeText,
						SortOrder:    8,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/pagination/anchor_text_for_previous`,
						ID:           "anchor_text_for_previous",
						Label:        `Anchor Text for Previous`,
						Comment:      `Alternative text for previous link in pagination menu. If empty, default arrow image will used.`,
						Type:         config.TypeText,
						SortOrder:    9,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/pagination/anchor_text_for_next`,
						ID:           "anchor_text_for_next",
						Label:        `Anchor Text for Next`,
						Comment:      `Alternative text for next link in pagination menu. If empty, default arrow image will used.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "email",
				Label:     `Transactional Emails`,
				Comment:   ``,
				SortOrder: 510,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/email/logo`,
						ID:           "logo",
						Label:        `Logo Image`,
						Comment:      `Allowed file types: jpg, jpeg, gif, png`,
						Type:         config.TypeImage,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Logo
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/email/logo_alt`,
						ID:           "logo_alt",
						Label:        `Logo Image Alt`,
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
	&config.Section{
		ID:        "dev",
		Label:     "Developer",
		SortOrder: 920,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "debug",
				Label:     `Debug`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.NewScopePerm(config.ScopeWebsite, config.ScopeStore),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/debug/template_hints`,
						ID:           "template_hints",
						Label:        `Template Path Hints`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeWebsite, config.ScopeStore),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `dev/debug/template_hints_blocks`,
						ID:           "template_hints_blocks",
						Label:        `Add Block Names to Hints`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    21,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeWebsite, config.ScopeStore),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "template",
				Label:     `Template Settings`,
				Comment:   ``,
				SortOrder: 25,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/template/allow_symlink`,
						ID:           "allow_symlink",
						Label:        `Allow Symlinks`,
						Comment:      `Warning! Enabling this feature is not recommended on production environments because it represents a potential security risk.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `dev/template/minify_html`,
						ID:           "minify_html",
						Label:        `Minify Html`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    25,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "translate_inline",
				Label:     `Translate Inline`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/translate_inline/active`,
						ID:           "active",
						Label:        `Enabled for Frontend`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Translate
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `dev/translate_inline/active_admin`,
						ID:           "active_admin",
						Label:        `Enabled for Admin`,
						Comment:      `Translate, blocks and other output caches should be disabled for both frontend and admin inline translations.`,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Translate
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "js",
				Label:     `JavaScript Settings`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/js/merge_files`,
						ID:           "merge_files",
						Label:        `Merge JavaScript Files`,
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
						// Path: `dev/js/enable_js_bundling`,
						ID:           "enable_js_bundling",
						Label:        `Enable Javascript Bundling`,
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
						// Path: `dev/js/minify_files`,
						ID:           "minify_files",
						Label:        `Minify JavaScript Files`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "css",
				Label:     `CSS Settings`,
				Comment:   ``,
				SortOrder: 110,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/css/merge_css_files`,
						ID:           "merge_css_files",
						Label:        `Merge CSS Files`,
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
						// Path: `dev/css/minify_files`,
						ID:           "minify_files",
						Label:        `Minify CSS Files`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "image",
				Label:     `Image Processing Settings`,
				Comment:   ``,
				SortOrder: 120,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/image/default_adapter`,
						ID:           "default_adapter",
						Label:        `Image Adapter`,
						Comment:      `When the adapter was changed, please, flush Catalog Images Cache.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image\Adapter
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Image\Adapter
					},
				},
			},

			&config.Group{
				ID:        "static",
				Label:     `Static Files Settings`,
				Comment:   ``,
				SortOrder: 130,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `dev/static/sign`,
						ID:           "sign",
						Label:        `Sign Static Files`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "general",
		Label:     "General",
		SortOrder: 10,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "country",
				Label:     `Country Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/country/allow`,
						ID:           "allow",
						Label:        `Allow Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `general/country/default`,
						ID:           "default",
						Label:        `Default Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `general/country/eu_countries`,
						ID:           "eu_countries",
						Label:        `European Union Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},
				},
			},

			&config.Group{
				ID:        "locale",
				Label:     `Locale Options`,
				Comment:   ``,
				SortOrder: 8,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/locale/timezone`,
						ID:           "timezone",
						Label:        `Timezone`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Locale\Timezone
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Timezone
					},

					&config.Field{
						// Path: `general/locale/code`,
						ID:           "code",
						Label:        `Locale`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale
					},

					&config.Field{
						// Path: `general/locale/firstday`,
						ID:           "firstday",
						Label:        `First Day of Week`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Weekdays
					},

					&config.Field{
						// Path: `general/locale/weekend`,
						ID:           "weekend",
						Label:        `Weekend Days`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    15,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Weekdays
					},
				},
			},

			&config.Group{
				ID:        "store_information",
				Label:     `Store Information`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/store_information/name`,
						ID:           "name",
						Label:        `Store Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/phone`,
						ID:           "phone",
						Label:        `Store Phone Number`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/country_id`,
						ID:           "country_id",
						Label:        `Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    25,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `general/store_information/region_id`,
						ID:           "region_id",
						Label:        `Region/State`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    27,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/postcode`,
						ID:           "postcode",
						Label:        `ZIP/Postal Code`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/city`,
						ID:           "city",
						Label:        `City`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    45,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/street_line1`,
						ID:           "street_line1",
						Label:        `Street Address`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    55,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/street_line2`,
						ID:           "street_line2",
						Label:        `Street Address Line 2`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `general/store_information/merchant_vat_number`,
						ID:           "merchant_vat_number",
						Label:        `VAT Number`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    61,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "single_store_mode",
				Label:     `Single-Store Mode`,
				Comment:   ``,
				SortOrder: 150,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/single_store_mode/enabled`,
						ID:           "enabled",
						Label:        `Enable Single-Store Mode`,
						Comment:      `This setting will not be taken into account if system has more than one store view.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "system",
		Label:     "System",
		SortOrder: 900,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "smtp",
				Label:     `Mail Sending Settings`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/smtp/disable`,
						ID:           "disable",
						Label:        `Disable Email Communications`,
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
						// Path: `system/smtp/host`,
						ID:           "host",
						Label:        `Host`,
						Comment:      `For Windows server only.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/smtp/port`,
						ID:           "port",
						Label:        `Port (25)`,
						Comment:      `For Windows server only.`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/smtp/set_return_path`,
						ID:           "set_return_path",
						Label:        `Set Return-Path`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesnocustom
					},

					&config.Field{
						// Path: `system/smtp/return_path_email`,
						ID:           "return_path_email",
						Label:        `Return-Path Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    80,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "admin",
		Label:     "Admin",
		SortOrder: 20,
		Scope:     config.NewScopePerm(config.ScopeDefault),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "emails",
				Label:     `Admin User Emails`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/emails/forgot_email_template`,
						ID:           "forgot_email_template",
						Label:        `Forgot Password Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `admin/emails/forgot_email_identity`,
						ID:           "forgot_email_identity",
						Label:        `Forgot and Reset Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `admin/emails/password_reset_link_expiration_period`,
						ID:           "password_reset_link_expiration_period",
						Label:        `Recovery Link Expiration Period (days)`,
						Comment:      `Please enter a number 1 or greater in this field.`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "startup",
				Label:     `Startup Page`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/startup/menu_item_id`,
						ID:           "menu_item_id",
						Label:        `Startup Page`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Admin\Page
					},
				},
			},

			&config.Group{
				ID:        "url",
				Label:     `Admin Base URL`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/url/use_custom`,
						ID:           "use_custom",
						Label:        `Use Custom Admin URL`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Usecustom
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `admin/url/custom`,
						ID:           "custom",
						Label:        `Custom Admin URL`,
						Comment:      `Make sure that base URL ends with '/' (slash), e.g. http://yourdomain/magento/`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Custom
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `admin/url/use_custom_path`,
						ID:           "use_custom_path",
						Label:        `Use Custom Admin Path`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Custompath
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `admin/url/custom_path`,
						ID:           "custom_path",
						Label:        `Custom Admin Path`,
						Comment:      `You will have to log in after you save your custom admin path.`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Custompath
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "security",
				Label:     `Security`,
				Comment:   ``,
				SortOrder: 35,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/security/use_form_key`,
						ID:           "use_form_key",
						Label:        `Add Secret Key to URLs`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Usesecretkey
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `admin/security/use_case_sensitive_login`,
						ID:           "use_case_sensitive_login",
						Label:        `Login is Case Sensitive`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `admin/security/session_lifetime`,
						ID:           "session_lifetime",
						Label:        `Admin Session Lifetime (seconds)`,
						Comment:      `Values less than 60 are ignored.`,
						Type:         config.Type,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "dashboard",
				Label:     `Dashboard`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/dashboard/enable_charts`,
						ID:           "enable_charts",
						Label:        `Enable Charts`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "web",
		Label:     "Web",
		SortOrder: 20,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "url",
				Label:     `Url Options`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     config.NewScopePerm(config.ScopeDefault),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/url/use_store`,
						ID:           "use_store",
						Label:        `Add Store Code to Urls`,
						Comment:      `<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Store
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/url/redirect_to_base`,
						ID:           "redirect_to_base",
						Label:        `Auto-redirect to Base URL`,
						Comment:      `I.e. redirect from http://example.com/store/ to http://www.example.com/store/`,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Web\Redirect
					},
				},
			},

			&config.Group{
				ID:        "seo",
				Label:     `Search Engine Optimization`,
				Comment:   ``,
				SortOrder: 5,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/seo/use_rewrites`,
						ID:           "use_rewrites",
						Label:        `Use Web Server Rewrites`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "unsecure",
				Label:     `Base URLs`,
				Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/unsecure/base_url`,
						ID:           "base_url",
						Label:        `Base URL`,
						Comment:      `Specify URL or {{base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/unsecure/base_link_url`,
						ID:           "base_link_url",
						Label:        `Base Link URL`,
						Comment:      `May start with {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/unsecure/base_static_url`,
						ID:           "base_static_url",
						Label:        `Base URL for Static View Files`,
						Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    25,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/unsecure/base_media_url`,
						ID:           "base_media_url",
						Label:        `Base URL for User Media Files`,
						Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "secure",
				Label:     `Base URLs (Secure)`,
				Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`,
				SortOrder: 20,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/secure/base_url`,
						ID:           "base_url",
						Label:        `Secure Base URL`,
						Comment:      `Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/secure/base_link_url`,
						ID:           "base_link_url",
						Label:        `Secure Base Link URL`,
						Comment:      `May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/secure/base_static_url`,
						ID:           "base_static_url",
						Label:        `Secure Base URL for Static View Files`,
						Comment:      `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    25,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/secure/base_media_url`,
						ID:           "base_media_url",
						Label:        `Secure Base URL for User Media Files`,
						Comment:      `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/secure/use_in_frontend`,
						ID:           "use_in_frontend",
						Label:        `Use Secure URLs in Frontend`,
						Comment:      `Enter https protocol to use Secure URLs in Frontend.`,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Secure
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/secure/use_in_adminhtml`,
						ID:           "use_in_adminhtml",
						Label:        `Use Secure URLs in Admin`,
						Comment:      `Enter https protocol to use Secure URLs in Admin.`,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Secure
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/secure/offloader_header`,
						ID:           "offloader_header",
						Label:        `Offloader header`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    70,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "default",
				Label:     `Default Pages`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/default/front`,
						ID:           "front",
						Label:        `Default Web URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `web/default/no_route`,
						ID:           "no_route",
						Label:        `Default No-route URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.ScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "session",
				Label:     `Session Validation Settings`,
				Comment:   ``,
				SortOrder: 60,
				Scope:     config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `web/session/use_remote_addr`,
						ID:           "use_remote_addr",
						Label:        `Validate REMOTE_ADDR`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/session/use_http_via`,
						ID:           "use_http_via",
						Label:        `Validate HTTP_VIA`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/session/use_http_x_forwarded_for`,
						ID:           "use_http_x_forwarded_for",
						Label:        `Validate HTTP_X_FORWARDED_FOR`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/session/use_http_user_agent`,
						ID:           "use_http_user_agent",
						Label:        `Validate HTTP_USER_AGENT`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `web/session/use_frontend_sid`,
						ID:           "use_frontend_sid",
						Label:        `Use SID on Frontend`,
						Comment:      `Allows customers to stay logged in when switching between different stores.`,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
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
						Default: `{"email_folder":"email"}`,
					},
				},
			},

			&config.Group{
				ID: "emails",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/emails/forgot_email_template`,
						ID:      "forgot_email_template",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `system_emails_forgot_email_template`,
					},

					&config.Field{
						// Path: `system/emails/forgot_email_identity`,
						ID:      "forgot_email_identity",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `general`,
					},
				},
			},

			&config.Group{
				ID: "dashboard",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/dashboard/enable_charts`,
						ID:      "enable_charts",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: true,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "general",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "validator_data",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/validator_data/input_types`,
						ID:      "input_types",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"price":"price","media_image":"media_image","gallery":"gallery"}`,
					},
				},
			},
		},
	},
)
