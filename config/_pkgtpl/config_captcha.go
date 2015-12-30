// +build ignore

package captcha

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "admin",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "captcha",
				Label:     `CAPTCHA`,
				SortOrder: 50,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/captcha/enable
						ID:        "enable",
						Label:     `Enable CAPTCHA in Admin`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: admin/captcha/font
						ID:        "font",
						Label:     `Font`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `linlibertine`,
						// SourceModel: Otnegam\Captcha\Model\Config\Font
					},

					&config.Field{
						// Path: admin/captcha/forms
						ID:        "forms",
						Label:     `Forms`,
						Type:      config.TypeMultiselect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `backend_forgotpassword`,
						// SourceModel: Otnegam\Captcha\Model\Config\Form\Backend
					},

					&config.Field{
						// Path: admin/captcha/mode
						ID:        "mode",
						Label:     `Displaying Mode`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `after_fail`,
						// SourceModel: Otnegam\Captcha\Model\Config\Mode
					},

					&config.Field{
						// Path: admin/captcha/failed_attempts_login
						ID:        "failed_attempts_login",
						Label:     `Number of Unsuccessful Attempts to Login`,
						Comment:   element.LongText(`If 0 is specified, CAPTCHA on the Login form will be always available.`),
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   3,
					},

					&config.Field{
						// Path: admin/captcha/timeout
						ID:        "timeout",
						Label:     `CAPTCHA Timeout (minutes)`,
						Type:      config.TypeText,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   7,
					},

					&config.Field{
						// Path: admin/captcha/length
						ID:        "length",
						Label:     `Number of Symbols`,
						Comment:   element.LongText(`Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`),
						Type:      config.TypeText,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `4-5`,
					},

					&config.Field{
						// Path: admin/captcha/symbols
						ID:        "symbols",
						Label:     `Symbols Used in CAPTCHA`,
						Comment:   element.LongText(`Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`),
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `ABCDEFGHJKMnpqrstuvwxyz23456789`,
					},

					&config.Field{
						// Path: admin/captcha/case_sensitive
						ID:        "case_sensitive",
						Label:     `Case Sensitive`,
						Type:      config.TypeSelect,
						SortOrder: 9,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID: "customer",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "captcha",
				Label:     `CAPTCHA`,
				SortOrder: 110,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: customer/captcha/enable
						ID:        "enable",
						Label:     `Enable CAPTCHA on Storefront`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: customer/captcha/font
						ID:        "font",
						Label:     `Font`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `linlibertine`,
						// SourceModel: Otnegam\Captcha\Model\Config\Font
					},

					&config.Field{
						// Path: customer/captcha/forms
						ID:        "forms",
						Label:     `Forms`,
						Comment:   element.LongText(`CAPTCHA for "Create user" and "Forgot password" forms is always enabled if chosen.`),
						Type:      config.TypeMultiselect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `user_forgotpassword`,
						// SourceModel: Otnegam\Captcha\Model\Config\Form\Frontend
					},

					&config.Field{
						// Path: customer/captcha/mode
						ID:        "mode",
						Label:     `Displaying Mode`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `after_fail`,
						// SourceModel: Otnegam\Captcha\Model\Config\Mode
					},

					&config.Field{
						// Path: customer/captcha/failed_attempts_login
						ID:        "failed_attempts_login",
						Label:     `Number of Unsuccessful Attempts to Login`,
						Comment:   element.LongText(`If 0 is specified, CAPTCHA on the Login form will be always available.`),
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   3,
					},

					&config.Field{
						// Path: customer/captcha/timeout
						ID:        "timeout",
						Label:     `CAPTCHA Timeout (minutes)`,
						Type:      config.TypeText,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   7,
					},

					&config.Field{
						// Path: customer/captcha/length
						ID:        "length",
						Label:     `Number of Symbols`,
						Comment:   element.LongText(`Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`),
						Type:      config.TypeText,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `4-5`,
					},

					&config.Field{
						// Path: customer/captcha/symbols
						ID:        "symbols",
						Label:     `Symbols Used in CAPTCHA`,
						Comment:   element.LongText(`Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`),
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `ABCDEFGHJKMnpqrstuvwxyz23456789`,
					},

					&config.Field{
						// Path: customer/captcha/case_sensitive
						ID:        "case_sensitive",
						Label:     `Case Sensitive`,
						Type:      config.TypeSelect,
						SortOrder: 9,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/allowed_resources
						ID:      `allowed_resources`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"captcha_folder":"captcha"}`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "admin",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "captcha",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/captcha/type
						ID:      `type`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `default`,
					},

					&config.Field{
						// Path: admin/captcha/failed_attempts_ip
						ID:      `failed_attempts_ip`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 1000,
					},

					&config.Field{
						// Path: admin/captcha/shown_to_logged_in_user
						ID:      `shown_to_logged_in_user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},

					&config.Field{
						// Path: admin/captcha/always_for
						ID:      `always_for`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"backend_forgotpassword":"1"}`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "customer",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "captcha",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: customer/captcha/type
						ID:      `type`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `default`,
					},

					&config.Field{
						// Path: customer/captcha/failed_attempts_ip
						ID:      `failed_attempts_ip`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 1000,
					},

					&config.Field{
						// Path: customer/captcha/shown_to_logged_in_user
						ID:      `shown_to_logged_in_user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"contact_us":"1"}`,
					},

					&config.Field{
						// Path: customer/captcha/always_for
						ID:      `always_for`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"user_create":"1","user_forgotpassword":"1","guest_checkout":"1","register_during_checkout":"1","contact_us":"1"}`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "captcha",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "_value",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: captcha/_value/fonts
						ID:      `fonts`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"linlibertine":{"label":"LinLibertine","path":"LinLibertineFont\/LinLibertine_Bd-2.8.1.ttf"}}`,
					},

					&config.Field{
						// Path: captcha/_value/frontend
						ID:      `frontend`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"areas":{"user_create":{"label":"Create user"},"user_login":{"label":"Login"},"user_forgotpassword":{"label":"Forgot password"},"guest_checkout":{"label":"Check Out as Guest"},"register_during_checkout":{"label":"Register during Checkout"},"contact_us":{"label":"Contact Us"}}}`,
					},

					&config.Field{
						// Path: captcha/_value/backend
						ID:      `backend`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"areas":{"backend_login":{"label":"Admin Login"},"backend_forgotpassword":{"label":"Admin Forgot Password"}}}`,
					},
				),
			},

			&config.Group{
				ID: "_attribute",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: captcha/_attribute/translate
						ID:      `translate`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `label`,
					},
				),
			},
		),
	},
)
