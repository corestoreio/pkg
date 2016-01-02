// +build ignore

package customer

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCustomerAccountShareScope => Share Customer Accounts.
// BackendModel: Otnegam\Customer\Model\Config\Share
// SourceModel: Otnegam\Customer\Model\Config\Share
var PathCustomerAccountShareScope = model.NewStr(`customer/account_share/scope`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountAutoGroupAssign => Enable Automatic Assignment to Customer Group.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountAutoGroupAssign = model.NewBool(`customer/create_account/auto_group_assign`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountTaxCalculationAddressType => Tax Calculation Based On.
// SourceModel: Otnegam\Customer\Model\Config\Source\Address\Type
var PathCustomerCreateAccountTaxCalculationAddressType = model.NewStr(`customer/create_account/tax_calculation_address_type`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountDefaultGroup => Default Group.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountDefaultGroup = model.NewStr(`customer/create_account/default_group`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivDomesticGroup => Group for Valid VAT ID - Domestic.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivDomesticGroup = model.NewStr(`customer/create_account/viv_domestic_group`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivIntraUnionGroup => Group for Valid VAT ID - Intra-Union.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivIntraUnionGroup = model.NewStr(`customer/create_account/viv_intra_union_group`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivInvalidGroup => Group for Invalid VAT ID.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivInvalidGroup = model.NewStr(`customer/create_account/viv_invalid_group`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivErrorGroup => Validation Error Group.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivErrorGroup = model.NewStr(`customer/create_account/viv_error_group`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivOnEachTransaction => Validate on Each Transaction.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVivOnEachTransaction = model.NewBool(`customer/create_account/viv_on_each_transaction`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVivDisableAutoGroupAssignDefault => Default Value for Disable Automatic Group Changes Based on VAT ID.
// BackendModel: Otnegam\Customer\Model\Config\Backend\CreateAccount\DisableAutoGroupAssignDefault
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVivDisableAutoGroupAssignDefault = model.NewBool(`customer/create_account/viv_disable_auto_group_assign_default`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountVatFrontendVisibility => Show VAT Number on Storefront.
// To show VAT number on Storefront, set Show VAT Number on Storefront option
// to Yes.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVatFrontendVisibility = model.NewBool(`customer/create_account/vat_frontend_visibility`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailDomain => Default Email Domain.
var PathCustomerCreateAccountEmailDomain = model.NewStr(`customer/create_account/email_domain`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailTemplate => Default Welcome Email.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailTemplate = model.NewStr(`customer/create_account/email_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailNoPasswordTemplate => Default Welcome Email Without Password.
// This email will be sent instead of the Default Welcome Email, if a customer
// was created without password.  Email template chosen based on theme
// fallback when "Default" option is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailNoPasswordTemplate = model.NewStr(`customer/create_account/email_no_password_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailIdentity => Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCustomerCreateAccountEmailIdentity = model.NewStr(`customer/create_account/email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountConfirm => Require Emails Confirmation.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountConfirm = model.NewBool(`customer/create_account/confirm`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailConfirmationTemplate => Confirmation Link Email.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailConfirmationTemplate = model.NewStr(`customer/create_account/email_confirmation_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountEmailConfirmedTemplate => Welcome Email.
// This email will be sent instead of the Default Welcome Email, after account
// confirmation.  Email template chosen based on theme fallback when "Default"
// option is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailConfirmedTemplate = model.NewStr(`customer/create_account/email_confirmed_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerCreateAccountGenerateHumanFriendlyId => Generate Human-Friendly Customer ID.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountGenerateHumanFriendlyId = model.NewBool(`customer/create_account/generate_human_friendly_id`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerPasswordForgotEmailTemplate => Forgot Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordForgotEmailTemplate = model.NewStr(`customer/password/forgot_email_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerPasswordRemindEmailTemplate => Remind Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordRemindEmailTemplate = model.NewStr(`customer/password/remind_email_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerPasswordResetPasswordTemplate => Reset Password Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordResetPasswordTemplate = model.NewStr(`customer/password/reset_password_template`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerPasswordForgotEmailIdentity => Password Template Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCustomerPasswordForgotEmailIdentity = model.NewStr(`customer/password/forgot_email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
// Please enter a number 1 or greater in this field.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Password\Link\Expirationperiod
var PathCustomerPasswordResetLinkExpirationPeriod = model.NewStr(`customer/password/reset_link_expiration_period`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressStreetLines => Number of Lines in a Street Address.
// Leave empty for default (2). Valid range: 1-4
// BackendModel: Otnegam\Customer\Model\Config\Backend\Address\Street
var PathCustomerAddressStreetLines = model.NewStr(`customer/address/street_lines`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressPrefixShow => Show Prefix.
// The title that goes before name (Mr., Mrs., etc.)
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressPrefixShow = model.NewStr(`customer/address/prefix_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressPrefixOptions => Prefix Dropdown Options.
// Semicolon (;) separated values.Put semicolon in the beginning for empty
// first option.Leave empty for open text field.
var PathCustomerAddressPrefixOptions = model.NewStr(`customer/address/prefix_options`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressMiddlenameShow => Show Middle Name (initial).
// Always optional.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerAddressMiddlenameShow = model.NewBool(`customer/address/middlename_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressSuffixShow => Show Suffix.
// The suffix that goes after name (Jr., Sr., etc.)
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressSuffixShow = model.NewStr(`customer/address/suffix_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressSuffixOptions => Suffix Dropdown Options.
// Semicolon (;) separated values.Put semicolon in the beginning for empty
// first option.Leave empty for open text field.
var PathCustomerAddressSuffixOptions = model.NewStr(`customer/address/suffix_options`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressDobShow => Show Date of Birth.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressDobShow = model.NewStr(`customer/address/dob_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressTaxvatShow => Show Tax/VAT Number.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressTaxvatShow = model.NewStr(`customer/address/taxvat_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressGenderShow => Show Gender.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressGenderShow = model.NewStr(`customer/address/gender_show`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerStartupRedirectDashboard => Redirect Customer to Account Dashboard after Logging in.
// Customer will stay on the current page if "No" is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerStartupRedirectDashboard = model.NewBool(`customer/startup/redirect_dashboard`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressTemplatesText => Text.
var PathCustomerAddressTemplatesText = model.NewStr(`customer/address_templates/text`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressTemplatesOneline => Text One Line.
var PathCustomerAddressTemplatesOneline = model.NewStr(`customer/address_templates/oneline`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressTemplatesHtml => HTML.
var PathCustomerAddressTemplatesHtml = model.NewStr(`customer/address_templates/html`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerAddressTemplatesPdf => PDF.
var PathCustomerAddressTemplatesPdf = model.NewStr(`customer/address_templates/pdf`, model.WithPkgCfg(PackageConfiguration))

// PathCustomerOnlineCustomersOnlineMinutesInterval => Online Minutes Interval.
// Leave empty for default (15 minutes).
var PathCustomerOnlineCustomersOnlineMinutesInterval = model.NewStr(`customer/online_customers/online_minutes_interval`, model.WithPkgCfg(PackageConfiguration))

// PathGeneralStoreInformationValidateVatNumber => .
var PathGeneralStoreInformationValidateVatNumber = model.NewStr(`general/store_information/validate_vat_number`, model.WithPkgCfg(PackageConfiguration))

// PathGeneralRestrictionAutocompleteOnStorefront => Enable Autocomplete on login/forgot password forms.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGeneralRestrictionAutocompleteOnStorefront = model.NewBool(`general/restriction/autocomplete_on_storefront`, model.WithPkgCfg(PackageConfiguration))
