// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/csfw/config/model"
)

// PathTransEmailIdentCustom1Email => Sender Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathTransEmailIdentCustom1Email = model.NewStr(`trans_email/ident_custom1/email`)

// PathTransEmailIdentCustom1Name => Sender Name.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
var PathTransEmailIdentCustom1Name = model.NewStr(`trans_email/ident_custom1/name`)

// PathTransEmailIdentCustom2Email => Sender Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathTransEmailIdentCustom2Email = model.NewStr(`trans_email/ident_custom2/email`)

// PathTransEmailIdentCustom2Name => Sender Name.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
var PathTransEmailIdentCustom2Name = model.NewStr(`trans_email/ident_custom2/name`)

// PathTransEmailIdentGeneralEmail => Sender Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathTransEmailIdentGeneralEmail = model.NewStr(`trans_email/ident_general/email`)

// PathTransEmailIdentGeneralName => Sender Name.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
var PathTransEmailIdentGeneralName = model.NewStr(`trans_email/ident_general/name`)

// PathTransEmailIdentSalesEmail => Sender Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathTransEmailIdentSalesEmail = model.NewStr(`trans_email/ident_sales/email`)

// PathTransEmailIdentSalesName => Sender Name.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
var PathTransEmailIdentSalesName = model.NewStr(`trans_email/ident_sales/name`)

// PathTransEmailIdentSupportEmail => Sender Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathTransEmailIdentSupportEmail = model.NewStr(`trans_email/ident_support/email`)

// PathTransEmailIdentSupportName => Sender Name.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Sender
var PathTransEmailIdentSupportName = model.NewStr(`trans_email/ident_support/name`)

// PathDesignThemeThemeId => Design Theme.
// If no value is specified, the system default will be used. The system
// default may be modified by third party extensions.
// BackendModel: Otnegam\Theme\Model\Design\Backend\Theme
// SourceModel: Otnegam\Framework\View\Design\Theme\Label::getLabelsCollectionForSystemConfiguration
var PathDesignThemeThemeId = model.NewStr(`design/theme/theme_id`)

// PathDesignThemeUaRegexp => User-Agent Exceptions.
// Search strings are either normal strings or regular exceptions (PCRE). They
// are matched in the same order as entered. Examples:Firefox/^mozilla/i
// BackendModel: Otnegam\Theme\Model\Design\Backend\Exceptions
var PathDesignThemeUaRegexp = model.NewStr(`design/theme/ua_regexp`)

// PathDesignPaginationPaginationFrame => Pagination Frame.
// How many links to display at once.
var PathDesignPaginationPaginationFrame = model.NewStr(`design/pagination/pagination_frame`)

// PathDesignPaginationPaginationFrameSkip => Pagination Frame Skip.
// If the current frame position does not cover utmost pages, will render link
// to current position plus/minus this value.
var PathDesignPaginationPaginationFrameSkip = model.NewStr(`design/pagination/pagination_frame_skip`)

// PathDesignPaginationAnchorTextForPrevious => Anchor Text for Previous.
// Alternative text for previous link in pagination menu. If empty, default
// arrow image will used.
var PathDesignPaginationAnchorTextForPrevious = model.NewStr(`design/pagination/anchor_text_for_previous`)

// PathDesignPaginationAnchorTextForNext => Anchor Text for Next.
// Alternative text for next link in pagination menu. If empty, default arrow
// image will used.
var PathDesignPaginationAnchorTextForNext = model.NewStr(`design/pagination/anchor_text_for_next`)

// PathDevDebugTemplateHintsStorefront => Enabled Template Path Hints for Storefront.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevDebugTemplateHintsStorefront = model.NewBool(`dev/debug/template_hints_storefront`)

// PathDevDebugTemplateHintsAdmin => Enabled Template Path Hints for Admin.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevDebugTemplateHintsAdmin = model.NewBool(`dev/debug/template_hints_admin`)

// PathDevDebugTemplateHintsBlocks => Add Block Names to Hints.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevDebugTemplateHintsBlocks = model.NewBool(`dev/debug/template_hints_blocks`)

