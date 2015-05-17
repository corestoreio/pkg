package captcha

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "admin",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "captcha",
				Label:     `CAPTCHA`,
				Comment:   ``,
				SortOrder: 50,
				Scope:     config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/captcha/enable`,
						ID:           "enable",
						Label:        `Enable CAPTCHA in Admin`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `admin/captcha/font`,
						ID:           "font",
						Label:        `Font`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      `linlibertine`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Font
					},

					&config.Field{
						// Path: `admin/captcha/forms`,
						ID:           "forms",
						Label:        `Forms`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      `backend_forgotpassword`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Form\Backend
					},

					&config.Field{
						// Path: `admin/captcha/mode`,
						ID:           "mode",
						Label:        `Displaying Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      `after_fail`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Mode
					},

					&config.Field{
						// Path: `admin/captcha/failed_attempts_login`,
						ID:           "failed_attempts_login",
						Label:        `Number of Unsuccessful Attempts to Login`,
						Comment:      `If 0 is specified, CAPTCHA on the Login form will be always available.`,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      3,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `admin/captcha/timeout`,
						ID:           "timeout",
						Label:        `CAPTCHA Timeout (minutes)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    6,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      7,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `admin/captcha/length`,
						ID:           "length",
						Label:        `Number of Symbols`,
						Comment:      `Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`,
						Type:         config.TypeText,
						SortOrder:    7,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      `4-5`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `admin/captcha/symbols`,
						ID:           "symbols",
						Label:        `Symbols Used in CAPTCHA`,
						Comment:      `Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`,
						Type:         config.TypeText,
						SortOrder:    8,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      `ABCDEFGHJKMnpqrstuvwxyz23456789`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `admin/captcha/case_sensitive`,
						ID:           "case_sensitive",
						Label:        `Case Sensitive`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    9,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "customer",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "captcha",
				Label:     `CAPTCHA`,
				Comment:   ``,
				SortOrder: 110,
				Scope:     config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/captcha/enable`,
						ID:           "enable",
						Label:        `Enable CAPTCHA on Frontend`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/captcha/font`,
						ID:           "font",
						Label:        `Font`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `linlibertine`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Font
					},

					&config.Field{
						// Path: `customer/captcha/forms`,
						ID:           "forms",
						Label:        `Forms`,
						Comment:      `CAPTCHA for "Create user" and "Forgot password" forms is always enabled if chosen.`,
						Type:         config.TypeMultiselect,
						SortOrder:    3,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `user_forgotpassword`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Form\Frontend
					},

					&config.Field{
						// Path: `customer/captcha/mode`,
						ID:           "mode",
						Label:        `Displaying Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `after_fail`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Captcha\Model\Config\Mode
					},

					&config.Field{
						// Path: `customer/captcha/failed_attempts_login`,
						ID:           "failed_attempts_login",
						Label:        `Number of Unsuccessful Attempts to Login`,
						Comment:      `If 0 is specified, CAPTCHA on the Login form will be always available.`,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      3,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/captcha/timeout`,
						ID:           "timeout",
						Label:        `CAPTCHA Timeout (minutes)`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    6,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      7,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/captcha/length`,
						ID:           "length",
						Label:        `Number of Symbols`,
						Comment:      `Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`,
						Type:         config.TypeText,
						SortOrder:    7,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `4-5`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/captcha/symbols`,
						ID:           "symbols",
						Label:        `Symbols Used in CAPTCHA`,
						Comment:      `Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`,
						Type:         config.TypeText,
						SortOrder:    8,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
						Default:      `ABCDEFGHJKMnpqrstuvwxyz23456789`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/captcha/case_sensitive`,
						ID:           "case_sensitive",
						Label:        `Case Sensitive`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    9,
						Visible:      true,
						Scope:        config.NewScopePerm(config.ScopeDefault, config.ScopeWebsite),
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
						Default: `{"captcha_folder":"captcha"}`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "admin",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "captcha",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `admin/captcha/type`,
						ID:      "type",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `default`,
					},

					&config.Field{
						// Path: `admin/captcha/failed_attempts_ip`,
						ID:      "failed_attempts_ip",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: 1000,
					},

					&config.Field{
						// Path: `admin/captcha/shown_to_logged_in_user`,
						ID:      "shown_to_logged_in_user",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `admin/captcha/always_for`,
						ID:      "always_for",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"backend_forgotpassword":"1"}`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "customer",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "captcha",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/captcha/type`,
						ID:      "type",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `default`,
					},

					&config.Field{
						// Path: `customer/captcha/failed_attempts_ip`,
						ID:      "failed_attempts_ip",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: 1000,
					},

					&config.Field{
						// Path: `customer/captcha/shown_to_logged_in_user`,
						ID:      "shown_to_logged_in_user",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"contact_us":"1"}`,
					},

					&config.Field{
						// Path: `customer/captcha/always_for`,
						ID:      "always_for",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"user_create":"1","user_forgotpassword":"1","guest_checkout":"1","register_during_checkout":"1","contact_us":"1"}`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "captcha",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "_value",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `captcha/_value/fonts`,
						ID:      "fonts",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"linlibertine":{"label":"LinLibertine","path":"LinLibertineFont\/LinLibertine_Bd-2.8.1.ttf"}}`,
					},

					&config.Field{
						// Path: `captcha/_value/frontend`,
						ID:      "frontend",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"areas":{"user_create":{"label":"Create user"},"user_login":{"label":"Login"},"user_forgotpassword":{"label":"Forgot password"},"guest_checkout":{"label":"Checkout as Guest"},"register_during_checkout":{"label":"Register during Checkout"},"contact_us":{"label":"Contact Us"}}}`,
					},

					&config.Field{
						// Path: `captcha/_value/backend`,
						ID:      "backend",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `{"areas":{"backend_login":{"label":"Admin Login"},"backend_forgotpassword":{"label":"Admin Forgot Password"}}}`,
					},
				},
			},

			&config.Group{
				ID: "_attribute",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `captcha/_attribute/translate`,
						ID:      "translate",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.ScopeDefault), // @todo search for that
						Default: `label`,
					},
				},
			},
		},
	},
)
