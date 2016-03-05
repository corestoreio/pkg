// +build ignore

package captcha

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "admin",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "captcha",
					Label:     `CAPTCHA`,
					SortOrder: 50,
					Scope:     scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/captcha/enable
							ID:        "enable",
							Label:     `Enable CAPTCHA in Admin`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/captcha/font
							ID:        "font",
							Label:     `Font`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `linlibertine`,
							// SourceModel: Magento\Captcha\Model\Config\Font
						},

						&element.Field{
							// Path: admin/captcha/forms
							ID:        "forms",
							Label:     `Forms`,
							Type:      element.TypeMultiselect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `backend_forgotpassword`,
							// SourceModel: Magento\Captcha\Model\Config\Form\Backend
						},

						&element.Field{
							// Path: admin/captcha/mode
							ID:        "mode",
							Label:     `Displaying Mode`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `after_fail`,
							// SourceModel: Magento\Captcha\Model\Config\Mode
						},

						&element.Field{
							// Path: admin/captcha/failed_attempts_login
							ID:        "failed_attempts_login",
							Label:     `Number of Unsuccessful Attempts to Login`,
							Comment:   text.Long(`If 0 is specified, CAPTCHA on the Login form will be always available.`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   3,
						},

						&element.Field{
							// Path: admin/captcha/timeout
							ID:        "timeout",
							Label:     `CAPTCHA Timeout (minutes)`,
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   7,
						},

						&element.Field{
							// Path: admin/captcha/length
							ID:        "length",
							Label:     `Number of Symbols`,
							Comment:   text.Long(`Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`),
							Type:      element.TypeText,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `4-5`,
						},

						&element.Field{
							// Path: admin/captcha/symbols
							ID:        "symbols",
							Label:     `Symbols Used in CAPTCHA`,
							Comment:   text.Long(`Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   `ABCDEFGHJKMnpqrstuvwxyz23456789`,
						},

						&element.Field{
							// Path: admin/captcha/case_sensitive
							ID:        "case_sensitive",
							Label:     `Case Sensitive`,
							Type:      element.TypeSelect,
							SortOrder: 9,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID: "customer",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "captcha",
					Label:     `CAPTCHA`,
					SortOrder: 110,
					Scope:     scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: customer/captcha/enable
							ID:        "enable",
							Label:     `Enable CAPTCHA on Storefront`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: customer/captcha/font
							ID:        "font",
							Label:     `Font`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `linlibertine`,
							// SourceModel: Magento\Captcha\Model\Config\Font
						},

						&element.Field{
							// Path: customer/captcha/forms
							ID:        "forms",
							Label:     `Forms`,
							Comment:   text.Long(`CAPTCHA for "Create user" and "Forgot password" forms is always enabled if chosen.`),
							Type:      element.TypeMultiselect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `user_forgotpassword`,
							// SourceModel: Magento\Captcha\Model\Config\Form\Frontend
						},

						&element.Field{
							// Path: customer/captcha/mode
							ID:        "mode",
							Label:     `Displaying Mode`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `after_fail`,
							// SourceModel: Magento\Captcha\Model\Config\Mode
						},

						&element.Field{
							// Path: customer/captcha/failed_attempts_login
							ID:        "failed_attempts_login",
							Label:     `Number of Unsuccessful Attempts to Login`,
							Comment:   text.Long(`If 0 is specified, CAPTCHA on the Login form will be always available.`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   3,
						},

						&element.Field{
							// Path: customer/captcha/timeout
							ID:        "timeout",
							Label:     `CAPTCHA Timeout (minutes)`,
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   7,
						},

						&element.Field{
							// Path: customer/captcha/length
							ID:        "length",
							Label:     `Number of Symbols`,
							Comment:   text.Long(`Please specify 8 symbols at the most. Range allowed (e.g. 3-5)`),
							Type:      element.TypeText,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `4-5`,
						},

						&element.Field{
							// Path: customer/captcha/symbols
							ID:        "symbols",
							Label:     `Symbols Used in CAPTCHA`,
							Comment:   text.Long(`Please use only letters (a-z or A-Z) or numbers (0-9) in this field. No spaces or other characters are allowed.<br />Similar looking characters (e.g. "i", "l", "1") decrease chance of correct recognition by customer.`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `ABCDEFGHJKMnpqrstuvwxyz23456789`,
						},

						&element.Field{
							// Path: customer/captcha/case_sensitive
							ID:        "case_sensitive",
							Label:     `Case Sensitive`,
							Type:      element.TypeSelect,
							SortOrder: 9,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "media_storage_configuration",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"captcha_folder":"captcha"}`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "admin",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "captcha",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/captcha/type
							ID:      `type`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `default`,
						},

						&element.Field{
							// Path: admin/captcha/failed_attempts_ip
							ID:      `failed_attempts_ip`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 1000,
						},

						&element.Field{
							// Path: admin/captcha/shown_to_logged_in_user
							ID:      `shown_to_logged_in_user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						&element.Field{
							// Path: admin/captcha/always_for
							ID:      `always_for`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"backend_forgotpassword":"1"}`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "customer",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "captcha",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: customer/captcha/type
							ID:      `type`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `default`,
						},

						&element.Field{
							// Path: customer/captcha/failed_attempts_ip
							ID:      `failed_attempts_ip`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 1000,
						},

						&element.Field{
							// Path: customer/captcha/shown_to_logged_in_user
							ID:      `shown_to_logged_in_user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"contact_us":"1"}`,
						},

						&element.Field{
							// Path: customer/captcha/always_for
							ID:      `always_for`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"user_create":"1","user_forgotpassword":"1","guest_checkout":"1","register_during_checkout":"1","contact_us":"1"}`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "captcha",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "_value",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: captcha/_value/fonts
							ID:      `fonts`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"linlibertine":{"label":"LinLibertine","path":"LinLibertineFont\/LinLibertine_Bd-2.8.1.ttf"}}`,
						},

						&element.Field{
							// Path: captcha/_value/frontend
							ID:      `frontend`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"areas":{"user_create":{"label":"Create user"},"user_login":{"label":"Login"},"user_forgotpassword":{"label":"Forgot password"},"guest_checkout":{"label":"Check Out as Guest"},"register_during_checkout":{"label":"Register during Checkout"},"contact_us":{"label":"Contact Us"}}}`,
						},

						&element.Field{
							// Path: captcha/_value/backend
							ID:      `backend`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"areas":{"backend_login":{"label":"Admin Login"},"backend_forgotpassword":{"label":"Admin Forgot Password"}}}`,
						},
					),
				},

				&element.Group{
					ID: "_attribute",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: captcha/_attribute/translate
							ID:      `translate`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `label`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