// PathDevTemplateAllowSymlink => Allow Symlinks.
// Warning! Enabling this feature is not recommended on production
// environments because it represents a potential security risk.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevTemplateAllowSymlink = model.NewBool(`dev/template/allow_symlink`)

// PathDevTemplateMinifyHtml => Minify Html.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevTemplateMinifyHtml = model.NewBool(`dev/template/minify_html`)

// PathDevTranslateInlineActive => Enabled for Storefront.
// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevTranslateInlineActive = model.NewBool(`dev/translate_inline/active`)

// PathDevTranslateInlineActiveAdmin => Enabled for Admin.
// Translate, blocks and other output caches should be disabled for both
// Storefront and Admin inline translations.
// BackendModel: Otnegam\Config\Model\Config\Backend\Translate
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevTranslateInlineActiveAdmin = model.NewBool(`dev/translate_inline/active_admin`)

// PathDevJsMergeFiles => Merge JavaScript Files.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevJsMergeFiles = model.NewBool(`dev/js/merge_files`)

// PathDevJsEnableJsBundling => Enable JavaScript Bundling.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevJsEnableJsBundling = model.NewBool(`dev/js/enable_js_bundling`)

// PathDevJsMinifyFiles => Minify JavaScript Files.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevJsMinifyFiles = model.NewBool(`dev/js/minify_files`)

// PathDevCssMergeCssFiles => Merge CSS Files.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevCssMergeCssFiles = model.NewBool(`dev/css/merge_css_files`)

// PathDevCssMinifyFiles => Minify CSS Files.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevCssMinifyFiles = model.NewBool(`dev/css/minify_files`)

// PathDevImageDefaultAdapter => Image Adapter.
// When the adapter was changed, please flush Catalog Images Cache.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image\Adapter
// SourceModel: Otnegam\Config\Model\Config\Source\Image\Adapter
var PathDevImageDefaultAdapter = model.NewStr(`dev/image/default_adapter`)

// PathDevStaticSign => Sign Static Files.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathDevStaticSign = model.NewBool(`dev/static/sign`)

// PathGeneralStoreInformationName => Store Name.
var PathGeneralStoreInformationName = model.NewStr(`general/store_information/name`)

// PathGeneralStoreInformationPhone => Store Phone Number.
var PathGeneralStoreInformationPhone = model.NewStr(`general/store_information/phone`)

// PathGeneralStoreInformationHours => Store Hours of Operation.
var PathGeneralStoreInformationHours = model.NewStr(`general/store_information/hours`)

// PathGeneralStoreInformationCountryId => Country.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralStoreInformationCountryId = model.NewStr(`general/store_information/country_id`)

// PathGeneralStoreInformationRegionId => Region/State.
var PathGeneralStoreInformationRegionId = model.NewStr(`general/store_information/region_id`)

// PathGeneralStoreInformationPostcode => ZIP/Postal Code.
var PathGeneralStoreInformationPostcode = model.NewStr(`general/store_information/postcode`)

// PathGeneralStoreInformationCity => City.
var PathGeneralStoreInformationCity = model.NewStr(`general/store_information/city`)

// PathGeneralStoreInformationStreetLine1 => Street Address.
var PathGeneralStoreInformationStreetLine1 = model.NewStr(`general/store_information/street_line1`)

// PathGeneralStoreInformationStreetLine2 => Street Address Line 2.
var PathGeneralStoreInformationStreetLine2 = model.NewStr(`general/store_information/street_line2`)

// PathGeneralStoreInformationMerchantVatNumber => VAT Number.
var PathGeneralStoreInformationMerchantVatNumber = model.NewStr(`general/store_information/merchant_vat_number`)

// PathGeneralSingleStoreModeEnabled => Enable Single-Store Mode.
// This setting will not be taken into account if system has more than one
// store view.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGeneralSingleStoreModeEnabled = model.NewBool(`general/single_store_mode/enabled`)

// PathSystemSmtpDisable => Disable Email Communications.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSystemSmtpDisable = model.NewBool(`system/smtp/disable`)

