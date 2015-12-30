// +build ignore

package customer

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCustomerAccountShareScope => Share Customer Accounts.
// BackendModel: Otnegam\Customer\Model\Config\Share
// SourceModel: Otnegam\Customer\Model\Config\Share
var PathCustomerAccountShareScope = model.NewStr(`customer/account_share/scope`)

// PathCustomerCreateAccountAutoGroupAssign => Enable Automatic Assignment to Customer Group.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountAutoGroupAssign = model.NewBool(`customer/create_account/auto_group_assign`)

// PathCustomerCreateAccountTaxCalculationAddressType => Tax Calculation Based On.
// SourceModel: Otnegam\Customer\Model\Config\Source\Address\Type
var PathCustomerCreateAccountTaxCalculationAddressType = model.NewStr(`customer/create_account/tax_calculation_address_type`)

// PathCustomerCreateAccountDefaultGroup => Default Group.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountDefaultGroup = model.NewStr(`customer/create_account/default_group`)

// PathCustomerCreateAccountVivDomesticGroup => Group for Valid VAT ID - Domestic.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivDomesticGroup = model.NewStr(`customer/create_account/viv_domestic_group`)

// PathCustomerCreateAccountVivIntraUnionGroup => Group for Valid VAT ID - Intra-Union.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivIntraUnionGroup = model.NewStr(`customer/create_account/viv_intra_union_group`)

// PathCustomerCreateAccountVivInvalidGroup => Group for Invalid VAT ID.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivInvalidGroup = model.NewStr(`customer/create_account/viv_invalid_group`)

// PathCustomerCreateAccountVivErrorGroup => Validation Error Group.
// SourceModel: Otnegam\Customer\Model\Config\Source\Group
var PathCustomerCreateAccountVivErrorGroup = model.NewStr(`customer/create_account/viv_error_group`)

// PathCustomerCreateAccountVivOnEachTransaction => Validate on Each Transaction.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVivOnEachTransaction = model.NewBool(`customer/create_account/viv_on_each_transaction`)

// PathCustomerCreateAccountVivDisableAutoGroupAssignDefault => Default Value for Disable Automatic Group Changes Based on VAT ID.
// BackendModel: Otnegam\Customer\Model\Config\Backend\CreateAccount\DisableAutoGroupAssignDefault
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVivDisableAutoGroupAssignDefault = model.NewBool(`customer/create_account/viv_disable_auto_group_assign_default`)

// PathCustomerCreateAccountVatFrontendVisibility => Show VAT Number on Storefront.
// To show VAT number on Storefront, set Show VAT Number on Storefront option
// to Yes.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountVatFrontendVisibility = model.NewBool(`customer/create_account/vat_frontend_visibility`)

// PathCustomerCreateAccountEmailDomain => Default Email Domain.
var PathCustomerCreateAccountEmailDomain = model.NewStr(`customer/create_account/email_domain`)

// PathCustomerCreateAccountEmailTemplate => Default Welcome Email.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailTemplate = model.NewStr(`customer/create_account/email_template`)

// PathCustomerCreateAccountEmailNoPasswordTemplate => Default Welcome Email Without Password.
// This email will be sent instead of the Default Welcome Email, if a customer
// was created without password.  Email template chosen based on theme
// fallback when "Default" option is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailNoPasswordTemplate = model.NewStr(`customer/create_account/email_no_password_template`)

// PathCustomerCreateAccountEmailIdentity => Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCustomerCreateAccountEmailIdentity = model.NewStr(`customer/create_account/email_identity`)

// PathCustomerCreateAccountConfirm => Require Emails Confirmation.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountConfirm = model.NewBool(`customer/create_account/confirm`)

// PathCustomerCreateAccountEmailConfirmationTemplate => Confirmation Link Email.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailConfirmationTemplate = model.NewStr(`customer/create_account/email_confirmation_template`)

// PathCustomerCreateAccountEmailConfirmedTemplate => Welcome Email.
// This email will be sent instead of the Default Welcome Email, after account
// confirmation.  Email template chosen based on theme fallback when "Default"
// option is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerCreateAccountEmailConfirmedTemplate = model.NewStr(`customer/create_account/email_confirmed_template`)

// PathCustomerCreateAccountGenerateHumanFriendlyId => Generate Human-Friendly Customer ID.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerCreateAccountGenerateHumanFriendlyId = model.NewBool(`customer/create_account/generate_human_friendly_id`)

// PathCustomerPasswordForgotEmailTemplate => Forgot Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordForgotEmailTemplate = model.NewStr(`customer/password/forgot_email_template`)

