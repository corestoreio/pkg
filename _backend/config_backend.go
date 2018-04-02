// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

// TODO: during development move each of this config stuff into its own package.

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		&element.Section{
			ID:        "advanced",
			Label:     `Advanced`,
			SortOrder: 910,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::advanced
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "modules_disable_output",
					Label:     `Disable Modules Output`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields:    element.MakeFields(),
				},
			),
		},
		&element.Section{
			ID:        "trans_email",
			Label:     `Store Email Addresses`,
			SortOrder: 90,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::trans_email
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "ident_custom1",
					Label:     `Custom Email 1`,
					SortOrder: 4,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: trans_email/ident_custom1/email
							ID:        "email",
							Label:     `Sender Email`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_custom1/name
							ID:        "name",
							Label:     `Sender Name`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        "ident_custom2",
					Label:     `Custom Email 2`,
					SortOrder: 5,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: trans_email/ident_custom2/email
							ID:        "email",
							Label:     `Sender Email`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_custom2/name
							ID:        "name",
							Label:     `Sender Name`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        "ident_general",
					Label:     `General Contact`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: trans_email/ident_general/email
							ID:        "email",
							Label:     `Sender Email`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_general/name
							ID:        "name",
							Label:     `Sender Name`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        "ident_sales",
					Label:     `Sales Representative`,
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: trans_email/ident_sales/email
							ID:        "email",
							Label:     `Sender Email`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_sales/name
							ID:        "name",
							Label:     `Sender Name`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        "ident_support",
					Label:     `Customer Support`,
					SortOrder: 3,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: trans_email/ident_support/email
							ID:        "email",
							Label:     `Sender Email`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_support/name
							ID:        "name",
							Label:     `Sender Name`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "design",
			Label:     `Design`,
			SortOrder: 30,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Config::config_design
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "theme",
					Label:     `Design Theme`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: design/theme/theme_id
							ID:        "theme_id",
							Label:     `Design Theme`,
							Comment:   `If no value is specified, the system default will be used. The system default may be modified by third party extensions.`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Theme\Model\Design\Backend\Theme
							// SourceModel: Magento\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
						},

						&element.Field{
							// Path: design/theme/ua_regexp
							ID:        "ua_regexp",
							Label:     `User-Agent Exceptions`,
							Comment:   `Search strings are either normal strings or regular exceptions (PCRE. They are matched in the same order as entered. Examples:<br /><span style="font-family:monospace">Firefox<br />/^mozilla/i</span>` ,
							Tooltip:   `Find a string in client user-agent header and switch to specific design theme for that browser.`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Theme\Model\Design\Backend\Exceptions
						},
					),
				},

				&element.Group{
					ID:        "pagination",
					Label:     `Pagination`,
					SortOrder: 500,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: design/pagination/pagination_frame
							ID:        "pagination_frame",
							Label:     `Pagination Frame`,
							Comment:   `How many links to display at once.`,
							Type:      element.TypeText,
							SortOrder: 7,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/pagination_frame_skip
							ID:        "pagination_frame_skip",
							Label:     `Pagination Frame Skip`,
							Comment:   `If the current frame position does not cover utmost pages, will render link to current position plus/minus this value.`,
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/anchor_text_for_previous
							ID:        "anchor_text_for_previous",
							Label:     `Anchor Text for Previous`,
							Comment:   `Alternative text for previous link in pagination menu. If empty, default arrow image will used.`,
							Type:      element.TypeText,
							SortOrder: 9,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/anchor_text_for_next
							ID:        "anchor_text_for_next",
							Label:     `Anchor Text for Next`,
							Comment:   `Alternative text for next link in pagination menu. If empty, default arrow image will used.`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "dev",
			Label:     `Developer`,
			SortOrder: 920,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::dev
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "debug",
					Label:     `Debug`,
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/debug/template_hints_storefront
							ID:        "template_hints_storefront",
							Label:     `Enabled Template Path Hints for Storefront`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/debug/template_hints_admin
							ID:        "template_hints_admin",
							Label:     `Enabled Template Path Hints for Admin`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/debug/template_hints_blocks
							ID:        "template_hints_blocks",
							Label:     `Add Block Names to Hints`,
							Type:      element.TypeSelect,
							SortOrder: 21,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "template",
					Label:     `Template Settings`,
					SortOrder: 25,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/template/allow_symlink
							ID:        "allow_symlink",
							Label:     `Allow Symlinks`,
							Comment:   `Warning! Enabling this feature is not recommended on production environments because it represents a potential security risk.`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/template/minify_html
							ID:        "minify_html",
							Label:     `Minify Html`,
							Type:      element.TypeSelect,
							SortOrder: 25,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "translate_inline",
					Label:     `Translate Inline`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/translate_inline/active
							ID:        "active",
							Label:     `Enabled for Storefront`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/translate_inline/active_admin
							ID:        "active_admin",
							Label:     `Enabled for Admin`,
							Comment:   `Translate, blocks and other output caches should be disabled for both Storefront and Admin inline translations.`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "js",
					Label:     `JavaScript Settings`,
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/js/merge_files
							ID:        "merge_files",
							Label:     `Merge JavaScript Files`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/js/enable_js_bundling
							ID:        "enable_js_bundling",
							Label:     `Enable JavaScript Bundling`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/js/minify_files
							ID:        "minify_files",
							Label:     `Minify JavaScript Files`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "css",
					Label:     `CSS Settings`,
					SortOrder: 110,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/css/merge_css_files
							ID:        "merge_css_files",
							Label:     `Merge CSS Files`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/css/minify_files
							ID:        "minify_files",
							Label:     `Minify CSS Files`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "image",
					Label:     `Image Processing Settings`,
					SortOrder: 120,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/image/default_adapter
							ID:        "default_adapter",
							Label:     `Image Adapter`,
							Comment:   `When the adapter was changed, please flush Catalog Images Cache.`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Image\Adapter
							// SourceModel: Magento\Config\Model\Config\Source\Image\Adapter
						},
					),
				},

				&element.Group{
					ID:        "static",
					Label:     `Static Files Settings`,
					SortOrder: 130,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: dev/static/sign
							ID:        "sign",
							Label:     `Sign Static Files`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		&element.Section{
			ID:        "system",
			Label:     `System`,
			SortOrder: 900,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Config::config_system
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "smtp",
					Label:     `Mail Sending Settings`,
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: system/smtp/disable
							ID:        "disable",
							Label:     `Disable Email Communications`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: system/smtp/host
							ID:        "host",
							Label:     `Host`,
							Comment:   `For Windows server only.`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: system/smtp/port
							ID:        "port",
							Label:     `Port (25)` ,
							Comment:   `For Windows server only.`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: system/smtp/set_return_path
							ID:        "set_return_path",
							Label:     `Set Return-Path`,
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesnocustom
						},

						&element.Field{
							// Path: system/smtp/return_path_email
							ID:        "return_path_email",
							Label:     `Return-Path Email`,
							Type:      element.TypeText,
							SortOrder: 80,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "admin",
			Label:     `Admin`,
			SortOrder: 20,
			Scopes:    scope.PermDefault,
			Resource:  0, // Magento_Config::config_admin
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "emails",
					Label:     `Admin User Emails`,
					SortOrder: 10,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: admin/emails/forgot_email_template
							ID:        "forgot_email_template",
							Label:     `Forgot Password Email Template`,
							Comment:   `Email template chosen based on theme fallback when "Default" option is selected.`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: admin/emails/forgot_email_identity
							ID:        "forgot_email_identity",
							Label:     `Forgot and Reset Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: admin/emails/password_reset_link_expiration_period
							ID:        "password_reset_link_expiration_period",
							Label:     `Recovery Link Expiration Period (days)` ,
							Comment:   `Please enter a number 1 or greater in this field.`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
						},
					),
				},

				&element.Group{
					ID:        "startup",
					Label:     `Startup Page`,
					SortOrder: 20,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: admin/startup/menu_item_id
							ID:        "menu_item_id",
							Label:     `Startup Page`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Admin\Page
						},
					),
				},

				&element.Group{
					ID:        "url",
					Label:     `Admin Base URL`,
					SortOrder: 30,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: admin/url/use_custom
							ID:        "use_custom",
							Label:     `Use Custom Admin URL`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usecustom
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/url/custom
							ID:        "custom",
							Label:     `Custom Admin URL`,
							Comment:   `Make sure that base URL ends with '/' (slash, e.g. http://yourdomain/magento/`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custom
						},

						&element.Field{
							// Path: admin/url/use_custom_path
							ID:        "use_custom_path",
							Label:     `Use Custom Admin Path`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/url/custom_path
							ID:        "custom_path",
							Label:     `Custom Admin Path`,
							Comment:   `You will have to sign in after you save your custom admin cfgpath.`,
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
						},
					),
				},

				&element.Group{
					ID:        "security",
					Label:     `Security`,
					SortOrder: 35,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: admin/security/use_form_key
							ID:        "use_form_key",
							Label:     `Add Secret Key to URLs`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usesecretkey
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/security/use_case_sensitive_login
							ID:        "use_case_sensitive_login",
							Label:     `Login is Case Sensitive`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/security/session_lifetime
							ID:        "session_lifetime",
							Label:     `Admin Session Lifetime (seconds)`,
							Comment:   `Values less than 60 are ignored.`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   true,
							Scopes:    scope.PermDefault,
						},
					),
				},

				&element.Group{
					ID:        "dashboard",
					Label:     `Dashboard`,
					SortOrder: 40,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: admin/dashboard/enable_charts
							ID:        "enable_charts",
							Label:     `Enable Charts`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID:        "web",
			Label:     `Web`,
			SortOrder: 20,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::web
			Groups: element.MakeGroups(
				&element.Group{
					ID:        "url",
					Label:     `Url Options`,
					SortOrder: 3,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/url/use_store
							ID:        "use_store",
							Label:     `Add Store Code to Urls`,
							Comment:   `<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc..`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermDefault,
							Default:   false,
							// BackendModel: Magento\Config\Model\Config\Backend\Store
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/url/redirect_to_base
							ID:        "redirect_to_base",
							Label:     `Auto-redirect to Base URL`,
							Comment:   `I.e. redirect from http://example.com/store/ to http://www.example.com/store/`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermDefault,
							Default:   1,
							// SourceModel: Magento\Config\Model\Config\Source\Web\Redirect
						},
					),
				},

				&element.Group{
					ID:        "seo",
					Label:     `Search Engine Optimization`,
					SortOrder: 5,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/seo/use_rewrites
							ID:        "use_rewrites",
							Label:     `Use Web Server Rewrites`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "unsecure",
					Label:     `Base URLs`,
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash e.g. http://example.com/magento/`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/unsecure/base_url
							ID:        "base_url",
							Label:     `Base URL`,
							Comment:   `Specify URL or {{base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_link_url
							ID:        "base_link_url",
							Label:     `Base Link URL`,
							Comment:   `May start with {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_static_url
							ID:        "base_static_url",
							Label:     `Base URL for Static View Files`,
							Comment:   `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_media_url
							ID:        "base_media_url",
							Label:     `Base URL for User Media Files`,
							Comment:   `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},
					),
				},

				&element.Group{
					ID:        "secure",
					Label:     `Base URLs (Secure`),
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash e.g. https://example.com/magento/`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/secure/base_url
							ID:        "base_url",
							Label:     `Secure Base URL`,
							Comment:   `Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_link_url
							ID:        "base_link_url",
							Label:     `Secure Base Link URL`,
							Comment:   `May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_static_url
							ID:        "base_static_url",
							Label:     `Secure Base URL for Static View Files`,
							Comment:   `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_media_url
							ID:        "base_media_url",
							Label:     `Secure Base URL for User Media Files`,
							Comment:   `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/use_in_frontend
							ID:        "use_in_frontend",
							Label:     `Use Secure URLs on Storefront`,
							Comment:   `Enter https protocol to use Secure URLs on Storefront.`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/use_in_adminhtml
							ID:        "use_in_adminhtml",
							Label:     `Use Secure URLs in Admin`,
							Comment:   `Enter https protocol to use Secure URLs in Admin.`,
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/enable_hsts
							ID:        "enable_hsts",
							Label:     `Enable HTTP Strict Transport Security (HSTS`),
							Comment:   `See <a href="https://www.owasp.org/index.php/HTTP_Strict_Transport_Security" target="_blank">HTTP Strict Transport Security</a> page for details.`,
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/enable_upgrade_insecure
							ID:        "enable_upgrade_insecure",
							Label:     `Upgrade Insecure Requests`,
							Comment:   `See <a href="http://www.w3.org/TR/upgrade-insecure-requests/" target="_blank">Upgrade Insecure Requests</a> page for details.`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/offloader_header
							ID:        "offloader_header",
							Label:     `Offloader header`,
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   true,
							Scopes:    scope.PermDefault,
						},
					),
				},

				&element.Group{
					ID:        "default",
					Label:     `Default Pages`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/default/front
							ID:        "front",
							Label:     `Default Web URL`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: web/default/no_route
							ID:        "no_route",
							Label:     `Default No-route URL`,
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
						},
					),
				},

				&element.Group{
					ID:        "session",
					Label:     `Session Validation Settings`,
					SortOrder: 60,
					Scopes:    scope.PermWebsite,
					Fields: element.MakeFields(
						&element.Field{
							// Path: web/session/use_remote_addr
							ID:        "use_remote_addr",
							Label:     `Validate REMOTE_ADDR`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_via
							ID:        "use_http_via",
							Label:     `Validate HTTP_VIA`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_x_forwarded_for
							ID:        "use_http_x_forwarded_for",
							Label:     `Validate HTTP_X_FORWARDED_FOR`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_user_agent
							ID:        "use_http_user_agent",
							Label:     `Validate HTTP_USER_AGENT`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_frontend_sid
							ID:        "use_frontend_sid",
							Label:     `Use SID on Storefront`,
							Comment:   `Allows customers to stay logged in when switching between different stores.`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "system",
			Groups: element.MakeGroups(
				&element.Group{
					ID: "media_storage_configuration",
					Fields: element.MakeFields(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
														Default: `{"email_folder":"email"}`,
						},
					),
				},

				&element.Group{
					ID: "emails",
					Fields: element.MakeFields(
						&element.Field{
							// Path: system/emails/forgot_email_template
							ID:      `forgot_email_template`,
							Type:    element.TypeHidden,	
														Default: `system_emails_forgot_email_template`,
						},

						&element.Field{
							// Path: system/emails/forgot_email_identity
							ID:      `forgot_email_identity`,
							Type:    element.TypeHidden,
														Default: `general`,
						},
					),
				},

				&element.Group{
					ID: "dashboard",
					Fields: element.MakeFields(
						&element.Field{
							// Path: system/dashboard/enable_charts
							ID:      `enable_charts`,
							Type:    element.TypeHidden,
														Default: true,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "general",
			Groups: element.MakeGroups(
				&element.Group{
					ID: "validator_data",
					Fields: element.MakeFields(
						&element.Field{
							// Path: general/validator_data/input_types
							ID:      `input_types`,
							Type:    element.TypeHidden,
														Default: `{"price":"price","media_image":"media_image","gallery":"gallery"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