// PathSystemSmtpHost => Host.
// For Windows server only.
var PathSystemSmtpHost = model.NewStr(`system/smtp/host`)

// PathSystemSmtpPort => Port (25).
// For Windows server only.
var PathSystemSmtpPort = model.NewStr(`system/smtp/port`)

// PathSystemSmtpSetReturnPath => Set Return-Path.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesnocustom
var PathSystemSmtpSetReturnPath = model.NewBool(`system/smtp/set_return_path`)

// PathSystemSmtpReturnPathEmail => Return-Path Email.
// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Address
var PathSystemSmtpReturnPathEmail = model.NewStr(`system/smtp/return_path_email`)

// PathAdminEmailsForgotEmailTemplate => Forgot Password Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathAdminEmailsForgotEmailTemplate = model.NewStr(`admin/emails/forgot_email_template`)

// PathAdminEmailsForgotEmailIdentity => Forgot and Reset Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathAdminEmailsForgotEmailIdentity = model.NewStr(`admin/emails/forgot_email_identity`)

// PathAdminEmailsPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
// Please enter a number 1 or greater in this field.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Password\Link\Expirationperiod
var PathAdminEmailsPasswordResetLinkExpirationPeriod = model.NewStr(`admin/emails/password_reset_link_expiration_period`)

// PathAdminStartupMenuItemId => Startup Page.
// SourceModel: Otnegam\Config\Model\Config\Source\Admin\Page
var PathAdminStartupMenuItemId = model.NewStr(`admin/startup/menu_item_id`)

// PathAdminUrlUseCustom => Use Custom Admin URL.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usecustom
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminUrlUseCustom = model.NewBool(`admin/url/use_custom`)

// PathAdminUrlCustom => Custom Admin URL.
// Make sure that base URL ends with '/' (slash), e.g.
// http://yourdomain/magento/
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custom
var PathAdminUrlCustom = model.NewStr(`admin/url/custom`)

// PathAdminUrlUseCustomPath => Use Custom Admin Path.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminUrlUseCustomPath = model.NewBool(`admin/url/use_custom_path`)

// PathAdminUrlCustomPath => Custom Admin Path.
// You will have to sign in after you save your custom admin path.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Custompath
var PathAdminUrlCustomPath = model.NewStr(`admin/url/custom_path`)

// PathAdminSecurityUseFormKey => Add Secret Key to URLs.
// BackendModel: Otnegam\Config\Model\Config\Backend\Admin\Usesecretkey
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminSecurityUseFormKey = model.NewBool(`admin/security/use_form_key`)

// PathAdminSecurityUseCaseSensitiveLogin => Login is Case Sensitive.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminSecurityUseCaseSensitiveLogin = model.NewBool(`admin/security/use_case_sensitive_login`)

// PathAdminSecuritySessionLifetime => Admin Session Lifetime (seconds).
// Values less than 60 are ignored.
var PathAdminSecuritySessionLifetime = model.NewStr(`admin/security/session_lifetime`)

// PathAdminDashboardEnableCharts => Enable Charts.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathAdminDashboardEnableCharts = model.NewBool(`admin/dashboard/enable_charts`)

// PathWebUrlUseStore => Add Store Code to Urls.
// Warning! When using Store Code in URLs, in some cases system may not work
// properly if URLs without Store Codes are specified in the third party
// services (e.g. PayPal etc.).
// BackendModel: Otnegam\Config\Model\Config\Backend\Store
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebUrlUseStore = model.NewBool(`web/url/use_store`)

// PathWebUrlRedirectToBase => Auto-redirect to Base URL.
// I.e. redirect from http://example.com/store/ to
// http://www.example.com/store/
// SourceModel: Otnegam\Config\Model\Config\Source\Web\Redirect
var PathWebUrlRedirectToBase = model.NewStr(`web/url/redirect_to_base`)

// PathWebSeoUseRewrites => Use Web Server Rewrites.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSeoUseRewrites = model.NewBool(`web/seo/use_rewrites`)

