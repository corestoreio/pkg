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
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

// TODO: during development move each of this config stuff into its own package.

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        path.NewRoute("advanced"),
			Label:     text.Chars(`Advanced`),
			SortOrder: 910,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Backend::advanced
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("modules_disable_output"),
					Label:     text.Chars(`Disable Modules Output`),
					SortOrder: 2,
					Scope:     scope.PermStore,
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("trans_email"),
			Label:     text.Chars(`Store Email Addresses`),
			SortOrder: 90,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Backend::trans_email
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("ident_custom1"),
					Label:     text.Chars(`Custom Email 1`),
					SortOrder: 4,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: trans_email/ident_custom1/email
							ID:        path.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_custom1/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("ident_custom2"),
					Label:     text.Chars(`Custom Email 2`),
					SortOrder: 5,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: trans_email/ident_custom2/email
							ID:        path.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_custom2/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("ident_general"),
					Label:     text.Chars(`General Contact`),
					SortOrder: 1,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: trans_email/ident_general/email
							ID:        path.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_general/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("ident_sales"),
					Label:     text.Chars(`Sales Representative`),
					SortOrder: 2,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: trans_email/ident_sales/email
							ID:        path.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_sales/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("ident_support"),
					Label:     text.Chars(`Customer Support`),
					SortOrder: 3,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: trans_email/ident_support/email
							ID:        path.NewRoute("email"),
							Label:     text.Chars(`Sender Email`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},

						&element.Field{
							// Path: trans_email/ident_support/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Sender Name`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("design"),
			Label:     text.Chars(`Design`),
			SortOrder: 30,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Config::config_design
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("theme"),
					Label:     text.Chars(`Design Theme`),
					SortOrder: 1,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/theme/theme_id
							ID:        path.NewRoute("theme_id"),
							Label:     text.Chars(`Design Theme`),
							Comment:   text.Chars(`If no value is specified, the system default will be used. The system default may be modified by third party extensions.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Theme\Model\Design\Backend\Theme
							// SourceModel: Magento\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
						},

						&element.Field{
							// Path: design/theme/ua_regexp
							ID:        path.NewRoute("ua_regexp"),
							Label:     text.Chars(`User-Agent Exceptions`),
							Comment:   text.Chars(`Search strings are either normal strings or regular exceptions (PCRE). They are matched in the same order as entered. Examples:<br /><span style="font-family:monospace">Firefox<br />/^mozilla/i</span>`),
							Tooltip:   text.Chars(`Find a string in client user-agent header and switch to specific design theme for that browser.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// BackendModel: Magento\Theme\Model\Design\Backend\Exceptions
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("pagination"),
					Label:     text.Chars(`Pagination`),
					SortOrder: 500,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/pagination/pagination_frame
							ID:        path.NewRoute("pagination_frame"),
							Label:     text.Chars(`Pagination Frame`),
							Comment:   text.Chars(`How many links to display at once.`),
							Type:      element.TypeText,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/pagination_frame_skip
							ID:        path.NewRoute("pagination_frame_skip"),
							Label:     text.Chars(`Pagination Frame Skip`),
							Comment:   text.Chars(`If the current frame position does not cover utmost pages, will render link to current position plus/minus this value.`),
							Type:      element.TypeText,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/anchor_text_for_previous
							ID:        path.NewRoute("anchor_text_for_previous"),
							Label:     text.Chars(`Anchor Text for Previous`),
							Comment:   text.Chars(`Alternative text for previous link in pagination menu. If empty, default arrow image will used.`),
							Type:      element.TypeText,
							SortOrder: 9,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: design/pagination/anchor_text_for_next
							ID:        path.NewRoute("anchor_text_for_next"),
							Label:     text.Chars(`Anchor Text for Next`),
							Comment:   text.Chars(`Alternative text for next link in pagination menu. If empty, default arrow image will used.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("dev"),
			Label:     text.Chars(`Developer`),
			SortOrder: 920,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Backend::dev
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("debug"),
					Label:     text.Chars(`Debug`),
					SortOrder: 20,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/debug/template_hints_storefront
							ID:        path.NewRoute("template_hints_storefront"),
							Label:     text.Chars(`Enabled Template Path Hints for Storefront`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/debug/template_hints_admin
							ID:        path.NewRoute("template_hints_admin"),
							Label:     text.Chars(`Enabled Template Path Hints for Admin`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/debug/template_hints_blocks
							ID:        path.NewRoute("template_hints_blocks"),
							Label:     text.Chars(`Add Block Names to Hints`),
							Type:      element.TypeSelect,
							SortOrder: 21,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("template"),
					Label:     text.Chars(`Template Settings`),
					SortOrder: 25,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/template/allow_symlink
							ID:        path.NewRoute("allow_symlink"),
							Label:     text.Chars(`Allow Symlinks`),
							Comment:   text.Chars(`Warning! Enabling this feature is not recommended on production environments because it represents a potential security risk.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/template/minify_html
							ID:        path.NewRoute("minify_html"),
							Label:     text.Chars(`Minify Html`),
							Type:      element.TypeSelect,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("translate_inline"),
					Label:     text.Chars(`Translate Inline`),
					SortOrder: 30,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/translate_inline/active
							ID:        path.NewRoute("active"),
							Label:     text.Chars(`Enabled for Storefront`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/translate_inline/active_admin
							ID:        path.NewRoute("active_admin"),
							Label:     text.Chars(`Enabled for Admin`),
							Comment:   text.Chars(`Translate, blocks and other output caches should be disabled for both Storefront and Admin inline translations.`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Translate
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("js"),
					Label:     text.Chars(`JavaScript Settings`),
					SortOrder: 100,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/js/merge_files
							ID:        path.NewRoute("merge_files"),
							Label:     text.Chars(`Merge JavaScript Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/js/enable_js_bundling
							ID:        path.NewRoute("enable_js_bundling"),
							Label:     text.Chars(`Enable JavaScript Bundling`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/js/minify_files
							ID:        path.NewRoute("minify_files"),
							Label:     text.Chars(`Minify JavaScript Files`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("css"),
					Label:     text.Chars(`CSS Settings`),
					SortOrder: 110,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/css/merge_css_files
							ID:        path.NewRoute("merge_css_files"),
							Label:     text.Chars(`Merge CSS Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: dev/css/minify_files
							ID:        path.NewRoute("minify_files"),
							Label:     text.Chars(`Minify CSS Files`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("image"),
					Label:     text.Chars(`Image Processing Settings`),
					SortOrder: 120,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/image/default_adapter
							ID:        path.NewRoute("default_adapter"),
							Label:     text.Chars(`Image Adapter`),
							Comment:   text.Chars(`When the adapter was changed, please flush Catalog Images Cache.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Image\Adapter
							// SourceModel: Magento\Config\Model\Config\Source\Image\Adapter
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("static"),
					Label:     text.Chars(`Static Files Settings`),
					SortOrder: 130,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: dev/static/sign
							ID:        path.NewRoute("sign"),
							Label:     text.Chars(`Sign Static Files`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("general"),
			Label:     text.Chars(`General`),
			SortOrder: 10,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Config::config_general
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("store_information"),
					Label:     text.Chars(`Store Information`),
					SortOrder: 100,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/store_information/name
							ID:        path.NewRoute("name"),
							Label:     text.Chars(`Store Name`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: general/store_information/phone
							ID:        path.NewRoute("phone"),
							Label:     text.Chars(`Store Phone Number`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: general/store_information/hours
							ID:        path.NewRoute("hours"),
							Label:     text.Chars(`Store Hours of Operation`),
							Type:      element.TypeText,
							SortOrder: 22,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: general/store_information/country_id
							ID:         path.NewRoute("country_id"),
							Label:      text.Chars(`Country`),
							Type:       element.TypeSelect,
							SortOrder:  25,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: general/store_information/region_id
							ID:        path.NewRoute("region_id"),
							Label:     text.Chars(`Region/State`),
							Type:      element.TypeText,
							SortOrder: 27,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: general/store_information/postcode
							ID:        path.NewRoute("postcode"),
							Label:     text.Chars(`ZIP/Postal Code`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: general/store_information/city
							ID:        path.NewRoute("city"),
							Label:     text.Chars(`City`),
							Type:      element.TypeText,
							SortOrder: 45,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: general/store_information/street_line1
							ID:        path.NewRoute("street_line1"),
							Label:     text.Chars(`Street Address`),
							Type:      element.TypeText,
							SortOrder: 55,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: general/store_information/street_line2
							ID:        path.NewRoute("street_line2"),
							Label:     text.Chars(`Street Address Line 2`),
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: general/store_information/merchant_vat_number
							ID:         path.NewRoute("merchant_vat_number"),
							Label:      text.Chars(`VAT Number`),
							Type:       element.TypeText,
							SortOrder:  61,
							Visible:    element.VisibleYes,
							Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							CanBeEmpty: true,
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("single_store_mode"),
					Label:     text.Chars(`Single-Store Mode`),
					SortOrder: 150,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/single_store_mode/enabled
							ID:        path.NewRoute("enabled"),
							Label:     text.Chars(`Enable Single-Store Mode`),
							Comment:   text.Chars(`This setting will not be taken into account if system has more than one store view.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("system"),
			Label:     text.Chars(`System`),
			SortOrder: 900,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Config::config_system
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("smtp"),
					Label:     text.Chars(`Mail Sending Settings`),
					SortOrder: 20,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/smtp/disable
							ID:        path.NewRoute("disable"),
							Label:     text.Chars(`Disable Email Communications`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: system/smtp/host
							ID:        path.NewRoute("host"),
							Label:     text.Chars(`Host`),
							Comment:   text.Chars(`For Windows server only.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: system/smtp/port
							ID:        path.NewRoute("port"),
							Label:     text.Chars(`Port (25)`),
							Comment:   text.Chars(`For Windows server only.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: system/smtp/set_return_path
							ID:        path.NewRoute("set_return_path"),
							Label:     text.Chars(`Set Return-Path`),
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesnocustom
						},

						&element.Field{
							// Path: system/smtp/return_path_email
							ID:        path.NewRoute("return_path_email"),
							Label:     text.Chars(`Return-Path Email`),
							Type:      element.TypeText,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("admin"),
			Label:     text.Chars(`Admin`),
			SortOrder: 20,
			Scope:     scope.NewPerm(scope.DefaultID),
			Resource:  0, // Magento_Config::config_admin
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("emails"),
					Label:     text.Chars(`Admin User Emails`),
					SortOrder: 10,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/emails/forgot_email_template
							ID:        path.NewRoute("forgot_email_template"),
							Label:     text.Chars(`Forgot Password Email Template`),
							Comment:   text.Chars(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: admin/emails/forgot_email_identity
							ID:        path.NewRoute("forgot_email_identity"),
							Label:     text.Chars(`Forgot and Reset Email Sender`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: admin/emails/password_reset_link_expiration_period
							ID:        path.NewRoute("password_reset_link_expiration_period"),
							Label:     text.Chars(`Recovery Link Expiration Period (days)`),
							Comment:   text.Chars(`Please enter a number 1 or greater in this field.`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("startup"),
					Label:     text.Chars(`Startup Page`),
					SortOrder: 20,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/startup/menu_item_id
							ID:        path.NewRoute("menu_item_id"),
							Label:     text.Chars(`Startup Page`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Admin\Page
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("url"),
					Label:     text.Chars(`Admin Base URL`),
					SortOrder: 30,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/url/use_custom
							ID:        path.NewRoute("use_custom"),
							Label:     text.Chars(`Use Custom Admin URL`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usecustom
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/url/custom
							ID:        path.NewRoute("custom"),
							Label:     text.Chars(`Custom Admin URL`),
							Comment:   text.Chars(`Make sure that base URL ends with '/' (slash), e.g. http://yourdomain/magento/`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custom
						},

						&element.Field{
							// Path: admin/url/use_custom_path
							ID:        path.NewRoute("use_custom_path"),
							Label:     text.Chars(`Use Custom Admin Path`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/url/custom_path
							ID:        path.NewRoute("custom_path"),
							Label:     text.Chars(`Custom Admin Path`),
							Comment:   text.Chars(`You will have to sign in after you save your custom admin path.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("security"),
					Label:     text.Chars(`Security`),
					SortOrder: 35,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/security/use_form_key
							ID:        path.NewRoute("use_form_key"),
							Label:     text.Chars(`Add Secret Key to URLs`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usesecretkey
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/security/use_case_sensitive_login
							ID:        path.NewRoute("use_case_sensitive_login"),
							Label:     text.Chars(`Login is Case Sensitive`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: admin/security/session_lifetime
							ID:        path.NewRoute("session_lifetime"),
							Label:     text.Chars(`Admin Session Lifetime (seconds)`),
							Comment:   text.Chars(`Values less than 60 are ignored.`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("dashboard"),
					Label:     text.Chars(`Dashboard`),
					SortOrder: 40,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: admin/dashboard/enable_charts
							ID:        path.NewRoute("enable_charts"),
							Label:     text.Chars(`Enable Charts`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID:        path.NewRoute("web"),
			Label:     text.Chars(`Web`),
			SortOrder: 20,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Backend::web
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("url"),
					Label:     text.Chars(`Url Options`),
					SortOrder: 3,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/url/use_store
							ID:        path.NewRoute("use_store"),
							Label:     text.Chars(`Add Store Code to Urls`),
							Comment:   text.Chars(`<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   false,
							// BackendModel: Magento\Config\Model\Config\Backend\Store
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/url/redirect_to_base
							ID:        path.NewRoute("redirect_to_base"),
							Label:     text.Chars(`Auto-redirect to Base URL`),
							Comment:   text.Chars(`I.e. redirect from http://example.com/store/ to http://www.example.com/store/`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   1,
							// SourceModel: Magento\Config\Model\Config\Source\Web\Redirect
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("seo"),
					Label:     text.Chars(`Search Engine Optimization`),
					SortOrder: 5,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/seo/use_rewrites
							ID:        path.NewRoute("use_rewrites"),
							Label:     text.Chars(`Use Web Server Rewrites`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("unsecure"),
					Label:     text.Chars(`Base URLs`),
					Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
					SortOrder: 10,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/unsecure/base_url
							ID:        path.NewRoute("base_url"),
							Label:     text.Chars(`Base URL`),
							Comment:   text.Chars(`Specify URL or {{base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_link_url
							ID:        path.NewRoute("base_link_url"),
							Label:     text.Chars(`Base Link URL`),
							Comment:   text.Chars(`May start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_static_url
							ID:        path.NewRoute("base_static_url"),
							Label:     text.Chars(`Base URL for Static View Files`),
							Comment:   text.Chars(`May be empty or start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/unsecure/base_media_url
							ID:        path.NewRoute("base_media_url"),
							Label:     text.Chars(`Base URL for User Media Files`),
							Comment:   text.Chars(`May be empty or start with {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("secure"),
					Label:     text.Chars(`Base URLs (Secure)`),
					Comment:   text.Chars(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`),
					SortOrder: 20,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/secure/base_url
							ID:        path.NewRoute("base_url"),
							Label:     text.Chars(`Secure Base URL`),
							Comment:   text.Chars(`Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_link_url
							ID:        path.NewRoute("base_link_url"),
							Label:     text.Chars(`Secure Base Link URL`),
							Comment:   text.Chars(`May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_static_url
							ID:        path.NewRoute("base_static_url"),
							Label:     text.Chars(`Secure Base URL for Static View Files`),
							Comment:   text.Chars(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 25,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/base_media_url
							ID:        path.NewRoute("base_media_url"),
							Label:     text.Chars(`Secure Base URL for User Media Files`),
							Comment:   text.Chars(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
						},

						&element.Field{
							// Path: web/secure/use_in_frontend
							ID:        path.NewRoute("use_in_frontend"),
							Label:     text.Chars(`Use Secure URLs on Storefront`),
							Comment:   text.Chars(`Enter https protocol to use Secure URLs on Storefront.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/use_in_adminhtml
							ID:        path.NewRoute("use_in_adminhtml"),
							Label:     text.Chars(`Use Secure URLs in Admin`),
							Comment:   text.Chars(`Enter https protocol to use Secure URLs in Admin.`),
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/enable_hsts
							ID:        path.NewRoute("enable_hsts"),
							Label:     text.Chars(`Enable HTTP Strict Transport Security (HSTS)`),
							Comment:   text.Chars(`See <a href="https://www.owasp.org/index.php/HTTP_Strict_Transport_Security" target="_blank">HTTP Strict Transport Security</a> page for details.`),
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/enable_upgrade_insecure
							ID:        path.NewRoute("enable_upgrade_insecure"),
							Label:     text.Chars(`Upgrade Insecure Requests`),
							Comment:   text.Chars(`See <a href="http://www.w3.org/TR/upgrade-insecure-requests/" target="_blank">Upgrade Insecure Requests</a> page for details.`),
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Secure
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/secure/offloader_header
							ID:        path.NewRoute("offloader_header"),
							Label:     text.Chars(`Offloader header`),
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("default"),
					Label:     text.Chars(`Default Pages`),
					SortOrder: 30,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/default/front
							ID:        path.NewRoute("front"),
							Label:     text.Chars(`Default Web URL`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: web/default/no_route
							ID:        path.NewRoute("no_route"),
							Label:     text.Chars(`Default No-route URL`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("session"),
					Label:     text.Chars(`Session Validation Settings`),
					SortOrder: 60,
					Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: web/session/use_remote_addr
							ID:        path.NewRoute("use_remote_addr"),
							Label:     text.Chars(`Validate REMOTE_ADDR`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_via
							ID:        path.NewRoute("use_http_via"),
							Label:     text.Chars(`Validate HTTP_VIA`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_x_forwarded_for
							ID:        path.NewRoute("use_http_x_forwarded_for"),
							Label:     text.Chars(`Validate HTTP_X_FORWARDED_FOR`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_http_user_agent
							ID:        path.NewRoute("use_http_user_agent"),
							Label:     text.Chars(`Validate HTTP_USER_AGENT`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: web/session/use_frontend_sid
							ID:        path.NewRoute("use_frontend_sid"),
							Label:     text.Chars(`Use SID on Storefront`),
							Comment:   text.Chars(`Allows customers to stay logged in when switching between different stores.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: path.NewRoute("system"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("media_storage_configuration"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      path.NewRoute(`allowed_resources`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"email_folder":"email"}`,
						},
					),
				},

				&element.Group{
					ID: path.NewRoute("emails"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/emails/forgot_email_template
							ID:      path.NewRoute(`forgot_email_template`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `system_emails_forgot_email_template`,
						},

						&element.Field{
							// Path: system/emails/forgot_email_identity
							ID:      path.NewRoute(`forgot_email_identity`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `general`,
						},
					),
				},

				&element.Group{
					ID: path.NewRoute("dashboard"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/dashboard/enable_charts
							ID:      path.NewRoute(`enable_charts`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("general"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("validator_data"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/validator_data/input_types
							ID:      path.NewRoute(`input_types`),
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
