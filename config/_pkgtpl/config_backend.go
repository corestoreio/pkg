// +build ignore

package backend

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "advanced",
		Label:     `Advanced`,
		SortOrder: 910,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::advanced
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "modules_disable_output",
				Label:     `Disable Modules Output`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields:    config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID:        "trans_email",
		Label:     `Store Email Addresses`,
		SortOrder: 90,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::trans_email
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "ident_custom1",
				Label:     `Custom Email 1`,
				SortOrder: 4,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_custom1/email
						ID:        "email",
						Label:     `Sender Email`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: trans_email/ident_custom1/name
						ID:        "name",
						Label:     `Sender Name`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
					},
				),
			},

			&config.Group{
				ID:        "ident_custom2",
				Label:     `Custom Email 2`,
				SortOrder: 5,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_custom2/email
						ID:        "email",
						Label:     `Sender Email`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: trans_email/ident_custom2/name
						ID:        "name",
						Label:     `Sender Name`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
					},
				),
			},

			&config.Group{
				ID:        "ident_general",
				Label:     `General Contact`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_general/email
						ID:        "email",
						Label:     `Sender Email`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: trans_email/ident_general/name
						ID:        "name",
						Label:     `Sender Name`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
					},
				),
			},

			&config.Group{
				ID:        "ident_sales",
				Label:     `Sales Representative`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_sales/email
						ID:        "email",
						Label:     `Sender Email`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: trans_email/ident_sales/name
						ID:        "name",
						Label:     `Sender Name`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
					},
				),
			},

			&config.Group{
				ID:        "ident_support",
				Label:     `Customer Support`,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_support/email
						ID:        "email",
						Label:     `Sender Email`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: trans_email/ident_support/name
						ID:        "name",
						Label:     `Sender Name`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "design",
		Label:     `Design`,
		SortOrder: 30,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Config::config_design
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "theme",
				Label:     `Design Theme`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: design/theme/theme_id
						ID:        "theme_id",
						Label:     `Design Theme`,
						Comment:   element.LongText(`If no value is specified, the system default will be used. The system default may be modified by third party extensions.`),
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Theme\Model\Design\Backend\Theme
						// SourceModel: Otnegam\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
					},

					&config.Field{
						// Path: design/theme/ua_regexp
						ID:        "ua_regexp",
						Label:     `User-Agent Exceptions`,
						Comment:   element.LongText(`Search strings are either normal strings or regular exceptions (PCRE). They are matched in the same order as entered. Examples:<br /><span style="font-family:monospace">Firefox<br />/^mozilla/i</span>`),
						Tooltip:   element.LongText(`Find a string in client user-agent header and switch to specific design theme for that browser.`),
						Type:      config.Type,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// BackendModel: Otnegam\Theme\Model\Design\Backend\Exceptions
					},
				),
			},

			&config.Group{
				ID:        "pagination",
				Label:     `Pagination`,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: design/pagination/pagination_frame
						ID:        "pagination_frame",
						Label:     `Pagination Frame`,
						Comment:   element.LongText(`How many links to display at once.`),
						Type:      config.TypeText,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/pagination/pagination_frame_skip
						ID:        "pagination_frame_skip",
						Label:     `Pagination Frame Skip`,
						Comment:   element.LongText(`If the current frame position does not cover utmost pages, will render link to current position plus/minus this value.`),
						Type:      config.TypeText,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/pagination/anchor_text_for_previous
						ID:        "anchor_text_for_previous",
						Label:     `Anchor Text for Previous`,
						Comment:   element.LongText(`Alternative text for previous link in pagination menu. If empty, default arrow image will used.`),
						Type:      config.TypeText,
						SortOrder: 9,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/pagination/anchor_text_for_next
						ID:        "anchor_text_for_next",
						Label:     `Anchor Text for Next`,
						Comment:   element.LongText(`Alternative text for next link in pagination menu. If empty, default arrow image will used.`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "dev",
		Label:     `Developer`,
		SortOrder: 920,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::dev
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "debug",
				Label:     `Debug`,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/debug/template_hints_storefront
						ID:        "template_hints_storefront",
						Label:     `Enabled Template Path Hints for Storefront`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/debug/template_hints_admin
						ID:        "template_hints_admin",
						Label:     `Enabled Template Path Hints for Admin`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/debug/template_hints_blocks
						ID:        "template_hints_blocks",
						Label:     `Add Block Names to Hints`,
						Type:      config.TypeSelect,
						SortOrder: 21,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "template",
				Label:     `Template Settings`,
				SortOrder: 25,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/template/allow_symlink
						ID:        "allow_symlink",
						Label:     `Allow Symlinks`,
						Comment:   element.LongText(`Warning! Enabling this feature is not recommended on production environments because it represents a potential security risk.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/template/minify_html
						ID:        "minify_html",
						Label:     `Minify Html`,
						Type:      config.TypeSelect,
						SortOrder: 25,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "translate_inline",
				Label:     `Translate Inline`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/translate_inline/active
						ID:        "active",
						Label:     `Enabled for Storefront`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/translate_inline/active_admin
						ID:        "active_admin",
						Label:     `Enabled for Admin`,
						Comment:   element.LongText(`Translate, blocks and other output caches should be disabled for both Storefront and Admin inline translations.`),
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "js",
				Label:     `JavaScript Settings`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/js/merge_files
						ID:        "merge_files",
						Label:     `Merge JavaScript Files`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/js/enable_js_bundling
						ID:        "enable_js_bundling",
						Label:     `Enable JavaScript Bundling`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/js/minify_files
						ID:        "minify_files",
						Label:     `Minify JavaScript Files`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "css",
				Label:     `CSS Settings`,
				SortOrder: 110,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/css/merge_css_files
						ID:        "merge_css_files",
						Label:     `Merge CSS Files`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: dev/css/minify_files
						ID:        "minify_files",
						Label:     `Minify CSS Files`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "image",
				Label:     `Image Processing Settings`,
				SortOrder: 120,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/image/default_adapter
						ID:        "default_adapter",
						Label:     `Image Adapter`,
						Comment:   element.LongText(`When the adapter was changed, please flush Catalog Images Cache.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Adapter
						// SourceModel: Otnegam\Config\Model\Config\Source\Image\Adapter
					},
				),
			},

			&config.Group{
				ID:        "static",
				Label:     `Static Files Settings`,
				SortOrder: 130,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: dev/static/sign
						ID:        "sign",
						Label:     `Sign Static Files`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "general",
		Label:     `General`,
		SortOrder: 10,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Config::config_general
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "country",
				Label:     `Country Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/country/allow
						ID:         "allow",
						Label:      `Allow Countries`,
						Type:       config.TypeMultiselect,
						SortOrder:  2,
						Visible:    config.VisibleYes,
						Scope:      scope.PermAll,
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: general/country/default
						ID:        "default",
						Label:     `Default Country`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: general/country/eu_countries
						ID:        "eu_countries",
						Label:     `European Union Countries`,
						Type:      config.TypeMultiselect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: general/country/destinations
						ID:        "destinations",
						Label:     `Top destinations`,
						Type:      config.TypeMultiselect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},
				),
			},

			&config.Group{
				ID:        "locale",
				Label:     `Locale Options`,
				SortOrder: 8,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/locale/timezone
						ID:        "timezone",
						Label:     `Timezone`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Locale\Timezone
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Timezone
					},

					&config.Field{
						// Path: general/locale/code
						ID:        "code",
						Label:     `Locale`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale
					},

					&config.Field{
						// Path: general/locale/firstday
						ID:        "firstday",
						Label:     `First Day of Week`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
					},

					&config.Field{
						// Path: general/locale/weekend
						ID:         "weekend",
						Label:      `Weekend Days`,
						Type:       config.TypeMultiselect,
						SortOrder:  15,
						Visible:    config.VisibleYes,
						Scope:      scope.PermAll,
						CanBeEmpty: true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
					},
				),
			},

			&config.Group{
				ID:        "store_information",
				Label:     `Store Information`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/store_information/name
						ID:        "name",
						Label:     `Store Name`,
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: general/store_information/phone
						ID:        "phone",
						Label:     `Store Phone Number`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: general/store_information/hours
						ID:        "hours",
						Label:     `Store Hours of Operation`,
						Type:      config.TypeText,
						SortOrder: 22,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: general/store_information/country_id
						ID:         "country_id",
						Label:      `Country`,
						Type:       config.TypeSelect,
						SortOrder:  25,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: general/store_information/region_id
						ID:        "region_id",
						Label:     `Region/State`,
						Type:      config.TypeText,
						SortOrder: 27,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: general/store_information/postcode
						ID:        "postcode",
						Label:     `ZIP/Postal Code`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: general/store_information/city
						ID:        "city",
						Label:     `City`,
						Type:      config.TypeText,
						SortOrder: 45,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: general/store_information/street_line1
						ID:        "street_line1",
						Label:     `Street Address`,
						Type:      config.TypeText,
						SortOrder: 55,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: general/store_information/street_line2
						ID:        "street_line2",
						Label:     `Street Address Line 2`,
						Type:      config.TypeText,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: general/store_information/merchant_vat_number
						ID:         "merchant_vat_number",
						Label:      `VAT Number`,
						Type:       config.TypeText,
						SortOrder:  61,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						CanBeEmpty: true,
					},
				),
			},

			&config.Group{
				ID:        "single_store_mode",
				Label:     `Single-Store Mode`,
				SortOrder: 150,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/single_store_mode/enabled
						ID:        "enabled",
						Label:     `Enable Single-Store Mode`,
						Comment:   element.LongText(`This setting will not be taken into account if system has more than one store view.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "system",
		Label:     `System`,
		SortOrder: 900,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Config::config_system
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "smtp",
				Label:     `Mail Sending Settings`,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/smtp/disable
						ID:        "disable",
						Label:     `Disable Email Communications`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: system/smtp/host
						ID:        "host",
						Label:     `Host`,
						Comment:   element.LongText(`For Windows server only.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: system/smtp/port
						ID:        "port",
						Label:     `Port (25)`,
						Comment:   element.LongText(`For Windows server only.`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: system/smtp/set_return_path
						ID:        "set_return_path",
						Label:     `Set Return-Path`,
						Type:      config.TypeSelect,
						SortOrder: 70,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesnocustom
					},

					&config.Field{
						// Path: system/smtp/return_path_email
						ID:        "return_path_email",
						Label:     `Return-Path Email`,
						Type:      config.TypeText,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "admin",
		Label:     `Admin`,
		SortOrder: 20,
		Scope:     scope.NewPerm(scope.DefaultID),
		Resource:  0, // Otnegam_Config::config_admin
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "emails",
				Label:     `Admin User Emails`,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/emails/forgot_email_template
						ID:        "forgot_email_template",
						Label:     `Forgot Password Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: admin/emails/forgot_email_identity
						ID:        "forgot_email_identity",
						Label:     `Forgot and Reset Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: admin/emails/password_reset_link_expiration_period
						ID:        "password_reset_link_expiration_period",
						Label:     `Recovery Link Expiration Period (days)`,
						Comment:   element.LongText(`Please enter a number 1 or greater in this field.`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
					},
				),
			},

			&config.Group{
				ID:        "startup",
				Label:     `Startup Page`,
				SortOrder: 20,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/startup/menu_item_id
						ID:        "menu_item_id",
						Label:     `Startup Page`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Admin\Page
					},
				),
			},

			&config.Group{
				ID:        "url",
				Label:     `Admin Base URL`,
				SortOrder: 30,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/url/use_custom
						ID:        "use_custom",
						Label:     `Use Custom Admin URL`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usecustom
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: admin/url/custom
						ID:        "custom",
						Label:     `Custom Admin URL`,
						Comment:   element.LongText(`Make sure that base URL ends with '/' (slash), e.g. http://yourdomain/magento/`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custom
					},

					&config.Field{
						// Path: admin/url/use_custom_path
						ID:        "use_custom_path",
						Label:     `Use Custom Admin Path`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: admin/url/custom_path
						ID:        "custom_path",
						Label:     `Custom Admin Path`,
						Comment:   element.LongText(`You will have to sign in after you save your custom admin path.`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
					},
				),
			},

			&config.Group{
				ID:        "security",
				Label:     `Security`,
				SortOrder: 35,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/security/use_form_key
						ID:        "use_form_key",
						Label:     `Add Secret Key to URLs`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usesecretkey
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: admin/security/use_case_sensitive_login
						ID:        "use_case_sensitive_login",
						Label:     `Login is Case Sensitive`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: admin/security/session_lifetime
						ID:        "session_lifetime",
						Label:     `Admin Session Lifetime (seconds)`,
						Comment:   element.LongText(`Values less than 60 are ignored.`),
						Type:      config.Type,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},

			&config.Group{
				ID:        "dashboard",
				Label:     `Dashboard`,
				SortOrder: 40,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: admin/dashboard/enable_charts
						ID:        "enable_charts",
						Label:     `Enable Charts`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID:        "web",
		Label:     `Web`,
		SortOrder: 20,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::web
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "url",
				Label:     `Url Options`,
				SortOrder: 3,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/url/use_store
						ID:        "use_store",
						Label:     `Add Store Code to Urls`,
						Comment:   element.LongText(`<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Store
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/url/redirect_to_base
						ID:        "redirect_to_base",
						Label:     `Auto-redirect to Base URL`,
						Comment:   element.LongText(`I.e. redirect from http://example.com/store/ to http://www.example.com/store/`),
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Web\Redirect
					},
				),
			},

			&config.Group{
				ID:        "seo",
				Label:     `Search Engine Optimization`,
				SortOrder: 5,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/seo/use_rewrites
						ID:        "use_rewrites",
						Label:     `Use Web Server Rewrites`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "unsecure",
				Label:     `Base URLs`,
				Comment:   element.LongText(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/unsecure/base_url
						ID:        "base_url",
						Label:     `Base URL`,
						Comment:   element.LongText(`Specify URL or {{base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/unsecure/base_link_url
						ID:        "base_link_url",
						Label:     `Base Link URL`,
						Comment:   element.LongText(`May start with {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/unsecure/base_static_url
						ID:        "base_static_url",
						Label:     `Base URL for Static View Files`,
						Comment:   element.LongText(`May be empty or start with {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 25,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/unsecure/base_media_url
						ID:        "base_media_url",
						Label:     `Base URL for User Media Files`,
						Comment:   element.LongText(`May be empty or start with {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},
				),
			},

			&config.Group{
				ID:        "secure",
				Label:     `Base URLs (Secure)`,
				Comment:   element.LongText(`Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`),
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/secure/base_url
						ID:        "base_url",
						Label:     `Secure Base URL`,
						Comment:   element.LongText(`Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/secure/base_link_url
						ID:        "base_link_url",
						Label:     `Secure Base Link URL`,
						Comment:   element.LongText(`May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/secure/base_static_url
						ID:        "base_static_url",
						Label:     `Secure Base URL for Static View Files`,
						Comment:   element.LongText(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 25,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/secure/base_media_url
						ID:        "base_media_url",
						Label:     `Secure Base URL for User Media Files`,
						Comment:   element.LongText(`May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`),
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
					},

					&config.Field{
						// Path: web/secure/use_in_frontend
						ID:        "use_in_frontend",
						Label:     `Use Secure URLs on Storefront`,
						Comment:   element.LongText(`Enter https protocol to use Secure URLs on Storefront.`),
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/secure/use_in_adminhtml
						ID:        "use_in_adminhtml",
						Label:     `Use Secure URLs in Admin`,
						Comment:   element.LongText(`Enter https protocol to use Secure URLs in Admin.`),
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/secure/enable_hsts
						ID:        "enable_hsts",
						Label:     `Enable HTTP Strict Transport Security (HSTS)`,
						Comment:   element.LongText(`See <a href="https://www.owasp.org/index.php/HTTP_Strict_Transport_Security" target="_blank">HTTP Strict Transport Security</a> page for details.`),
						Type:      config.TypeSelect,
						SortOrder: 70,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/secure/enable_upgrade_insecure
						ID:        "enable_upgrade_insecure",
						Label:     `Upgrade Insecure Requests`,
						Comment:   element.LongText(`See <a href="http://www.w3.org/TR/upgrade-insecure-requests/" target="_blank">Upgrade Insecure Requests</a> page for details.`),
						Type:      config.TypeSelect,
						SortOrder: 80,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/secure/offloader_header
						ID:        "offloader_header",
						Label:     `Offloader header`,
						Type:      config.TypeText,
						SortOrder: 90,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},

			&config.Group{
				ID:        "default",
				Label:     `Default Pages`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/default/front
						ID:        "front",
						Label:     `Default Web URL`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: web/default/no_route
						ID:        "no_route",
						Label:     `Default No-route URL`,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},

			&config.Group{
				ID:        "session",
				Label:     `Session Validation Settings`,
				SortOrder: 60,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: web/session/use_remote_addr
						ID:        "use_remote_addr",
						Label:     `Validate REMOTE_ADDR`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/session/use_http_via
						ID:        "use_http_via",
						Label:     `Validate HTTP_VIA`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/session/use_http_x_forwarded_for
						ID:        "use_http_x_forwarded_for",
						Label:     `Validate HTTP_X_FORWARDED_FOR`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/session/use_http_user_agent
						ID:        "use_http_user_agent",
						Label:     `Validate HTTP_USER_AGENT`,
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: web/session/use_frontend_sid
						ID:        "use_frontend_sid",
						Label:     `Use SID on Storefront`,
						Comment:   element.LongText(`Allows customers to stay logged in when switching between different stores.`),
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
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
						Default: `{"email_folder":"email"}`,
					},
				),
			},

			&config.Group{
				ID: "emails",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/emails/forgot_email_template
						ID:      `forgot_email_template`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `system_emails_forgot_email_template`,
					},

					&config.Field{
						// Path: system/emails/forgot_email_identity
						ID:      `forgot_email_identity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `general`,
					},
				),
			},

			&config.Group{
				ID: "dashboard",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/dashboard/enable_charts
						ID:      `enable_charts`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "general",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "validator_data",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/validator_data/input_types
						ID:      `input_types`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"price":"price","media_image":"media_image","gallery":"gallery"}`,
					},
				),
			},
		),
	},
)