// PathWebUnsecureBaseUrl => Base URL.
// Specify URL or {{base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebUnsecureBaseUrl = model.NewStr(`web/unsecure/base_url`)

// PathWebUnsecureBaseLinkUrl => Base Link URL.
// May start with {{unsecure_base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebUnsecureBaseLinkUrl = model.NewStr(`web/unsecure/base_link_url`)

// PathWebUnsecureBaseStaticUrl => Base URL for Static View Files.
// May be empty or start with {{unsecure_base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebUnsecureBaseStaticUrl = model.NewStr(`web/unsecure/base_static_url`)

// PathWebUnsecureBaseMediaUrl => Base URL for User Media Files.
// May be empty or start with {{unsecure_base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebUnsecureBaseMediaUrl = model.NewStr(`web/unsecure/base_media_url`)

// PathWebSecureBaseUrl => Secure Base URL.
// Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebSecureBaseUrl = model.NewStr(`web/secure/base_url`)

// PathWebSecureBaseLinkUrl => Secure Base Link URL.
// May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebSecureBaseLinkUrl = model.NewStr(`web/secure/base_link_url`)

// PathWebSecureBaseStaticUrl => Secure Base URL for Static View Files.
// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
// placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebSecureBaseStaticUrl = model.NewStr(`web/secure/base_static_url`)

// PathWebSecureBaseMediaUrl => Secure Base URL for User Media Files.
// May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}}
// placeholder.
// BackendModel: Otnegam\Config\Model\Config\Backend\Baseurl
var PathWebSecureBaseMediaUrl = model.NewStr(`web/secure/base_media_url`)

// PathWebSecureUseInFrontend => Use Secure URLs on Storefront.
// Enter https protocol to use Secure URLs on Storefront.
// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSecureUseInFrontend = model.NewBool(`web/secure/use_in_frontend`)

// PathWebSecureUseInAdminhtml => Use Secure URLs in Admin.
// Enter https protocol to use Secure URLs in Admin.
// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSecureUseInAdminhtml = model.NewBool(`web/secure/use_in_adminhtml`)

// PathWebSecureEnableHsts => Enable HTTP Strict Transport Security (HSTS).
// See HTTP Strict Transport Security page for details.
// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSecureEnableHsts = model.NewBool(`web/secure/enable_hsts`)

// PathWebSecureEnableUpgradeInsecure => Upgrade Insecure Requests.
// See Upgrade Insecure Requests page for details.
// BackendModel: Otnegam\Config\Model\Config\Backend\Secure
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSecureEnableUpgradeInsecure = model.NewBool(`web/secure/enable_upgrade_insecure`)

// PathWebSecureOffloaderHeader => Offloader header.
var PathWebSecureOffloaderHeader = model.NewStr(`web/secure/offloader_header`)

// PathWebDefaultFront => Default Web URL.
var PathWebDefaultFront = model.NewStr(`web/default/front`)

// PathWebDefaultNoRoute => Default No-route URL.
var PathWebDefaultNoRoute = model.NewStr(`web/default/no_route`)

// PathWebSessionUseRemoteAddr => Validate REMOTE_ADDR.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSessionUseRemoteAddr = model.NewBool(`web/session/use_remote_addr`)

// PathWebSessionUseHttpVia => Validate HTTP_VIA.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSessionUseHttpVia = model.NewBool(`web/session/use_http_via`)

// PathWebSessionUseHttpXForwardedFor => Validate HTTP_X_FORWARDED_FOR.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSessionUseHttpXForwardedFor = model.NewBool(`web/session/use_http_x_forwarded_for`)

// PathWebSessionUseHttpUserAgent => Validate HTTP_USER_AGENT.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSessionUseHttpUserAgent = model.NewBool(`web/session/use_http_user_agent`)

// PathWebSessionUseFrontendSid => Use SID on Storefront.
// Allows customers to stay logged in when switching between different stores.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathWebSessionUseFrontendSid = model.NewBool(`web/session/use_frontend_sid`)
