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
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// TransEmailIdentCustom1Email => Sender Email.
	// Path: trans_email/ident_custom1/email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	TransEmailIdentCustom1Email model.Str

	// TransEmailIdentCustom1Name => Sender Name.
	// Path: trans_email/ident_custom1/name
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentCustom1Name model.Str

	// TransEmailIdentCustom2Email => Sender Email.
	// Path: trans_email/ident_custom2/email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	TransEmailIdentCustom2Email model.Str

	// TransEmailIdentCustom2Name => Sender Name.
	// Path: trans_email/ident_custom2/name
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentCustom2Name model.Str

	// TransEmailIdentGeneralEmail => Sender Email.
	// Path: trans_email/ident_general/email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	TransEmailIdentGeneralEmail model.Str

	// TransEmailIdentGeneralName => Sender Name.
	// Path: trans_email/ident_general/name
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentGeneralName model.Str

	// TransEmailIdentSalesEmail => Sender Email.
	// Path: trans_email/ident_sales/email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	TransEmailIdentSalesEmail model.Str

	// TransEmailIdentSalesName => Sender Name.
	// Path: trans_email/ident_sales/name
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentSalesName model.Str

	// TransEmailIdentSupportEmail => Sender Email.
	// Path: trans_email/ident_support/email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	TransEmailIdentSupportEmail model.Str

	// TransEmailIdentSupportName => Sender Name.
	// Path: trans_email/ident_support/name
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentSupportName model.Str

	// DesignThemeThemeId => Design Theme.
	// If no value is specified, the system default will be used. The system
	// default may be modified by third party extensions.
	// Path: design/theme/theme_id
	// BackendModel: Otnegam\Theme\Model\Design\Backend\Theme
	// SourceModel: Otnegam\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
	DesignThemeThemeId model.Str

	// DesignThemeUaRegexp => User-Agent Exceptions.
	// Search strings are either normal strings or regular exceptions (PCRE). They
	// are matched in the same order as entered. Examples:Firefox/^mozilla/i
	// Path: design/theme/ua_regexp
	// BackendModel: Otnegam\Theme\Model\Design\Backend\Exceptions
	DesignThemeUaRegexp model.Str

	// DesignPaginationPaginationFrame => Pagination Frame.
	// How many links to display at once.
	// Path: design/pagination/pagination_frame
	DesignPaginationPaginationFrame model.Str

	// DesignPaginationPaginationFrameSkip => Pagination Frame Skip.
	// If the current frame position does not cover utmost pages, will render link
	// to current position plus/minus this value.
	// Path: design/pagination/pagination_frame_skip
	DesignPaginationPaginationFrameSkip model.Str

	// DesignPaginationAnchorTextForPrevious => Anchor Text for Previous.
	// Alternative text for previous link in pagination menu. If empty, default
	// arrow image will used.
	// Path: design/pagination/anchor_text_for_previous
	DesignPaginationAnchorTextForPrevious model.Str

	// DesignPaginationAnchorTextForNext => Anchor Text for Next.
	// Alternative text for next link in pagination menu. If empty, default arrow
	// image will used.
	// Path: design/pagination/anchor_text_for_next
	DesignPaginationAnchorTextForNext model.Str

	// DevDebugTemplateHintsStorefront => Enabled Template Path Hints for Storefront.
	// Path: dev/debug/template_hints_storefront
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsStorefront model.Bool

	// DevDebugTemplateHintsAdmin => Enabled Template Path Hints for Admin.
	// Path: dev/debug/template_hints_admin
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsAdmin model.Bool

	// DevDebugTemplateHintsBlocks => Add Block Names to Hints.
	// Path: dev/debug/template_hints_blocks
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsBlocks model.Bool

	// DevTemplateAllowSymlink => Allow Symlinks.
	// Warning! Enabling this feature is not recommended on production
	// environments because it represents a potential security risk.
	// Path: dev/template/allow_symlink
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevTemplateAllowSymlink model.Bool

	// DevTemplateMinifyHtml => Minify Html.
	// Path: dev/template/minify_html
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevTemplateMinifyHtml model.Bool

	// DevTranslateInlineActive => Enabled for Storefront.
	// Path: dev/translate_inline/active
	// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevTranslateInlineActive model.Bool

	// DevTranslateInlineActiveAdmin => Enabled for Admin.
	// Translate, blocks and other output caches should be disabled for both
	// Storefront and Admin inline translations.
	// Path: dev/translate_inline/active_admin
	// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevTranslateInlineActiveAdmin model.Bool

	// DevJsMergeFiles => Merge JavaScript Files.
	// Path: dev/js/merge_files
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevJsMergeFiles model.Bool

	// DevJsEnableJsBundling => Enable JavaScript Bundling.
	// Path: dev/js/enable_js_bundling
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevJsEnableJsBundling model.Bool

	// DevJsMinifyFiles => Minify JavaScript Files.
	// Path: dev/js/minify_files
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevJsMinifyFiles model.Bool

	// DevCssMergeCssFiles => Merge CSS Files.
	// Path: dev/css/merge_css_files
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevCssMergeCssFiles model.Bool

	// DevCssMinifyFiles => Minify CSS Files.
	// Path: dev/css/minify_files
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevCssMinifyFiles model.Bool

	// DevImageDefaultAdapter => Image Adapter.
	// When the adapter was changed, please flush Catalog Images Cache.
	// Path: dev/image/default_adapter
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Adapter
	// SourceModel: Otnegam\Config\Model\Config\Source\Image\Adapter
	DevImageDefaultAdapter model.Str

	// DevStaticSign => Sign Static Files.
	// Path: dev/static/sign
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	DevStaticSign model.Bool

	// GeneralStoreInformationName => Store Name.
	// Path: general/store_information/name
	GeneralStoreInformationName model.Str

	// GeneralStoreInformationPhone => Store Phone Number.
	// Path: general/store_information/phone
	GeneralStoreInformationPhone model.Str

	// GeneralStoreInformationHours => Store Hours of Operation.
	// Path: general/store_information/hours
	GeneralStoreInformationHours model.Str

	// GeneralStoreInformationCountryId => Country.
	// Path: general/store_information/country_id
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralStoreInformationCountryId model.Str

	// GeneralStoreInformationRegionId => Region/State.
	// Path: general/store_information/region_id
	GeneralStoreInformationRegionId model.Str

	// GeneralStoreInformationPostcode => ZIP/Postal Code.
	// Path: general/store_information/postcode
	GeneralStoreInformationPostcode model.Str

	// GeneralStoreInformationCity => City.
	// Path: general/store_information/city
	GeneralStoreInformationCity model.Str

	// GeneralStoreInformationStreetLine1 => Street Address.
	// Path: general/store_information/street_line1
	GeneralStoreInformationStreetLine1 model.Str

	// GeneralStoreInformationStreetLine2 => Street Address Line 2.
	// Path: general/store_information/street_line2
	GeneralStoreInformationStreetLine2 model.Str

	// GeneralStoreInformationMerchantVatNumber => VAT Number.
	// Path: general/store_information/merchant_vat_number
	GeneralStoreInformationMerchantVatNumber model.Str

	// GeneralSingleStoreModeEnabled => Enable Single-Store Mode.
	// This setting will not be taken into account if system has more than one
	// store view.
	// Path: general/single_store_mode/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	GeneralSingleStoreModeEnabled model.Bool

	// SystemSmtpDisable => Disable Email Communications.
	// Path: system/smtp/disable
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SystemSmtpDisable model.Bool

	// SystemSmtpHost => Host.
	// For Windows server only.
	// Path: system/smtp/host
	SystemSmtpHost model.Str

	// SystemSmtpPort => Port (25).
	// For Windows server only.
	// Path: system/smtp/port
	SystemSmtpPort model.Str

	// SystemSmtpSetReturnPath => Set Return-Path.
	// Path: system/smtp/set_return_path
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesnocustom
	SystemSmtpSetReturnPath model.Bool

	// SystemSmtpReturnPathEmail => Return-Path Email.
	// Path: system/smtp/return_path_email
	// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
	SystemSmtpReturnPathEmail model.Str

	// AdminEmailsForgotEmailTemplate => Forgot Password Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: admin/emails/forgot_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	AdminEmailsForgotEmailTemplate model.Str

	// AdminEmailsForgotEmailIdentity => Forgot and Reset Email Sender.
	// Path: admin/emails/forgot_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	AdminEmailsForgotEmailIdentity model.Str

	// AdminEmailsPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
	// Please enter a number 1 or greater in this field.
	// Path: admin/emails/password_reset_link_expiration_period
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
	AdminEmailsPasswordResetLinkExpirationPeriod model.Str

	// AdminStartupMenuItemId => Startup Page.
	// Path: admin/startup/menu_item_id
	// SourceModel: Otnegam\Config\Model\Config\Source\Admin\Page
	AdminStartupMenuItemId model.Str

	// AdminUrlUseCustom => Use Custom Admin URL.
	// Path: admin/url/use_custom
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usecustom
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminUrlUseCustom model.Bool

	// AdminUrlCustom => Custom Admin URL.
	// Make sure that base URL ends with '/' (slash), e.g.
	// http://yourdomain/magento/
	// Path: admin/url/custom
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custom
	AdminUrlCustom model.Str

	// AdminUrlUseCustomPath => Use Custom Admin Path.
	// Path: admin/url/use_custom_path
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminUrlUseCustomPath model.Bool

	// AdminUrlCustomPath => Custom Admin Path.
	// You will have to sign in after you save your custom admin path.
	// Path: admin/url/custom_path
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
	AdminUrlCustomPath model.Str

	// AdminSecurityUseFormKey => Add Secret Key to URLs.
	// Path: admin/security/use_form_key
	// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usesecretkey
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminSecurityUseFormKey model.Bool

	// AdminSecurityUseCaseSensitiveLogin => Login is Case Sensitive.
	// Path: admin/security/use_case_sensitive_login
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminSecurityUseCaseSensitiveLogin model.Bool

	// AdminSecuritySessionLifetime => Admin Session Lifetime (seconds).
	// Values less than 60 are ignored.
	// Path: admin/security/session_lifetime
	AdminSecuritySessionLifetime model.Str

	// AdminDashboardEnableCharts => Enable Charts.
	// Path: admin/dashboard/enable_charts
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	AdminDashboardEnableCharts model.Bool

	// WebUrlUseStore => Add Store Code to Urls.
	// Warning! When using Store Code in URLs, in some cases system may not work
	// properly if URLs without Store Codes are specified in the third party
	// services (e.g. PayPal etc.).
	// Path: web/url/use_store
	// BackendModel: Otnegam\Config\Model\Config\Backend\Store
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebUrlUseStore model.Bool

	// WebUrlRedirectToBase => Auto-redirect to Base URL.
	// I.e. redirect from http://example.com/store/ to
	// http://www.example.com/store/
	// Path: web/url/redirect_to_base
	// SourceModel: Otnegam\Config\Model\Config\Source\Web\Redirect
	WebUrlRedirectToBase ConfigRedirectToBase

	// WebSeoUseRewrites => Use Web Server Rewrites.
	// Path: web/seo/use_rewrites
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSeoUseRewrites model.Bool

	// WebUnsecureBaseUrl => Base URL.
	// Specify URL or {{base_url}} placeholder.
	// Path: web/unsecure/base_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseUrl model.Str

	// WebUnsecureBaseLinkUrl => Base Link URL.
	// May start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_link_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseLinkUrl model.Str

	// WebUnsecureBaseStaticUrl => Base URL for Static View Files.
	// May be empty or start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_static_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseStaticUrl model.Str

	// WebUnsecureBaseMediaUrl => Base URL for User Media Files.
	// May be empty or start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_media_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseMediaUrl model.Str

	// WebSecureBaseUrl => Secure Base URL.
	// Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.
	// Path: web/secure/base_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebSecureBaseUrl model.Str

	// WebSecureBaseLinkUrl => Secure Base Link URL.
	// May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.
	// Path: web/secure/base_link_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebSecureBaseLinkUrl model.Str

	// WebSecureBaseStaticUrl => Secure Base URL for Static View Files.
	// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
	// placeholder.
	// Path: web/secure/base_static_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebSecureBaseStaticUrl model.Str

	// WebSecureBaseMediaUrl => Secure Base URL for User Media Files.
	// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
	// placeholder.
	// Path: web/secure/base_media_url
	// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
	WebSecureBaseMediaUrl model.Str

	// WebSecureUseInFrontend => Use Secure URLs on Storefront.
	// Enter https protocol to use Secure URLs on Storefront.
	// Path: web/secure/use_in_frontend
	// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSecureUseInFrontend model.Bool

	// WebSecureUseInAdminhtml => Use Secure URLs in Admin.
	// Enter https protocol to use Secure URLs in Admin.
	// Path: web/secure/use_in_adminhtml
	// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSecureUseInAdminhtml model.Bool

	// WebSecureEnableHsts => Enable HTTP Strict Transport Security (HSTS).
	// See HTTP Strict Transport Security page for details.
	// Path: web/secure/enable_hsts
	// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSecureEnableHsts model.Bool

	// WebSecureEnableUpgradeInsecure => Upgrade Insecure Requests.
	// See Upgrade Insecure Requests page for details.
	// Path: web/secure/enable_upgrade_insecure
	// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSecureEnableUpgradeInsecure model.Bool

	// WebSecureOffloaderHeader => Offloader header.
	// Path: web/secure/offloader_header
	WebSecureOffloaderHeader model.Str

	// WebDefaultFront => Default Web URL.
	// Path: web/default/front
	WebDefaultFront model.Str

	// WebDefaultNoRoute => Default No-route URL.
	// Path: web/default/no_route
	WebDefaultNoRoute model.Str

	// WebSessionUseRemoteAddr => Validate REMOTE_ADDR.
	// Path: web/session/use_remote_addr
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSessionUseRemoteAddr model.Bool

	// WebSessionUseHttpVia => Validate HTTP_VIA.
	// Path: web/session/use_http_via
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSessionUseHttpVia model.Bool

	// WebSessionUseHttpXForwardedFor => Validate HTTP_X_FORWARDED_FOR.
	// Path: web/session/use_http_x_forwarded_for
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSessionUseHttpXForwardedFor model.Bool

	// WebSessionUseHttpUserAgent => Validate HTTP_USER_AGENT.
	// Path: web/session/use_http_user_agent
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSessionUseHttpUserAgent model.Bool

	// WebSessionUseFrontendSid => Use SID on Storefront.
	// Allows customers to stay logged in when switching between different stores.
	// Path: web/session/use_frontend_sid
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WebSessionUseFrontendSid model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.TransEmailIdentCustom1Email = model.NewStr(`trans_email/ident_custom1/email`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentCustom1Name = model.NewStr(`trans_email/ident_custom1/name`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentCustom2Email = model.NewStr(`trans_email/ident_custom2/email`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentCustom2Name = model.NewStr(`trans_email/ident_custom2/name`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentGeneralEmail = model.NewStr(`trans_email/ident_general/email`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentGeneralName = model.NewStr(`trans_email/ident_general/name`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentSalesEmail = model.NewStr(`trans_email/ident_sales/email`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentSalesName = model.NewStr(`trans_email/ident_sales/name`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentSupportEmail = model.NewStr(`trans_email/ident_support/email`, model.WithConfigStructure(cfgStruct))
	pp.TransEmailIdentSupportName = model.NewStr(`trans_email/ident_support/name`, model.WithConfigStructure(cfgStruct))
	pp.DesignThemeThemeId = model.NewStr(`design/theme/theme_id`, model.WithConfigStructure(cfgStruct))
	pp.DesignThemeUaRegexp = model.NewStr(`design/theme/ua_regexp`, model.WithConfigStructure(cfgStruct))
	pp.DesignPaginationPaginationFrame = model.NewStr(`design/pagination/pagination_frame`, model.WithConfigStructure(cfgStruct))
	pp.DesignPaginationPaginationFrameSkip = model.NewStr(`design/pagination/pagination_frame_skip`, model.WithConfigStructure(cfgStruct))
	pp.DesignPaginationAnchorTextForPrevious = model.NewStr(`design/pagination/anchor_text_for_previous`, model.WithConfigStructure(cfgStruct))
	pp.DesignPaginationAnchorTextForNext = model.NewStr(`design/pagination/anchor_text_for_next`, model.WithConfigStructure(cfgStruct))
	pp.DevDebugTemplateHintsStorefront = model.NewBool(`dev/debug/template_hints_storefront`, model.WithConfigStructure(cfgStruct))
	pp.DevDebugTemplateHintsAdmin = model.NewBool(`dev/debug/template_hints_admin`, model.WithConfigStructure(cfgStruct))
	pp.DevDebugTemplateHintsBlocks = model.NewBool(`dev/debug/template_hints_blocks`, model.WithConfigStructure(cfgStruct))
	pp.DevTemplateAllowSymlink = model.NewBool(`dev/template/allow_symlink`, model.WithConfigStructure(cfgStruct))
	pp.DevTemplateMinifyHtml = model.NewBool(`dev/template/minify_html`, model.WithConfigStructure(cfgStruct))
	pp.DevTranslateInlineActive = model.NewBool(`dev/translate_inline/active`, model.WithConfigStructure(cfgStruct))
	pp.DevTranslateInlineActiveAdmin = model.NewBool(`dev/translate_inline/active_admin`, model.WithConfigStructure(cfgStruct))
	pp.DevJsMergeFiles = model.NewBool(`dev/js/merge_files`, model.WithConfigStructure(cfgStruct))
	pp.DevJsEnableJsBundling = model.NewBool(`dev/js/enable_js_bundling`, model.WithConfigStructure(cfgStruct))
	pp.DevJsMinifyFiles = model.NewBool(`dev/js/minify_files`, model.WithConfigStructure(cfgStruct))
	pp.DevCssMergeCssFiles = model.NewBool(`dev/css/merge_css_files`, model.WithConfigStructure(cfgStruct))
	pp.DevCssMinifyFiles = model.NewBool(`dev/css/minify_files`, model.WithConfigStructure(cfgStruct))
	pp.DevImageDefaultAdapter = model.NewStr(`dev/image/default_adapter`, model.WithConfigStructure(cfgStruct))
	pp.DevStaticSign = model.NewBool(`dev/static/sign`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationName = model.NewStr(`general/store_information/name`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationPhone = model.NewStr(`general/store_information/phone`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationHours = model.NewStr(`general/store_information/hours`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationCountryId = model.NewStr(`general/store_information/country_id`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationRegionId = model.NewStr(`general/store_information/region_id`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationPostcode = model.NewStr(`general/store_information/postcode`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationCity = model.NewStr(`general/store_information/city`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationStreetLine1 = model.NewStr(`general/store_information/street_line1`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationStreetLine2 = model.NewStr(`general/store_information/street_line2`, model.WithConfigStructure(cfgStruct))
	pp.GeneralStoreInformationMerchantVatNumber = model.NewStr(`general/store_information/merchant_vat_number`, model.WithConfigStructure(cfgStruct))
	pp.GeneralSingleStoreModeEnabled = model.NewBool(`general/single_store_mode/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SystemSmtpDisable = model.NewBool(`system/smtp/disable`, model.WithConfigStructure(cfgStruct))
	pp.SystemSmtpHost = model.NewStr(`system/smtp/host`, model.WithConfigStructure(cfgStruct))
	pp.SystemSmtpPort = model.NewStr(`system/smtp/port`, model.WithConfigStructure(cfgStruct))
	pp.SystemSmtpSetReturnPath = model.NewBool(`system/smtp/set_return_path`, model.WithConfigStructure(cfgStruct))
	pp.SystemSmtpReturnPathEmail = model.NewStr(`system/smtp/return_path_email`, model.WithConfigStructure(cfgStruct))
	pp.AdminEmailsForgotEmailTemplate = model.NewStr(`admin/emails/forgot_email_template`, model.WithConfigStructure(cfgStruct))
	pp.AdminEmailsForgotEmailIdentity = model.NewStr(`admin/emails/forgot_email_identity`, model.WithConfigStructure(cfgStruct))
	pp.AdminEmailsPasswordResetLinkExpirationPeriod = model.NewStr(`admin/emails/password_reset_link_expiration_period`, model.WithConfigStructure(cfgStruct))
	pp.AdminStartupMenuItemId = model.NewStr(`admin/startup/menu_item_id`, model.WithConfigStructure(cfgStruct))
	pp.AdminUrlUseCustom = model.NewBool(`admin/url/use_custom`, model.WithConfigStructure(cfgStruct))
	pp.AdminUrlCustom = model.NewStr(`admin/url/custom`, model.WithConfigStructure(cfgStruct))
	pp.AdminUrlUseCustomPath = model.NewBool(`admin/url/use_custom_path`, model.WithConfigStructure(cfgStruct))
	pp.AdminUrlCustomPath = model.NewStr(`admin/url/custom_path`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityUseFormKey = model.NewBool(`admin/security/use_form_key`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecurityUseCaseSensitiveLogin = model.NewBool(`admin/security/use_case_sensitive_login`, model.WithConfigStructure(cfgStruct))
	pp.AdminSecuritySessionLifetime = model.NewStr(`admin/security/session_lifetime`, model.WithConfigStructure(cfgStruct))
	pp.AdminDashboardEnableCharts = model.NewBool(`admin/dashboard/enable_charts`, model.WithConfigStructure(cfgStruct))
	pp.WebUrlUseStore = model.NewBool(`web/url/use_store`, model.WithConfigStructure(cfgStruct))
	pp.WebUrlRedirectToBase = NewConfigRedirectToBase(`web/url/redirect_to_base`, model.WithConfigStructure(cfgStruct))
	pp.WebSeoUseRewrites = model.NewBool(`web/seo/use_rewrites`, model.WithConfigStructure(cfgStruct))
	pp.WebUnsecureBaseUrl = model.NewStr(`web/unsecure/base_url`, model.WithConfigStructure(cfgStruct))
	pp.WebUnsecureBaseLinkUrl = model.NewStr(`web/unsecure/base_link_url`, model.WithConfigStructure(cfgStruct))
	pp.WebUnsecureBaseStaticUrl = model.NewStr(`web/unsecure/base_static_url`, model.WithConfigStructure(cfgStruct))
	pp.WebUnsecureBaseMediaUrl = model.NewStr(`web/unsecure/base_media_url`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureBaseUrl = model.NewStr(`web/secure/base_url`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureBaseLinkUrl = model.NewStr(`web/secure/base_link_url`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureBaseStaticUrl = model.NewStr(`web/secure/base_static_url`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureBaseMediaUrl = model.NewStr(`web/secure/base_media_url`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureUseInFrontend = model.NewBool(`web/secure/use_in_frontend`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureUseInAdminhtml = model.NewBool(`web/secure/use_in_adminhtml`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureEnableHsts = model.NewBool(`web/secure/enable_hsts`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureEnableUpgradeInsecure = model.NewBool(`web/secure/enable_upgrade_insecure`, model.WithConfigStructure(cfgStruct))
	pp.WebSecureOffloaderHeader = model.NewStr(`web/secure/offloader_header`, model.WithConfigStructure(cfgStruct))
	pp.WebDefaultFront = model.NewStr(`web/default/front`, model.WithConfigStructure(cfgStruct))
	pp.WebDefaultNoRoute = model.NewStr(`web/default/no_route`, model.WithConfigStructure(cfgStruct))
	pp.WebSessionUseRemoteAddr = model.NewBool(`web/session/use_remote_addr`, model.WithConfigStructure(cfgStruct))
	pp.WebSessionUseHttpVia = model.NewBool(`web/session/use_http_via`, model.WithConfigStructure(cfgStruct))
	pp.WebSessionUseHttpXForwardedFor = model.NewBool(`web/session/use_http_x_forwarded_for`, model.WithConfigStructure(cfgStruct))
	pp.WebSessionUseHttpUserAgent = model.NewBool(`web/session/use_http_user_agent`, model.WithConfigStructure(cfgStruct))
	pp.WebSessionUseFrontendSid = model.NewBool(`web/session/use_frontend_sid`, model.WithConfigStructure(cfgStruct))

	return pp
}
