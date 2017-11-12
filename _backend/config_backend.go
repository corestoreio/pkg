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
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

// TODO: during development move each of this config stuff into its own package.

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        cfgpath.NewRoute("advanced"),
			Label:     text.Chars(`Advanced`),
			SortOrder: 910,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::advanced
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("modules_disable_output"),
					Label:     text.Chars(`Disable Modules Output`),
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID:        cfgpath.NewRoute("trans_email"),
			Label:     text.Chars(`Store Email Addresses`),
			SortOrder: 90,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::trans_email
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("ident_custom1"),
					Label:     text.Chars(`Custom Email 1`),
					SortOrder: 4,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_custom1/email
							ID:        cfgpath.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						element.Field{
							// Path: trans_email/ident_custom1/name
							ID:        cfgpath.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("ident_custom2"),
					Label:     text.Chars(`Custom Email 2`),
					SortOrder: 5,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_custom2/email
							ID:        cfgpath.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						element.Field{
							// Path: trans_email/ident_custom2/name
							ID:        cfgpath.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("ident_general"),
					Label:     text.Chars(`General Contact`),
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_general/email
							ID:        cfgpath.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						element.Field{
							// Path: trans_email/ident_general/name
							ID:        cfgpath.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("ident_sales"),
					Label:     text.Chars(`Sales Representative`),
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_sales/email
							ID:        cfgpath.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						element.Field{
							// Path: trans_email/ident_sales/name
							ID:        cfgpath.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("ident_support"),
					Label:     text.Chars(`Customer Support`),
					SortOrder: 3,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_support/email
							ID:        cfgpath.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						element.Field{
							// Path: trans_email/ident_support/name
							ID:        cfgpath.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},
			),
		},
		element.Section{
			ID:        cfgpath.NewRoute("design"),
			Label:     text.Chars(`Design`),
			SortOrder: 30,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Config::config_design
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("theme"),
					Label:     text.Chars(`Design Theme`),
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: design/theme/theme_id
							ID:        cfgpath.NewRoute("theme_id"),
							Label:     text.Chars(`Design Theme`),
							Comment:   text.Chars(`If no value is specified, the system default will be used. The system default may be modified by third party extensions.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Theme\Model\Design\Backend\Theme
							// SourceModel: Magento\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
						},

						element.Field{
							// Path: design/theme/ua_regexp
							ID:        cfgpath.NewRoute("ua_regexp"),
							Label:     text.Chars(`User-Agent Exceptions`),
							Comment:   text.Chars(`Search strings are either normal strings or regular exceptions (PCRE). They are matched in the same order as entered. Examples:<br /><span style="font-family:monospace">Firefox<br />/^mozilla/i</span>`),
							Tooltip:   text.Chars(`Find a string in client user-agent header and switch to specific design theme for that browser.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Theme\Model\Design\Backend\Exceptions
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("pagination"),
					Label:     text.Chars(`Pagination`),
					SortOrder: 500,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: design/pagination/pagination_frame
							ID:        cfgpath.NewRoute("pagination_frame"),
							Label:     text.Chars(`Pagination Frame`),
							Comment:   text.Chars(`How many links to display at once.`),
							Type:      element.TypeText,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/pagination/pagination_frame_skip
							ID:        cfgpath.NewRoute("pagination_frame_skip"),
							Label:     text.Chars(`Pagination Frame Skip`),
							Comment:   text.Chars(`If the current frame position does not cover utmost pages, will render link to current position plus/minus this value.`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/pagination/anchor_text_for_previous
							ID:        cfgpath.NewRoute("anchor_text_for_previous"),
							Label:     text.Chars(`Anchor Text for Previous`),
							Comment:   text.Chars(`Alternative text for previous link in pagination menu. If empty, default arrow image will used.`),
							Type:      element.TypeText,
							SortOrder: 9,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/pagination/anchor_text_for_next
							ID:        cfgpath.NewRoute("anchor_text_for_next"),
							Label:     text.Chars(`Anchor Text for Next`),
							Comment:   text.Chars(`Alternative text for next link in pagination menu. If empty, default arrow image will used.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		element.Section{
			ID:        cfgpath.NewRoute("dev"),
			Label:     text.Chars(`Developer`),
			SortOrder: 920,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::dev
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("debug"),
					Label:     text.Chars(`Debug`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/debug/template_hints_storefront
							ID:        cfgpath.NewRoute("template_hints_storefront"),
							Label:     text.Chars(`Enabled Template Path Hints for Storefront`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/debug/template_hints_admin
							ID:        cfgpath.NewRoute("template_hints_admin"),
							Label:     text.Chars(`Enabled Template Path Hints for Admin`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/debug/template_hints_blocks
							ID:        cfgpath.NewRoute("template_hints_blocks"),
							Label:     text.Chars(`Add Block Names to Hints`),
							Type:      element.TypeSelect,
							SortOrder: 21,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("template"),
					Label:     text.Chars(`Template Settings`),
					SortOrder: 25,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/template/allow_symlink
							ID:        cfgpath.NewRoute("allow_symlink"),
							Label:     text.Chars(`Allow Symlinks`),
							Comment:   text.Chars(`Warning! Enabling this feature is not recommended on production environments because it represents a potential security risk.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/template/minify_html
							ID:        cfgpath.NewRoute("minify_html"),
							Label:     text.Chars(`Minify Html`),
							Type:      element.TypeSelect,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("translate_inline"),
					Label:     text.Chars(`Translate Inline`),
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/translate_inline/active
							ID:        cfgpath.NewRoute("active"),
							Label:     text.Chars(`Enabled for Storefront`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/translate_inline/active_admin
							ID:        cfgpath.NewRoute("active_admin"),
							Label:     text.Chars(`Enabled for Admin`),
							Comment:   text.Chars(`Translate, blocks and other output caches should be disabled for both Storefront and Admin inline translations.`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("js"),
					Label:     text.Chars(`JavaScript Settings`),
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/js/merge_files
							ID:        cfgpath.NewRoute("merge_files"),
							Label:     text.Chars(`Merge JavaScript Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/js/enable_js_bundling
							ID:        cfgpath.NewRoute("enable_js_bundling"),
							Label:     text.Chars(`Enable JavaScript Bundling`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/js/minify_files
							ID:        cfgpath.NewRoute("minify_files"),
							Label:     text.Chars(`Minify JavaScript Files`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("css"),
					Label:     text.Chars(`CSS Settings`),
					SortOrder: 110,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/css/merge_css_files
							ID:        cfgpath.NewRoute("merge_css_files"),
							Label:     text.Chars(`Merge CSS Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: dev/css/minify_files
							ID:        cfgpath.NewRoute("minify_files"),
							Label:     text.Chars(`Minify CSS Files`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("image"),
					Label:     text.Chars(`Image Processing Settings`),
					SortOrder: 120,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/image/default_adapter
							ID:        cfgpath.NewRoute("default_adapter"),
							Label:     text.Chars(`Image Adapter`),
							Comment:   text.Chars(`When the adapter was changed, please flush Catalog Images Cache.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Image\Adapter
							// SourceModel: Magento\Config\Model\Config\Source\Image\Adapter
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("static"),
					Label:     text.Chars(`Static Files Settings`),
					SortOrder: 130,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: dev/static/sign
							ID:        cfgpath.NewRoute("sign"),
							Label:     text.Chars(`Sign Static Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		element.Section{
			ID:        cfgpath.NewRoute("system"),
			Label:     text.Chars(`System`),
			SortOrder: 900,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Config::config_system
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("smtp"),
					Label:     text.Chars(`Mail Sending Settings`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/smtp/disable
							ID:        cfgpath.NewRoute("disable"),
							Label:     text.Chars(`Disable Email Communications`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: system/smtp/host
							ID:        cfgpath.NewRoute("host"),
							Label:     text.Chars(`Host`),
							Comment:   text.Chars(`For Windows server only.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: system/smtp/port
							ID:        cfgpath.NewRoute("port"),
							Label:     text.Chars(`Port (25)`),
							Comment:   text.Chars(`For Windows server only.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: system/smtp/set_return_path
							ID:        cfgpath.NewRoute("set_return_path"),
							Label:     text.Chars(`Set Return-Path`),
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesnocustom
						},

						element.Field{
							// Path: system/smtp/return_path_email
							ID:        cfgpath.NewRoute("return_path_email"),
							Label:     text.Chars(`Return-Path Email`),
							Type:      element.TypeText,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},
					),
				},
			),
		},
		element.Section{
			ID:        cfgpath.NewRoute("admin"),
			Label:     text.Chars(`Admin`),
			SortOrder: 20,
			Scopes:    scope.PermDefault,
			Resource:  0, // Magento_Config::config_admin
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("emails"),
					Label:     text.Chars(`Admin User Emails`),
					SortOrder: 10,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: admin/emails/forgot_email_template
							ID:        cfgpath.NewRoute("forgot_email_template"),
							Label:     text.Chars(`Forgot Password Email Template`),
							Comment:   text.Chars(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: admin/emails/forgot_email_identity
							ID:        cfgpath.NewRoute("forgot_email_identity"),
							Label:     text.Chars(`Forgot and Reset Email Sender`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: admin/emails/password_reset_link_expiration_period
							ID:        cfgpath.NewRoute("password_reset_link_expiration_period"),
							Label:     text.Chars(`Recovery Link Expiration Period (days)`),
							Comment:   text.Chars(`Please enter a number 1 or greater in this field.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("startup"),
					Label:     text.Chars(`Startup Page`),
					SortOrder: 20,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: admin/startup/menu_item_id
							ID:        cfgpath.NewRoute("menu_item_id"),
							Label:     text.Chars(`Startup Page`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Admin\Page
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("url"),
					Label:     text.Chars(`Admin Base URL`),
					SortOrder: 30,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: admin/url/use_custom
							ID:        cfgpath.NewRoute("use_custom"),
							Label:     text.Chars(`Use Custom Admin URL`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usecustom
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: admin/url/custom
							ID:        cfgpath.NewRoute("custom"),
							Label:     text.Chars(`Custom Admin URL`),
							Comment:   text.Chars(`Make sure that base URL ends with '/' (slash), e.g. http://yourdomain/magento/`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custom
						},

						element.Field{
							// Path: admin/url/use_custom_path
							ID:        cfgpath.NewRoute("use_custom_path"),
							Label:     text.Chars(`Use Custom Admin Path`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: admin/url/custom_path
							ID:        cfgpath.NewRoute("custom_path"),
							Label:     text.Chars(`Custom Admin Path`),
							Comment:   text.Chars(`You will have to sign in after you save your custom admin cfgpath.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("security"),
					Label:     text.Chars(`Security`),
					SortOrder: 35,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: admin/security/use_form_key
							ID:        cfgpath.NewRoute("use_form_key"),
							Label:     text.Chars(`Add Secret Key to URLs`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usesecretkey
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: admin/security/use_case_sensitive_login
							ID:        cfgpath.NewRoute("use_case_sensitive_login"),
							Label:     text.Chars(`Login is Case Sensitive`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: admin/security/session_lifetime
							ID:        cfgpath.NewRoute("session_lifetime"),
							Label:     text.Chars(`Admin Session Lifetime (seconds)`),
							Comment:   text.Chars(`Values less than 60 are ignored.`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("dashboard"),
					Label:     text.Chars(`Dashboard`),
					SortOrder: 40,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: admin/dashboard/enable_charts
							ID:        cfgpath.NewRoute("enable_charts"),
							Label:     text.Chars(`Enable Charts`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		element.Section{
			ID:        cfgpath.NewRoute("web"),
			Label:     text.Chars(`Web`),
			SortOrder: 20,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::web
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        cfgpath.NewRoute("url"),
					Label:     text.Chars(`Url Options`),
					SortOrder: 3,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/url/use_store
							ID:        cfgpath.NewRoute("use_store"),
							Label:     text.Chars(`Add Store Code to Urls`),
							Comment:   text.Chars(`<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   false,
							// BackendModel: Magento\Config\Model\Config\Backend\Store
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/url/redirect_to_base
							ID:        cfgpath.NewRoute("redirect_to_base"),
							Label:     text.Chars(`Auto-redirect to Base URL`),
							Comment:   text.Chars(`I.e. redirect from http://example.com/store/ to http://www.example.com/store/`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   1,
							// SourceModel: Magento\Config\Model\Config\Source\Web\Redirect
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("seo"),
					Label:     text.Chars(`Search Engine Optimization`),
					SortOrder: 5,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/seo/use_rewrites
							ID:        cfgpath.NewRoute("use_rewrites"),
							Label:     text.Chars(`Use Web Server Rewrites`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("unsecure"),
					Label:     text.Chars(`Base URLs`),
					Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/unsecure/base_url
							ID:        cfgpath.NewRoute("base_url"),
							Label:     text.Chars(`Base URL`),
							Comment:   text.Chars(`Specify URL or {{base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/unsecure/base_link_url
							ID:        cfgpath.NewRoute("base_link_url"),
							Label:     text.Chars(`Base Link URL`),
							Comment:   text.Chars(`May start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/unsecure/base_static_url
							ID:        cfgpath.NewRoute("base_static_url"),
							Label:     text.Chars(`Base URL for Static View Files`),
							Comment:   text.Chars(`May be empty or start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/unsecure/base_media_url
							ID:        cfgpath.NewRoute("base_media_url"),
							Label:     text.Chars(`Base URL for User Media Files`),
							Comment:   text.Chars(`May be empty or start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("secure"),
					Label:     text.Chars(`Base URLs (Secure)`),
					Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/secure/base_url
							ID:        cfgpath.NewRoute("base_url"),
							Label:     text.Chars(`Secure Base URL`),
							Comment:   text.Chars(`Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/secure/base_link_url
							ID:        cfgpath.NewRoute("base_link_url"),
							Label:     text.Chars(`Secure Base Link URL`),
							Comment:   text.Chars(`May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/secure/base_static_url
							ID:        cfgpath.NewRoute("base_static_url"),
							Label:     text.Chars(`Secure Base URL for Static View Files`),
							Comment:   text.Chars(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/secure/base_media_url
							ID:        cfgpath.NewRoute("base_media_url"),
							Label:     text.Chars(`Secure Base URL for User Media Files`),
							Comment:   text.Chars(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						element.Field{
							// Path: web/secure/use_in_frontend
							ID:        cfgpath.NewRoute("use_in_frontend"),
							Label:     text.Chars(`Use Secure URLs on Storefront`),
							Comment:   text.Chars(`Enter https protocol to use Secure URLs on Storefront.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/secure/use_in_adminhtml
							ID:        cfgpath.NewRoute("use_in_adminhtml"),
							Label:     text.Chars(`Use Secure URLs in Admin`),
							Comment:   text.Chars(`Enter https protocol to use Secure URLs in Admin.`),
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/secure/enable_hsts
							ID:        cfgpath.NewRoute("enable_hsts"),
							Label:     text.Chars(`Enable HTTP Strict Transport Security (HSTS)`),
							Comment:   text.Chars(`See <a href="https://www.owasp.org/index.php/HTTP_Strict_Transport_Security" target="_blank">HTTP Strict Transport Security</a> page for details.`),
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/secure/enable_upgrade_insecure
							ID:        cfgpath.NewRoute("enable_upgrade_insecure"),
							Label:     text.Chars(`Upgrade Insecure Requests`),
							Comment:   text.Chars(`See <a href="http://www.w3.org/TR/upgrade-insecure-requests/" target="_blank">Upgrade Insecure Requests</a> page for details.`),
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/secure/offloader_header
							ID:        cfgpath.NewRoute("offloader_header"),
							Label:     text.Chars(`Offloader header`),
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("default"),
					Label:     text.Chars(`Default Pages`),
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/default/front
							ID:        cfgpath.NewRoute("front"),
							Label:     text.Chars(`Default Web URL`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: web/default/no_route
							ID:        cfgpath.NewRoute("no_route"),
							Label:     text.Chars(`Default No-route URL`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},

				element.Group{
					ID:        cfgpath.NewRoute("session"),
					Label:     text.Chars(`Session Validation Settings`),
					SortOrder: 60,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: web/session/use_remote_addr
							ID:        cfgpath.NewRoute("use_remote_addr"),
							Label:     text.Chars(`Validate REMOTE_ADDR`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/session/use_http_via
							ID:        cfgpath.NewRoute("use_http_via"),
							Label:     text.Chars(`Validate HTTP_VIA`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/session/use_http_x_forwarded_for
							ID:        cfgpath.NewRoute("use_http_x_forwarded_for"),
							Label:     text.Chars(`Validate HTTP_X_FORWARDED_FOR`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/session/use_http_user_agent
							ID:        cfgpath.NewRoute("use_http_user_agent"),
							Label:     text.Chars(`Validate HTTP_USER_AGENT`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: web/session/use_frontend_sid
							ID:        cfgpath.NewRoute("use_frontend_sid"),
							Label:     text.Chars(`Use SID on Storefront`),
							Comment:   text.Chars(`Allows customers to stay logged in when switching between different stores.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: cfgpath.NewRoute("system"),
			Groups: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute("media_storage_configuration"),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      cfgpath.NewRoute(`allowed_resources`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"email_folder":"email"}`,
						},
					),
				},

				element.Group{
					ID: cfgpath.NewRoute("emails"),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/emails/forgot_email_template
							ID:      cfgpath.NewRoute(`forgot_email_template`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `system_emails_forgot_email_template`,
						},

						element.Field{
							// Path: system/emails/forgot_email_identity
							ID:      cfgpath.NewRoute(`forgot_email_identity`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `general`,
						},
					),
				},

				element.Group{
					ID: cfgpath.NewRoute("dashboard"),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/dashboard/enable_charts
							ID:      cfgpath.NewRoute(`enable_charts`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
		element.Section{
			ID: cfgpath.NewRoute("general"),
			Groups: element.NewGroupSlice(
				element.Group{
					ID: cfgpath.NewRoute("validator_data"),
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/validator_data/input_types
							ID:      cfgpath.NewRoute(`input_types`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"price":"price","media_image":"media_image","gallery":"gallery"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