// PathCustomerPasswordRemindEmailTemplate => Remind Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordRemindEmailTemplate = model.NewStr(`customer/password/remind_email_template`)

// PathCustomerPasswordResetPasswordTemplate => Reset Password Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCustomerPasswordResetPasswordTemplate = model.NewStr(`customer/password/reset_password_template`)

// PathCustomerPasswordForgotEmailIdentity => Password Template Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCustomerPasswordForgotEmailIdentity = model.NewStr(`customer/password/forgot_email_identity`)

// PathCustomerPasswordResetLinkExpirationPeriod => Recovery Link Expiration Period (days).
// Please enter a number 1 or greater in this field.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Password\Link\Expirationperiod
var PathCustomerPasswordResetLinkExpirationPeriod = model.NewStr(`customer/password/reset_link_expiration_period`)

// PathCustomerAddressStreetLines => Number of Lines in a Street Address.
// Leave empty for default (2). Valid range: 1-4
// BackendModel: Otnegam\Customer\Model\Config\Backend\Address\Street
var PathCustomerAddressStreetLines = model.NewStr(`customer/address/street_lines`)

// PathCustomerAddressPrefixShow => Show Prefix.
// The title that goes before name (Mr., Mrs., etc.)
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressPrefixShow = model.NewStr(`customer/address/prefix_show`)

// PathCustomerAddressPrefixOptions => Prefix Dropdown Options.
// Semicolon (;) separated values.Put semicolon in the beginning for empty
// first option.Leave empty for open text field.
var PathCustomerAddressPrefixOptions = model.NewStr(`customer/address/prefix_options`)

// PathCustomerAddressMiddlenameShow => Show Middle Name (initial).
// Always optional.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerAddressMiddlenameShow = model.NewBool(`customer/address/middlename_show`)

// PathCustomerAddressSuffixShow => Show Suffix.
// The suffix that goes after name (Jr., Sr., etc.)
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Address
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressSuffixShow = model.NewStr(`customer/address/suffix_show`)

// PathCustomerAddressSuffixOptions => Suffix Dropdown Options.
// Semicolon (;) separated values.Put semicolon in the beginning for empty
// first option.Leave empty for open text field.
var PathCustomerAddressSuffixOptions = model.NewStr(`customer/address/suffix_options`)

// PathCustomerAddressDobShow => Show Date of Birth.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressDobShow = model.NewStr(`customer/address/dob_show`)

// PathCustomerAddressTaxvatShow => Show Tax/VAT Number.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressTaxvatShow = model.NewStr(`customer/address/taxvat_show`)

// PathCustomerAddressGenderShow => Show Gender.
// BackendModel: Otnegam\Customer\Model\Config\Backend\Show\Customer
// SourceModel: Otnegam\Config\Model\Config\Source\Nooptreq
var PathCustomerAddressGenderShow = model.NewStr(`customer/address/gender_show`)

// PathCustomerStartupRedirectDashboard => Redirect Customer to Account Dashboard after Logging in.
// Customer will stay on the current page if "No" is selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCustomerStartupRedirectDashboard = model.NewBool(`customer/startup/redirect_dashboard`)

// PathCustomerAddressTemplatesText => Text.
var PathCustomerAddressTemplatesText = model.NewStr(`customer/address_templates/text`)

// PathCustomerAddressTemplatesOneline => Text One Line.
var PathCustomerAddressTemplatesOneline = model.NewStr(`customer/address_templates/oneline`)

// PathCustomerAddressTemplatesHtml => HTML.
var PathCustomerAddressTemplatesHtml = model.NewStr(`customer/address_templates/html`)

// PathCustomerAddressTemplatesPdf => PDF.
var PathCustomerAddressTemplatesPdf = model.NewStr(`customer/address_templates/pdf`)

// PathCustomerOnlineCustomersOnlineMinutesInterval => Online Minutes Interval.
// Leave empty for default (15 minutes).
var PathCustomerOnlineCustomersOnlineMinutesInterval = model.NewStr(`customer/online_customers/online_minutes_interval`)

// PathGeneralStoreInformationValidateVatNumber => .
var PathGeneralStoreInformationValidateVatNumber = model.NewStr(`general/store_information/validate_vat_number`)

// PathGeneralRestrictionAutocompleteOnStorefront => Enable Autocomplete on login/forgot password forms.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGeneralRestrictionAutocompleteOnStorefront = model.NewBool(`general/restriction/autocomplete_on_storefront`)
