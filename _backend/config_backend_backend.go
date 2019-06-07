// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sync"

	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/config/source"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// TODO: during development move each of this config stuff into its own package.

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	sync.Mutex
	// TransEmailIdentCustom1Email => Sender Email.
	// Path: trans_email/ident_custom1/email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	TransEmailIdentCustom1Email cfgmodel.Str

	// TransEmailIdentCustom1Name => Sender Name.
	// Path: trans_email/ident_custom1/name
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentCustom1Name cfgmodel.Str

	// TransEmailIdentCustom2Email => Sender Email.
	// Path: trans_email/ident_custom2/email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	TransEmailIdentCustom2Email cfgmodel.Str

	// TransEmailIdentCustom2Name => Sender Name.
	// Path: trans_email/ident_custom2/name
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentCustom2Name cfgmodel.Str

	// TransEmailIdentGeneralEmail => Sender Email.
	// Path: trans_email/ident_general/email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	TransEmailIdentGeneralEmail cfgmodel.Str

	// TransEmailIdentGeneralName => Sender Name.
	// Path: trans_email/ident_general/name
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentGeneralName cfgmodel.Str

	// TransEmailIdentSalesEmail => Sender Email.
	// Path: trans_email/ident_sales/email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	TransEmailIdentSalesEmail cfgmodel.Str

	// TransEmailIdentSalesName => Sender Name.
	// Path: trans_email/ident_sales/name
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentSalesName cfgmodel.Str

	// TransEmailIdentSupportEmail => Sender Email.
	// Path: trans_email/ident_support/email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	TransEmailIdentSupportEmail cfgmodel.Str

	// TransEmailIdentSupportName => Sender Name.
	// Path: trans_email/ident_support/name
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Sender
	TransEmailIdentSupportName cfgmodel.Str

	// DesignThemeThemeID => Design Theme.
	// If no value is specified, the system default will be used. The system
	// default may be modified by third party extensions.
	// Path: design/theme/theme_id
	// BackendModel: Magento\Theme\Model\Design\Backend\Theme
	// SourceModel: Magento\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
	DesignThemeThemeID cfgmodel.Str

	// DesignThemeUaRegexp => User-Agent Exceptions.
	// Search strings are either normal strings or regular exceptions (PCRE). They
	// are matched in the same order as entered. Examples:Firefox/^mozilla/i
	// Path: design/theme/ua_regexp
	// BackendModel: Magento\Theme\Model\Design\Backend\Exceptions
	DesignThemeUaRegexp cfgmodel.Str

	// DesignPaginationPaginationFrame => Pagination Frame.
	// How many links to display at once.
	// Path: design/pagination/pagination_frame
	DesignPaginationPaginationFrame cfgmodel.Str

	// DesignPaginationPaginationFrameSkip => Pagination Frame Skip.
	// If the current frame position does not cover utmost pages, will render link
	// to current position plus/minus this value.
	// Path: design/pagination/pagination_frame_skip
	DesignPaginationPaginationFrameSkip cfgmodel.Str

	// DesignPaginationAnchorTextForPrevious => Anchor Text for Previous.
	// Alternative text for previous link in pagination menu. If empty, default
	// arrow image will used.
	// Path: design/pagination/anchor_text_for_previous
	DesignPaginationAnchorTextForPrevious cfgmodel.Str

	// DesignPaginationAnchorTextForNext => Anchor Text for Next.
	// Alternative text for next link in pagination menu. If empty, default arrow
	// image will used.
	// Path: design/pagination/anchor_text_for_next
	DesignPaginationAnchorTextForNext cfgmodel.Str

	// DevDebugTemplateHintsStorefront => Enabled Template Path Hints for Storefront.
	// Path: dev/debug/template_hints_storefront
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsStorefront cfgmodel.Bool

	// DevDebugTemplateHintsAdmin => Enabled Template Path Hints for Admin.
	// Path: dev/debug/template_hints_admin
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsAdmin cfgmodel.Bool

	// DevDebugTemplateHintsBlocks => Add Block Names to Hints.
	// Path: dev/debug/template_hints_blocks
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevDebugTemplateHintsBlocks cfgmodel.Bool

	// DevTemplateAllowSymlink => Allow Symlinks.
	// Warning! Enabling this feature is not recommended on production
	// environments because it represents a potential security risk.
	// Path: dev/template/allow_symlink
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevTemplateAllowSymlink cfgmodel.Bool

	// DevTemplateMinifyHTML => Minify HTML.
	// Path: dev/template/minify_html
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevTemplateMinifyHTML cfgmodel.Bool

	// DevTranslateInlineActive => Enabled for Storefront.
	// Path: dev/translate_inline/active
	// BackendModel: Magento\Config\Model\Config\Backend\Translate
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevTranslateInlineActive cfgmodel.Bool

	// DevTranslateInlineActiveAdmin => Enabled for Admin.
	// Translate, blocks and other output caches should be disabled for both
	// Storefront and Admin inline translations.
	// Path: dev/translate_inline/active_admin
	// BackendModel: Magento\Config\Model\Config\Backend\Translate
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevTranslateInlineActiveAdmin cfgmodel.Bool

	// DevJsMergeFiles => Merge JavaScript Files.
	// Path: dev/js/merge_files
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevJsMergeFiles cfgmodel.Bool

	// DevJsEnableJsBundling => Enable JavaScript Bundling.
	// Path: dev/js/enable_js_bundling
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevJsEnableJsBundling cfgmodel.Bool

	// DevJsMinifyFiles => Minify JavaScript Files.
	// Path: dev/js/minify_files
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevJsMinifyFiles cfgmodel.Bool

	// DevCSSMergeCSSFiles => Merge CSS Files.
	// Path: dev/css/merge_css_files
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevCSSMergeCSSFiles cfgmodel.Bool

	// DevCSSMinifyFiles => Minify CSS Files.
	// Path: dev/css/minify_files
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevCSSMinifyFiles cfgmodel.Bool

	// DevImageDefaultAdapter => Image Adapter.
	// When the adapter was changed, please flush Catalog Images Cache.
	// Path: dev/image/default_adapter
	// BackendModel: Magento\Config\Model\Config\Backend\Image\Adapter
	// SourceModel: Magento\Config\Model\Config\Source\Image\Adapter
	DevImageDefaultAdapter cfgmodel.Str

	// DevStaticSign => Sign Static Files.
	// Path: dev/static/sign
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	DevStaticSign cfgmodel.Bool

	// SystemSMTPDisable => Disable Email Communications.
	// Path: system/smtp/disable
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemSMTPDisable cfgmodel.Bool

	// SystemSMTPHost => Host.
	// For Windows server only.
	// Path: system/smtp/host
	SystemSMTPHost cfgmodel.Str

	// SystemSMTPPort => Port (25).
	// For Windows server only.
	// Path: system/smtp/port
	SystemSMTPPort cfgmodel.Str

	// SystemSMTPSetReturnPath => Set Return-Path.
	// Path: system/smtp/set_return_path
	// SourceModel: Magento\Config\Model\Config\Source\Yesnocustom
	SystemSMTPSetReturnPath cfgmodel.Bool

	// SystemSMTPReturnPathEmail => Return-Path Email.
	// Path: system/smtp/return_path_email
	// BackendModel: Magento\Config\Model\Config\Backend\Email\Address
	SystemSMTPReturnPathEmail cfgmodel.Str

	// AdminEmailsForgotEmailTemplate => Forgot Password Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: admin/emails/forgot_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	AdminEmailsForgotEmailTemplate cfgmodel.Str

	// AdminEmailsForgotEmailIdentity => Forgot and Reset Email Sender.
	// Path: admin/emails/forgot_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	AdminEmailsForgotEmailIdentity cfgmodel.Str

	// AdminEmailsPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
	// Please enter a number 1 or greater in this field.
	// Path: admin/emails/password_reset_link_expiration_period
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
	AdminEmailsPasswordResetLinkExpirationPeriod cfgmodel.Str

	// AdminStartupMenuItemID => Startup Page.
	// Path: admin/startup/menu_item_id
	// SourceModel: Magento\Config\Model\Config\Source\Admin\Page
	AdminStartupMenuItemID cfgmodel.Str

	// AdminURLUseCustom => Use Custom Admin URL.
	// Path: admin/url/use_custom
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usecustom
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminURLUseCustom cfgmodel.Bool

	// AdminURLCustom => Custom Admin URL.
	// Make sure that base URL ends with '/' (slash), e.g.
	// http://yourdomain/magento/
	// Path: admin/url/custom
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custom
	AdminURLCustom cfgmodel.Str

	// AdminURLUseCustomPath => Use Custom Admin Path.
	// Path: admin/url/use_custom_path
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminURLUseCustomPath cfgmodel.Bool

	// AdminURLCustomPath => Custom Admin Path.
	// You will have to sign in after you save your custom admin path.
	// Path: admin/url/custom_path
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Custompath
	AdminURLCustomPath cfgmodel.Str

	// AdminSecurityUseFormKey => Add Secret Key to URLs.
	// Path: admin/security/use_form_key
	// BackendModel: Magento\Config\Model\Config\Backend\Admin\Usesecretkey
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminSecurityUseFormKey cfgmodel.Bool

	// AdminSecurityUseCaseSensitiveLogin => Login is Case Sensitive.
	// Path: admin/security/use_case_sensitive_login
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminSecurityUseCaseSensitiveLogin cfgmodel.Bool

	// AdminSecuritySessionLifetime => Admin Session Lifetime (seconds).
	// Values less than 60 are ignored.
	// Path: admin/security/session_lifetime
	AdminSecuritySessionLifetime cfgmodel.Str

	// AdminDashboardEnableCharts => Enable Charts.
	// Path: admin/dashboard/enable_charts
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	AdminDashboardEnableCharts cfgmodel.Bool

	// WebURLUseStore => Add Store Code to URLs.
	// Warning! When using Store Code in URLs, in some cases system may not work
	// properly if URLs without Store Codes are specified in the third party
	// services (e.g. PayPal etc.).
	// Path: web/url/use_store
	// BackendModel: Magento\Config\Model\Config\Backend\Store
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebURLUseStore cfgmodel.Bool

	// WebURLRedirectToBase => Auto-redirect to Base URL.
	// I.e. redirect from http://example.com/store/ to
	// http://www.example.com/store/
	// Path: web/url/redirect_to_base
	// SourceModel: Magento\Config\Model\Config\Source\Web\Redirect
	WebURLRedirectToBase ConfigRedirectToBase

	// WebSeoUseRewrites => Use Web Server Rewrites.
	// Path: web/seo/use_rewrites
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSeoUseRewrites cfgmodel.Bool

	// WebUnsecureBaseURL => Base URL.
	// Specify URL or {{base_url}} placeholder.
	// Path: web/unsecure/base_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseURL cfgmodel.BaseURL

	// WebUnsecureBaseLinkURL => Base Link URL.
	// May start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_link_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseLinkURL cfgmodel.BaseURL

	// WebUnsecureBaseStaticURL => Base URL for Static View Files.
	// May be empty or start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_static_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseStaticURL cfgmodel.BaseURL

	// WebUnsecureBaseMediaURL => Base URL for User Media Files.
	// May be empty or start with {{unsecure_base_url}} placeholder.
	// Path: web/unsecure/base_media_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebUnsecureBaseMediaURL cfgmodel.BaseURL

	// WebSecureBaseURL => Secure Base URL.
	// Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.
	// Path: web/secure/base_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebSecureBaseURL cfgmodel.BaseURL

	// WebSecureBaseLinkURL => Secure Base Link URL.
	// May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.
	// Path: web/secure/base_link_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebSecureBaseLinkURL cfgmodel.BaseURL

	// WebSecureBaseStaticURL => Secure Base URL for Static View Files.
	// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
	// placeholder.
	// Path: web/secure/base_static_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebSecureBaseStaticURL cfgmodel.BaseURL

	// WebSecureBaseMediaURL => Secure Base URL for User Media Files.
	// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
	// placeholder.
	// Path: web/secure/base_media_url
	// BackendModel: Magento\Config\Model\Config\Backend\Baseurl
	WebSecureBaseMediaURL cfgmodel.BaseURL

	// WebSecureUseInFrontend => Use Secure URLs on Storefront.
	// Enter https protocol to use Secure URLs on Storefront.
	// Path: web/secure/use_in_frontend
	// BackendModel: Magento\Config\Model\Config\Backend\Secure
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSecureUseInFrontend cfgmodel.Bool

	// WebSecureUseInAdminhtml => Use Secure URLs in Admin.
	// Enter https protocol to use Secure URLs in Admin.
	// Path: web/secure/use_in_adminhtml
	// BackendModel: Magento\Config\Model\Config\Backend\Secure
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSecureUseInAdminhtml cfgmodel.Bool

	// WebSecureEnableHsts => Enable HTTP Strict Transport Security (HSTS).
	// See HTTP Strict Transport Security page for details.
	// Path: web/secure/enable_hsts
	// BackendModel: Magento\Config\Model\Config\Backend\Secure
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSecureEnableHsts cfgmodel.Bool

	// WebSecureEnableUpgradeInsecure => Upgrade Insecure Requests.
	// See Upgrade Insecure Requests page for details.
	// Path: web/secure/enable_upgrade_insecure
	// BackendModel: Magento\Config\Model\Config\Backend\Secure
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSecureEnableUpgradeInsecure cfgmodel.Bool

	// WebSecureOffloaderHeader => Offloader header.
	// Path: web/secure/offloader_header
	WebSecureOffloaderHeader cfgmodel.Str

	// WebDefaultFront => Default Web URL.
	// Path: web/default/front
	WebDefaultFront cfgmodel.Str

	// WebDefaultNoRoute => Default No-route URL.
	// Path: web/default/no_route
	WebDefaultNoRoute cfgmodel.Str

	// WebSessionUseRemoteAddr => Validate REMOTE_ADDR.
	// Path: web/session/use_remote_addr
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSessionUseRemoteAddr cfgmodel.Bool

	// WebSessionUseHTTPVia => Validate HTTP_VIA.
	// Path: web/session/use_http_via
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSessionUseHTTPVia cfgmodel.Bool

	// WebSessionUseHTTPXForwardedFor => Validate HTTP_X_FORWARDED_FOR.
	// Path: web/session/use_http_x_forwarded_for
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSessionUseHTTPXForwardedFor cfgmodel.Bool

	// WebSessionUseHTTPUserAgent => Validate HTTP_USER_AGENT.
	// Path: web/session/use_http_user_agent
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSessionUseHTTPUserAgent cfgmodel.Bool

	// WebSessionUseFrontendSid => Use SID on Storefront.
	// Allows customers to stay logged in when switching between different stores.
	// Path: web/session/use_frontend_sid
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	WebSessionUseFrontendSid cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.TransEmailIdentCustom1Email = cfgmodel.NewStr(`trans_email/ident_custom1/email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentCustom1Name = cfgmodel.NewStr(`trans_email/ident_custom1/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentCustom2Email = cfgmodel.NewStr(`trans_email/ident_custom2/email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentCustom2Name = cfgmodel.NewStr(`trans_email/ident_custom2/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentGeneralEmail = cfgmodel.NewStr(`trans_email/ident_general/email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentGeneralName = cfgmodel.NewStr(`trans_email/ident_general/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentSalesEmail = cfgmodel.NewStr(`trans_email/ident_sales/email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentSalesName = cfgmodel.NewStr(`trans_email/ident_sales/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentSupportEmail = cfgmodel.NewStr(`trans_email/ident_support/email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.TransEmailIdentSupportName = cfgmodel.NewStr(`trans_email/ident_support/name`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignThemeThemeID = cfgmodel.NewStr(`design/theme/theme_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignThemeUaRegexp = cfgmodel.NewStr(`design/theme/ua_regexp`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignPaginationPaginationFrame = cfgmodel.NewStr(`design/pagination/pagination_frame`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignPaginationPaginationFrameSkip = cfgmodel.NewStr(`design/pagination/pagination_frame_skip`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignPaginationAnchorTextForPrevious = cfgmodel.NewStr(`design/pagination/anchor_text_for_previous`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignPaginationAnchorTextForNext = cfgmodel.NewStr(`design/pagination/anchor_text_for_next`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DevDebugTemplateHintsStorefront = cfgmodel.NewBool(`dev/debug/template_hints_storefront`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevDebugTemplateHintsAdmin = cfgmodel.NewBool(`dev/debug/template_hints_admin`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevDebugTemplateHintsBlocks = cfgmodel.NewBool(`dev/debug/template_hints_blocks`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevTemplateAllowSymlink = cfgmodel.NewBool(`dev/template/allow_symlink`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevTemplateMinifyHTML = cfgmodel.NewBool(`dev/template/minify_html`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevTranslateInlineActive = cfgmodel.NewBool(`dev/translate_inline/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevTranslateInlineActiveAdmin = cfgmodel.NewBool(`dev/translate_inline/active_admin`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevJsMergeFiles = cfgmodel.NewBool(`dev/js/merge_files`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevJsEnableJsBundling = cfgmodel.NewBool(`dev/js/enable_js_bundling`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevJsMinifyFiles = cfgmodel.NewBool(`dev/js/minify_files`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevCSSMergeCSSFiles = cfgmodel.NewBool(`dev/css/merge_css_files`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevCSSMinifyFiles = cfgmodel.NewBool(`dev/css/minify_files`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.DevImageDefaultAdapter = cfgmodel.NewStr(`dev/image/default_adapter`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.DevStaticSign = cfgmodel.NewBool(`dev/static/sign`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.SystemSMTPDisable = cfgmodel.NewBool(`system/smtp/disable`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.EnableDisable))
	pp.SystemSMTPHost = cfgmodel.NewStr(`system/smtp/host`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemSMTPPort = cfgmodel.NewStr(`system/smtp/port`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemSMTPSetReturnPath = cfgmodel.NewBool(`system/smtp/set_return_path`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.SystemSMTPReturnPathEmail = cfgmodel.NewStr(`system/smtp/return_path_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminEmailsForgotEmailTemplate = cfgmodel.NewStr(`admin/emails/forgot_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminEmailsForgotEmailIdentity = cfgmodel.NewStr(`admin/emails/forgot_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminEmailsPasswordResetLinkExpirationPeriod = cfgmodel.NewStr(`admin/emails/password_reset_link_expiration_period`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminStartupMenuItemID = cfgmodel.NewStr(`admin/startup/menu_item_id`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminURLUseCustom = cfgmodel.NewBool(`admin/url/use_custom`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.AdminURLCustom = cfgmodel.NewStr(`admin/url/custom`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminURLUseCustomPath = cfgmodel.NewBool(`admin/url/use_custom_path`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.AdminURLCustomPath = cfgmodel.NewStr(`admin/url/custom_path`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityUseFormKey = cfgmodel.NewBool(`admin/security/use_form_key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminSecurityUseCaseSensitiveLogin = cfgmodel.NewBool(`admin/security/use_case_sensitive_login`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.AdminSecuritySessionLifetime = cfgmodel.NewStr(`admin/security/session_lifetime`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.AdminDashboardEnableCharts = cfgmodel.NewBool(`admin/dashboard/enable_charts`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebURLUseStore = cfgmodel.NewBool(`web/url/use_store`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebURLRedirectToBase = NewConfigRedirectToBase(`web/url/redirect_to_base`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSeoUseRewrites = cfgmodel.NewBool(`web/seo/use_rewrites`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebUnsecureBaseURL = cfgmodel.NewBaseURL(`web/unsecure/base_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebUnsecureBaseLinkURL = cfgmodel.NewBaseURL(`web/unsecure/base_link_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebUnsecureBaseStaticURL = cfgmodel.NewBaseURL(`web/unsecure/base_static_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebUnsecureBaseMediaURL = cfgmodel.NewBaseURL(`web/unsecure/base_media_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureBaseURL = cfgmodel.NewBaseURL(`web/secure/base_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureBaseLinkURL = cfgmodel.NewBaseURL(`web/secure/base_link_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureBaseStaticURL = cfgmodel.NewBaseURL(`web/secure/base_static_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureBaseMediaURL = cfgmodel.NewBaseURL(`web/secure/base_media_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureUseInFrontend = cfgmodel.NewBool(`web/secure/use_in_frontend`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSecureUseInAdminhtml = cfgmodel.NewBool(`web/secure/use_in_adminhtml`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSecureEnableHsts = cfgmodel.NewBool(`web/secure/enable_hsts`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureEnableUpgradeInsecure = cfgmodel.NewBool(`web/secure/enable_upgrade_insecure`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSecureOffloaderHeader = cfgmodel.NewStr(`web/secure/offloader_header`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebDefaultFront = cfgmodel.NewStr(`web/default/front`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebDefaultNoRoute = cfgmodel.NewStr(`web/default/no_route`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.WebSessionUseRemoteAddr = cfgmodel.NewBool(`web/session/use_remote_addr`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSessionUseHTTPVia = cfgmodel.NewBool(`web/session/use_http_via`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSessionUseHTTPXForwardedFor = cfgmodel.NewBool(`web/session/use_http_x_forwarded_for`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSessionUseHTTPUserAgent = cfgmodel.NewBool(`web/session/use_http_user_agent`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))
	pp.WebSessionUseFrontendSid = cfgmodel.NewBool(`web/session/use_frontend_sid`, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))

	return pp
}
